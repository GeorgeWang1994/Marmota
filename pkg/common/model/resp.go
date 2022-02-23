package model

type RpcResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg_opt"`
	Data interface{} `json:"data"`
}

type NullRpcRequest struct {
}
