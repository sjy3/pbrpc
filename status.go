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
/*
 * @Author: Malin Xie
 * @Description:
 * @Date: 2021-07-26 17:09:25
 */
package baidurpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jhunters/timewheel"
)

// HttpStatusView
type HttpStatusView struct {
	server *TcpServer
}

func (hsv *HttpStatusView) Status(c context.Context) (*RPCStatus, context.Context) {
	s := hsv.server
	result := &RPCStatus{}
	if s.serverMeta.Host != nil {
		result.Host = *s.serverMeta.Host
	}
	if s.serverMeta.Port != nil {
		result.Port = int32(*s.serverMeta.Port)
	}
	if s.serverMeta.IdleTimeoutSeconds != nil {
		result.TimeoutSenconds = int32(*s.serverMeta.IdleTimeoutSeconds)
	}

	rpcServices := s.services
	methods := make([]*RPCMethod, len(rpcServices))
	var i int = 0
	for sname, service := range rpcServices {
		m := &RPCMethod{Service: service.GetServiceName(), Method: service.GetMethodName()}
		// query meta info
		serviceMeta, ok := s.servicesMeta[sname]
		if ok {
			if serviceMeta.InPbFieldMetas != nil {
				metaString, _ := json.Marshal(serviceMeta.InPbFieldMetas)
				m.InTypeMeta = string(metaString)
			}
			if serviceMeta.RetrunPbFieldMetas != nil {
				metaString, _ := json.Marshal(serviceMeta.RetrunPbFieldMetas)
				m.ReturnTypeMeta = string(metaString)
			}
		}

		methods[i] = m
		i++
	}
	result.Methods = methods
	return result, c
}

func (hsv *HttpStatusView) QpsDataStatus(c context.Context, method *RPCMethod) (*QpsData, context.Context) {
	serviceId := GetServiceId(method.Service, method.Method)
	ret := &QpsData{Qpsinfo: make(map[int64]int32)}
	requestStatus, ok := hsv.server.requestStatus.Methods[serviceId]
	if ok {
		ret.Qpsinfo = requestStatus.QpsStatus
	}
	// add current current
	ret.Qpsinfo[time.Now().Unix()] += 0
	return ret, c
}

// RPCRequestStatus
type RPCRequestStatus struct {
	Methods map[string]*RPCMethodReuqestStatus

	reqeustChan chan request

	closeChan chan bool

	expireAfterSecs int16

	started bool

	tw *timewheel.TimeWheel
}

type request struct {
	method string
	t      time.Time
	count  int
}

// RPCMethodReuqestStatus
type RPCMethodReuqestStatus struct {
	QpsStatus map[int64]int32
}

// NewRPCRequestStatus
func NewRPCRequestStatus(services map[string]Service) *RPCRequestStatus {
	ret := &RPCRequestStatus{
		Methods:     make(map[string]*RPCMethodReuqestStatus, len(services)),
		reqeustChan: make(chan request, 1024),
		closeChan:   make(chan bool),
	}

	for name := range services {
		ret.Methods[name] = &RPCMethodReuqestStatus{QpsStatus: make(map[int64]int32, 1024)}
	}

	return ret
}

// Start
func (r *RPCRequestStatus) Start() error {
	Infof("RPC method reuqest status record starting. expire time within %d seconds ", r.expireAfterSecs)
	r.started = true

	// start time wheel to delete expire data
	tw, err := timewheel.New(1*time.Second, uint16(r.expireAfterSecs))
	if err != nil {
		r.started = false
		return err
	}
	r.tw = tw
	r.tw.Start()

	for {
		select {
		case m := <-r.reqeustChan:
			status, ok := r.Methods[m.method]
			if !ok {
				status = &RPCMethodReuqestStatus{QpsStatus: make(map[int64]int32, 1024)}
				r.Methods[m.method] = status
			}
			k := m.t.Unix()
			count, ok := status.QpsStatus[k]
			if !ok {
				count = int32(m.count)
				// add task
				task := timewheel.Task{
					Data: m,
					TimeoutCallback: func(tt timewheel.Task) { // call back function on time out
						k := tt.Data.(request)
						r.expire(k.method, k.t)

					}}
				// add task and return unique task id
				r.tw.AddTask(time.Duration(r.expireAfterSecs)*time.Second, task) // add delay task
			} else {
				count += int32(m.count)
			}
			status.QpsStatus[k] = count

		case <-r.closeChan:
			r.started = false
			return nil
		}
	}

}

// RequestIn
func (r *RPCRequestStatus) RequestIn(methodName string, t time.Time, count int) error {
	if !r.started {
		return fmt.Errorf("RequestIn failed status not started")
	}
	req := request{method: methodName, t: t, count: count}
	r.reqeustChan <- req

	return nil
}

// expire  remove qps data after time expire
func (r *RPCRequestStatus) expire(methodName string, t time.Time) {
	status, ok := r.Methods[methodName]
	if ok {
		delete(status.QpsStatus, t.Unix())
	}
}

// Stop
func (r *RPCRequestStatus) Stop() {
	if !r.started {
		return
	}
	r.started = false
	r.closeChan <- true

	r.tw.Stop()
}
