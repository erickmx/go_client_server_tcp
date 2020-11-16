package main

type TcpData struct {
	ID    uint64 `json:"id"`
	Count uint64 `json:"count"`
}

func (tcpData TcpData) New(id uint64) *TcpData {
	return &TcpData{
		ID:    id,
		Count: 0,
	}
}

func (tcpData TcpData) From(original *TcpData) *TcpData {
	return &TcpData{
		ID:    original.ID,
		Count: original.Count,
	}
}
