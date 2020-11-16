package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:4000")

	if err != nil {
		panic("Could not connect to the server")
	}

	defer conn.Close()

	tcpReader := bufio.NewReader(conn)
	scan := bufio.NewScanner(os.Stdin)

	clientData := &TcpData{}
	err = json.NewDecoder(tcpReader).Decode(clientData)
	if err != nil {
		panic("Error parsing data")
	}

	onSignalInterruption(conn, clientData)
	closeClient := make(chan bool)
	go RunClientProcess(closeClient, conn, clientData)
	for {
		if scan.Scan() {
			closeClient <- true
			CloseClientConnection(conn, clientData)
			conn.Close()
			return
		}
	}
}

func onSignalInterruption(server net.Conn, clientData *TcpData) {
	signalChannel := make(chan os.Signal)

	signal.Notify(signalChannel, syscall.SIGTERM)
	go func() {
		<-signalChannel
		CloseClientConnection(server, clientData)
		os.Exit(0)
	}()
}

func CloseClientConnection(server net.Conn, clientData *TcpData) {
	json.NewEncoder(server).Encode(clientData)
}

func RunClientProcess(close chan bool, conn net.Conn, clientData *TcpData) {
	for {
		select {
		case <-close:
			CloseClientConnection(conn, clientData)
			return
		default:
			fmt.Printf("%d: %d\n", clientData.ID, clientData.Count)
			clientData.Count++
			time.Sleep(time.Millisecond * 500)
		}
	}
}
