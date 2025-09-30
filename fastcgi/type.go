package fastcgi

import "net"

type Socket struct {
	socketPath string
	listener   net.Listener
}

type RequestContext struct {
	RequestId uint16
	Params    map[string]string
	StdIn     []byte
	StdOut    []byte
}

const DefaultSocketPath = "/var/run/php-fubar.sock"

type RecordType uint8

const (
	BeginRequest RecordType = iota + 1
	AbortRequest
	EndRequest
	Params
	StdIn
	StdOut
	StdErr
	Data
)

type Header struct {
	version       uint8
	Type          RecordType
	RequestId     uint16
	ContentLength uint16
	PaddingLength uint8
	_             uint8
}

type BeginRequestBody struct {
	Role  uint16
	Flags byte
	_     [5]byte
}

type EndRequestBody struct {
	AppStatus      uint32
	ProtocolStatus uint8
	_              [3]byte
}
