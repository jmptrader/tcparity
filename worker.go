package main

import (
	"net"
)

// worker processes requests, and sends them to the response channel
func worker(request net.Conn, resChan chan net.Conn, errChan chan error) {

}
