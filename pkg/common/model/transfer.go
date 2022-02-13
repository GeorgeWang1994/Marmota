package model

type TransferResponse struct {
	Message string
	Total   int
	Invalid int
	Latency int64
}
