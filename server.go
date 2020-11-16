package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type ServerAdmin struct {
	procecess []*TcpData
}

func NewServerAdmin(limit uint64) *ServerAdmin {
	server := &ServerAdmin{}
	var process *TcpData

	for i := uint64(0); i < limit; i++ {
		process = &TcpData{}
		process.ID = i
		server.procecess = append(server.procecess, process)
	}

	return server
}

func main() {
	limit := uint64(5)
	server := NewServerAdmin(limit)
	listener, err := net.Listen("tcp", "localhost:4000")

	if err != nil {
		panic("Could not run server")
	}

	defer listener.Close()

	go AttendProcecess(server)
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("could not process request")
		}

		go RequestReceived(server, conn)
	}
}

func CountProcecess(server *ServerAdmin) {
	for _, process := range server.procecess {
		fmt.Printf("%d: %d\n", process.ID, process.Count)
	}
	fmt.Println("=============================================")
}

func IncrementCounts(procecess []*TcpData) {
	for _, process := range procecess {
		process.Count++
	}
}

func AttendProcecess(server *ServerAdmin) {
	for {
		CountProcecess(server)
		IncrementCounts(server.procecess)
		time.Sleep(time.Millisecond * 500)
	}
}

func RequestReceived(server *ServerAdmin, conn net.Conn) {
	defer conn.Close()
	var clientData *TcpData = server.procecess[0]
	server.procecess = server.procecess[1:]
	json.NewEncoder(conn).Encode(clientData)

	dataReader := bufio.NewScanner(conn)
	var dataBytes []byte
	for {
		if dataReader.Scan() {
			dataBytes = []byte(dataReader.Text())
			json.Unmarshal(dataBytes, clientData)
			server.procecess = append(server.procecess, clientData)
			return
		}
	}
}
