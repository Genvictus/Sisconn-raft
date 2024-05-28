package transport

type RPCResponse struct {
	Response interface{}
	Error error
}

type RPC struct {
	Command interface{}
	RespChan chan <- RPC
}