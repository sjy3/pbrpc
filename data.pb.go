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

// Code generated by protoc-gen-go.
// source: Request.proto
// DO NOT EDIT!
/*
enum Errno {
    // Errno caused by client
    ENOSERVICE              = 1001;  // Service not found
    ENOMETHOD               = 1002;  // Method not found
    EREQUEST                = 1003;  // Bad Request
    ERPCAUTH                = 1004;  // Unauthorized, can't be called EAUTH
                                     // directly which is defined in MACOSX
    ETOOMANYFAILS           = 1005;  // Too many sub calls failed
    EPCHANFINISH            = 1006;  // [Internal] ParallelChannel finished
    EBACKUPREQUEST          = 1007;  // Sending backup request
    ERPCTIMEDOUT            = 1008;  // RPC call is timed out
    EFAILEDSOCKET           = 1009;  // Broken socket
    EHTTP                   = 1010;  // Bad http call
    EOVERCROWDED            = 1011;  // The server is overcrowded
    ERTMPPUBLISHABLE        = 1012;  // RtmpRetryingClientStream is publishable
    ERTMPCREATESTREAM       = 1013;  // createStream was rejected by the RTMP server
    EEOF                    = 1014;  // Got EOF
    EUNUSED                 = 1015;  // The socket was not needed
    ESSL                    = 1016;  // SSL related error
    EH2RUNOUTSTREAMS        = 1017;  // The H2 socket was run out of streams
    EREJECT                 = 1018;  // The Request is rejected

    // Errno caused by server
    EINTERNAL               = 2001;  // Internal Server Error
    ERESPONSE               = 2002;  // Bad Response
    ELOGOFF                 = 2003;  // Server is stopping
    ELIMIT                  = 2004;  // Reached server's limit on resources
    ECLOSE                  = 2005;  // Close socket initiatively
    EITP                    = 2006;  // Failed Itp response
}


*/

/*
Package baidurpc is a generated protocol buffer package.

It is generated from these files:
	Request.proto

It has these top-level messages:
	Request
	Response
	ChunkInfo
	RpcMeta
*/
package baidurpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Request struct {
	ServiceName      *string `protobuf:"bytes,1,req,name=serviceName" json:"serviceName,omitempty"`
	MethodName       *string `protobuf:"bytes,2,req,name=methodName" json:"methodName,omitempty"`
	LogId            *int64  `protobuf:"varint,3,opt,name=logId" json:"logId,omitempty"`
	ExtraParam       []byte  `protobuf:"bytes,4,opt,name=extraParam" json:"extraParam,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetServiceName() string {
	if m != nil && m.ServiceName != nil {
		return *m.ServiceName
	}
	return ""
}

func (m *Request) GetMethodName() string {
	if m != nil && m.MethodName != nil {
		return *m.MethodName
	}
	return ""
}

func (m *Request) GetLogId() int64 {
	if m != nil && m.LogId != nil {
		return *m.LogId
	}
	return 0
}

func (m *Request) GetExtraParam() []byte {
	if m != nil {
		return m.ExtraParam
	}
	return nil
}

type Response struct {
	ErrorCode        *int32  `protobuf:"varint,1,opt,name=errorCode" json:"errorCode,omitempty"`
	ErrorText        *string `protobuf:"bytes,2,opt,name=errorText" json:"errorText,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Response) GetErrorCode() int32 {
	if m != nil && m.ErrorCode != nil {
		return *m.ErrorCode
	}
	return 0
}

func (m *Response) GetErrorText() string {
	if m != nil && m.ErrorText != nil {
		return *m.ErrorText
	}
	return ""
}

type ChunkInfo struct {
	StreamId         *int64 `protobuf:"varint,1,req,name=streamId" json:"streamId,omitempty"`
	ChunkId          *int64 `protobuf:"varint,2,req,name=chunkId" json:"chunkId,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *ChunkInfo) Reset()                    { *m = ChunkInfo{} }
func (m *ChunkInfo) String() string            { return proto.CompactTextString(m) }
func (*ChunkInfo) ProtoMessage()               {}
func (*ChunkInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ChunkInfo) GetStreamId() int64 {
	if m != nil && m.StreamId != nil {
		return *m.StreamId
	}
	return 0
}

func (m *ChunkInfo) GetChunkId() int64 {
	if m != nil && m.ChunkId != nil {
		return *m.ChunkId
	}
	return 0
}

type RpcMeta struct {
	Request            *Request   `protobuf:"bytes,1,opt,name=request" json:"request,omitempty"`
	Response           *Response  `protobuf:"bytes,2,opt,name=response" json:"response,omitempty"`
	CompressType       *int32     `protobuf:"varint,3,opt,name=compressType" json:"compressType,omitempty"`
	CorrelationId      *int64     `protobuf:"varint,4,opt,name=correlationId" json:"correlationId,omitempty"`
	AttachmentSize     *int32     `protobuf:"varint,5,opt,name=attachmentSize" json:"attachmentSize,omitempty"`
	ChunkInfo          *ChunkInfo `protobuf:"bytes,6,opt,name=chunkInfo" json:"chunkInfo,omitempty"`
	AuthenticationData []byte     `protobuf:"bytes,7,opt,name=authenticationData" json:"authenticationData,omitempty"`
	XXX_unrecognized   []byte     `json:"-"`
}

func (m *RpcMeta) Reset()                    { *m = RpcMeta{} }
func (m *RpcMeta) String() string            { return proto.CompactTextString(m) }
func (*RpcMeta) ProtoMessage()               {}
func (*RpcMeta) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *RpcMeta) GetRequest() *Request {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *RpcMeta) GetResponse() *Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *RpcMeta) GetCompressType() int32 {
	if m != nil && m.CompressType != nil {
		return *m.CompressType
	}
	return 0
}

func (m *RpcMeta) GetCorrelationId() int64 {
	if m != nil && m.CorrelationId != nil {
		return *m.CorrelationId
	}
	return 0
}

func (m *RpcMeta) GetAttachmentSize() int32 {
	if m != nil && m.AttachmentSize != nil {
		return *m.AttachmentSize
	}
	return 0
}

func (m *RpcMeta) GetChunkInfo() *ChunkInfo {
	if m != nil {
		return m.ChunkInfo
	}
	return nil
}

func (m *RpcMeta) GetAuthenticationData() []byte {
	if m != nil {
		return m.AuthenticationData
	}
	return nil
}

func init() {
	proto.RegisterType((*Request)(nil), "baidurpc.Request")
	proto.RegisterType((*Response)(nil), "baidurpc.Response")
	proto.RegisterType((*ChunkInfo)(nil), "baidurpc.ChunkInfo")
	proto.RegisterType((*RpcMeta)(nil), "baidurpc.RpcMeta")
}

func init() { proto.RegisterFile("Request.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 361 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x91, 0xdf, 0x6a, 0xdb, 0x30,
	0x14, 0xc6, 0xb1, 0x3d, 0xcf, 0xf6, 0x49, 0x32, 0xd8, 0xd9, 0x2e, 0xc4, 0x18, 0xc3, 0x98, 0x31,
	0x0c, 0x03, 0xc3, 0xf2, 0x06, 0x23, 0x63, 0xe0, 0x8b, 0x8d, 0xa2, 0xe6, 0x05, 0x54, 0xf9, 0xb4,
	0x36, 0x8d, 0x2d, 0x57, 0x96, 0x4b, 0xda, 0xbb, 0x3e, 0x57, 0x5f, 0xae, 0x44, 0x89, 0xff, 0xa4,
	0xf4, 0xf2, 0xfc, 0xbe, 0x4f, 0xd2, 0x77, 0x3e, 0xc1, 0x8a, 0xd3, 0x5d, 0x4f, 0x9d, 0xc9, 0x5a,
	0xad, 0x8c, 0xc2, 0xf0, 0x4a, 0x54, 0x45, 0xaf, 0x5b, 0x99, 0x3c, 0x39, 0x10, 0x9c, 0x34, 0x8c,
	0x61, 0xd1, 0x91, 0xbe, 0xaf, 0x24, 0xfd, 0x17, 0x35, 0x31, 0x27, 0x76, 0xd3, 0x88, 0xcf, 0x11,
	0x7e, 0x03, 0xa8, 0xc9, 0x94, 0xaa, 0xb0, 0x06, 0xd7, 0x1a, 0x66, 0x04, 0x3f, 0x83, 0xbf, 0x53,
	0x37, 0x79, 0xc1, 0xbc, 0xd8, 0x49, 0x3d, 0x7e, 0x1c, 0x0e, 0xa7, 0x68, 0x6f, 0xb4, 0xb8, 0x10,
	0x5a, 0xd4, 0xec, 0x5d, 0xec, 0xa4, 0x4b, 0x3e, 0x23, 0xc9, 0x5f, 0x08, 0x39, 0x75, 0xad, 0x6a,
	0x3a, 0xc2, 0xaf, 0x10, 0x91, 0xd6, 0x4a, 0x6f, 0x54, 0x71, 0x48, 0xe0, 0xa4, 0x3e, 0x9f, 0xc0,
	0xa8, 0x6e, 0x69, 0x6f, 0x98, 0x1b, 0x3b, 0x69, 0xc4, 0x27, 0x90, 0xfc, 0x86, 0x68, 0x53, 0xf6,
	0xcd, 0x6d, 0xde, 0x5c, 0x2b, 0xfc, 0x02, 0x61, 0x67, 0x34, 0x89, 0x3a, 0x2f, 0xec, 0x26, 0x1e,
	0x1f, 0x67, 0x64, 0x10, 0x48, 0x6b, 0x2c, 0xec, 0x0e, 0x1e, 0x1f, 0xc6, 0xe4, 0xd9, 0x85, 0x80,
	0xb7, 0xf2, 0x1f, 0x19, 0x81, 0x3f, 0x21, 0xd0, 0xc7, 0x66, 0x6c, 0x90, 0xc5, 0xfa, 0x63, 0x36,
	0xd4, 0x96, 0x9d, 0x2a, 0xe3, 0x83, 0x03, 0x33, 0x08, 0xf5, 0x69, 0x07, 0x1b, 0x6c, 0xb1, 0xc6,
	0xb9, 0xfb, 0xa8, 0xf0, 0xd1, 0x83, 0x09, 0x2c, 0xa5, 0xaa, 0x5b, 0x4d, 0x5d, 0xb7, 0x7d, 0x68,
	0xc9, 0x16, 0xe6, 0xf3, 0x33, 0x86, 0xdf, 0x61, 0x25, 0x95, 0xd6, 0xb4, 0x13, 0xa6, 0x52, 0x4d,
	0x5e, 0xd8, 0xea, 0x3c, 0x7e, 0x0e, 0xf1, 0x07, 0x7c, 0x10, 0xc6, 0x08, 0x59, 0xd6, 0xd4, 0x98,
	0xcb, 0xea, 0x91, 0x98, 0x6f, 0xef, 0x7a, 0x45, 0xf1, 0x17, 0x44, 0x72, 0x68, 0x87, 0xbd, 0xb7,
	0x11, 0x3f, 0x4d, 0x11, 0xc7, 0xe2, 0xf8, 0xe4, 0xc2, 0x0c, 0x50, 0xf4, 0xa6, 0xa4, 0xc6, 0x54,
	0xd2, 0x3e, 0xf7, 0x47, 0x18, 0xc1, 0x02, 0xfb, 0x81, 0x6f, 0x28, 0x2f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x7d, 0x01, 0x53, 0x40, 0x66, 0x02, 0x00, 0x00,
}
