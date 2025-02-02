// Go support for Protocol Buffers RPC which compatiable with https://github.com/Baidu-ecom/Jprotobuf-rpc-socket
//
// Copyright 2002-2007 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package baidurpc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jhunters/timewheel"
	"google.golang.org/protobuf/proto"
)

var (
	defaultTimewheelInterval = 10 * time.Millisecond
	defaultTimewheelSlot     = 300

	errNeedInit               = errors.New("[client-001]Session is not initialized, Please use NewRpcInvocation() to create instance")
	errResponseNil            = errors.New("[client-003]No response result, mybe net work break error")
	LOG_SERVER_REPONSE_ERROR  = "[client-002]Server response error. code=%d, msg='%s'"
	LOG_CLIENT_TIMECOUST_INFO = "[client-101]Server name '%s' method '%s' process cost '%.5g' seconds"

	closedTimeOut = time.Duration(0)
)

const (
	ST_READ_TIMEOUT = 62
)

/*
RPC client invoke
*/
type RpcClient struct {
	Session Connection
	tw      *timewheel.TimeWheel

	// 单次请求唯一标识
	correlationId int64
	// async request state map
	requestCallState sync.Map //  use sync map for cocurrent access

	// to close loop receive
	closeChan chan bool

	asyncMode bool
}

// URL with host and port attribute
type URL struct {
	Host *string
	Port *int
}

// SetHost set host name
func (u *URL) SetHost(host *string) *URL {
	u.Host = host
	return u
}

// SetPort set port
func (u *URL) SetPort(port *int) *URL {
	u.Port = port
	return u
}

// RpcInvocation define rpc invocation
type RpcInvocation struct {
	ServiceName       *string
	MethodName        *string
	ParameterIn       *proto.Message
	Attachment        []byte
	LogId             *int64
	CompressType      *int32
	AuthenticateData  []byte
	ChunkSize         uint32
	TraceId           int64
	SpanId            int64
	ParentSpanId      int64
	RpcRequestMetaExt map[string]string
}

// NewRpcCient new rpc client
func NewRpcCient(connection Connection) (*RpcClient, error) {
	return NewRpcCientWithTimeWheelSetting(connection, defaultTimewheelInterval, uint16(defaultTimewheelSlot))
}

// NewRpcCientWithTimeWheelSetting new rpc client with set timewheel settings
func NewRpcCientWithTimeWheelSetting(connection Connection, timewheelInterval time.Duration, timewheelSlot uint16) (*RpcClient, error) {
	c := RpcClient{}
	c.Session = connection

	// async mode not support under pooled connection
	_, pooled := connection.(*TCPConnectionPool)
	c.asyncMode = !pooled

	// initial timewheel to process async request on time out event handle
	c.tw, _ = timewheel.New(timewheelInterval, timewheelSlot)
	c.tw.Start()
	c.closeChan = make(chan bool, 1)
	c.requestCallState = sync.Map{} // make(map[int64]chan *RpcDataPackage)

	if c.asyncMode { // only enabled on async mode
		go c.startLoopReceive()
	}
	return &c, nil
}

// NewRpcInvocation create RpcInvocation with service name and method name
func NewRpcInvocation(serviceName, methodName *string) *RpcInvocation {
	r := new(RpcInvocation)
	r.init(serviceName, methodName)

	return r
}

func (r *RpcInvocation) init(serviceName, methodName *string) {

	*r = RpcInvocation{}
	r.ServiceName = serviceName
	r.MethodName = methodName
	compressType := COMPRESS_NO
	r.CompressType = &compressType
	r.ParameterIn = nil
}

// SetParameterIn
func (r *RpcInvocation) SetParameterIn(parameterIn proto.Message) {
	r.ParameterIn = &parameterIn
}

// GetRequestRpcDataPackage
func (r *RpcInvocation) GetRequestRpcDataPackage() (*RpcDataPackage, error) {

	rpcDataPackage := new(RpcDataPackage)
	rpcDataPackage.ServiceName(*r.ServiceName)
	rpcDataPackage.MethodName(*r.MethodName)
	rpcDataPackage.MagicCode(MAGIC_CODE)
	rpcDataPackage.AuthenticationData(r.AuthenticateData)
	rpcDataPackage.chunkSize = r.ChunkSize
	rpcDataPackage.TraceId(r.TraceId)
	rpcDataPackage.SpanId(r.SpanId)
	rpcDataPackage.ParentSpanId(r.ParentSpanId)
	rpcDataPackage.RpcRequestMetaExt(r.RpcRequestMetaExt)
	if r.CompressType != nil {
		rpcDataPackage.CompressType(*r.CompressType)
	}
	if r.LogId != nil {
		rpcDataPackage.LogId(*r.LogId)
	}

	rpcDataPackage.SetAttachment(r.Attachment)

	if r.ParameterIn != nil {
		data, err := proto.Marshal(*r.ParameterIn)
		if err != nil {
			return nil, err
		}
		rpcDataPackage.SetData(data)
	}

	return rpcDataPackage, nil
}

// define client methods
// Close close client with time wheel
func (c *RpcClient) Close() {
	c.closeChan <- true
	if c.tw != nil {
		c.tw.Stop()
	}
}

func (c *RpcClient) startLoopReceive() {
	for {

		select {
		case <-c.closeChan:
			// exit loop
			return
		default:
			dataPackage, err := c.safeReceive()
			if err != nil {

				netErr, ok := err.(*net.OpError)
				if ok {
					// if met network error, wait some time to retry or call client close method to close loop if met net error
					// error maybe about broken network or closed network
					log.Println(netErr)
					c.Session.Reconnect() // try reconnect

				}
				time.Sleep(200 * time.Millisecond)

			}

			if dataPackage != nil && dataPackage.Meta != nil {
				correlationId := dataPackage.Meta.GetCorrelationId()
				v, exist := c.requestCallState.LoadAndDelete(correlationId) // [correlationId]
				if !exist {
					// bad response correlationId
					Errorf("bad correlationId '%d' not exist ", correlationId)
					continue
				}
				ch := v.(chan *RpcDataPackage)
				go func() {
					ch <- dataPackage
				}()
			}
		}

	}
}

func (c *RpcClient) safeReceive() (*RpcDataPackage, error) {
	defer func() {
		if p := recover(); p != nil {
			Warningf("receive catched panic error %v", p)
		}
	}()
	return c.Session.Receive()
}

// asyncRequest
func (c *RpcClient) asyncRequest(timeout time.Duration, request *RpcDataPackage, ch chan *RpcDataPackage) {
	// create a task bind with key, data and  time out call back function.
	t := &timewheel.Task{
		Data: nil, // business data
		TimeoutCallback: func(task timewheel.Task) { // call back function on time out
			// process someting after time out happened.
			errorcode := int32(ST_READ_TIMEOUT)
			request.ErrorCode(errorcode)
			errormsg := fmt.Sprintf("request time out of %v", task.Delay())
			request.ErrorText(errormsg)
			ch <- request
		}}

	// add task and return unique task id
	taskid, err := c.tw.AddTask(timeout, *t) // add delay task
	if err != nil {
		errorcode := int32(ST_ERROR)
		request.ErrorCode(errorcode)
		errormsg := err.Error()
		request.ErrorText(errormsg)

		ch <- request
		return
	}

	defer func() {
		c.tw.RemoveTask(taskid)
		if e := recover(); e != nil {
			Warningf("asyncRequest failed with error %v", e)
		}
	}()

	rsp, err := c.doSendReceive(request, ch)
	if err != nil {
		errorcode := int32(ST_ERROR)
		request.ErrorCode(errorcode)
		errormsg := err.Error()
		request.ErrorText(errormsg)

		ch <- request
		return
	}

	ch <- rsp
}

func (c *RpcClient) doSendReceive(rpcDataPackage *RpcDataPackage, ch <-chan *RpcDataPackage) (*RpcDataPackage, error) {
	if c.asyncMode {
		err := c.Session.Send(rpcDataPackage)
		if err != nil {
			return nil, err
		}
		// async wait response
		return <-ch, nil
	}
	// not async mode use block request
	return c.Session.SendReceive(rpcDataPackage)

}

// SendRpcRequest send rpc request to remote server
func (c *RpcClient) SendRpcRequest(rpcInvocation *RpcInvocation, responseMessage proto.Message) (*RpcDataPackage, error) {
	return c.SendRpcRequestWithTimeout(closedTimeOut, rpcInvocation, responseMessage)

}

// SendRpcRequest send rpc request to remote server
func (c *RpcClient) SendRpcRequestWithTimeout(timeout time.Duration, rpcInvocation *RpcInvocation, responseMessage proto.Message) (*RpcDataPackage, error) {
	if c.Session == nil {
		return nil, errNeedInit
	}

	now := time.Now().UnixNano()

	rpcDataPackage, err := rpcInvocation.GetRequestRpcDataPackage()
	if err != nil {
		return nil, err
	}

	// set request unique id
	correlationId := atomic.AddInt64(&c.correlationId, 1)
	rpcDataPackage.CorrelationId(correlationId)

	var rsp *RpcDataPackage
	if c.asyncMode {
		ch := make(chan *RpcDataPackage, 1)
		c.requestCallState.Store(correlationId, ch)
		// c.requestCallState[correlationId] = ch

		if timeout > 0 {
			go c.asyncRequest(timeout, rpcDataPackage, ch)
			rsp = <-ch
		} else {
			rsp, err = c.doSendReceive(rpcDataPackage, ch)
		}

	} else {
		if timeout > 0 {
			ch := make(chan *RpcDataPackage, 1)
			go c.asyncRequest(timeout, rpcDataPackage, ch)
			defer close(ch)
			// wait for message
			rsp = <-ch
		} else {
			rsp, err = c.Session.SendReceive(rpcDataPackage)
		}
	}

	if err != nil {
		errorcode := int32(ST_ERROR)
		rpcDataPackage.ErrorCode(errorcode)
		errormsg := err.Error()
		rpcDataPackage.ErrorText(errormsg)
		return rpcDataPackage, err
	}

	r := rsp
	if r == nil {
		return nil, errResponseNil //to ignore this nil value
	}

	errorCode := r.GetMeta().GetResponse().GetErrorCode()
	if errorCode > 0 {
		errMsg := fmt.Sprintf(LOG_SERVER_REPONSE_ERROR,
			errorCode, r.GetMeta().GetResponse().GetErrorText())
		return r, errors.New(errMsg)
	}

	response := r.GetData()
	if response != nil {
		err = proto.Unmarshal(response, responseMessage)
		if err != nil {
			return r, err
		}
	}

	took := TimetookInSeconds(now)
	Infof(LOG_CLIENT_TIMECOUST_INFO, *rpcInvocation.ServiceName, *rpcInvocation.MethodName, took)

	return r, nil

}

// RpcResult Rpc response result from client request api under asynchronous way
type RpcResult struct {
	rpcData *RpcDataPackage
	err     error
	message proto.Message
}

func (rr *RpcResult) Get() proto.Message {
	return rr.message
}

func (rr *RpcResult) GetRpcDataPackage() *RpcDataPackage {
	return rr.rpcData
}

func (rr *RpcResult) GetErr() error {
	return rr.err
}

// SendRpcRequestAsyc send rpc request to remote server in asynchronous way
func (c *RpcClient) SendRpcRequestAsyc(rpcInvocation *RpcInvocation, responseMessage proto.Message) <-chan *RpcResult {
	ch := make(chan *RpcResult, 1)

	go func() {
		defer func() {
			if p := recover(); p != nil {
				if err, ok := p.(error); ok {
					r := &RpcResult{nil, err, responseMessage}
					ch <- r
				}
			}
		}()

		resp, err := c.SendRpcRequest(rpcInvocation, responseMessage)
		result := &RpcResult{resp, err, responseMessage}
		ch <- result
	}()

	return ch
}
