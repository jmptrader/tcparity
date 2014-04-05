package main

import (
	"log"
	"net"
)

// manager is reponsible for coordinating the application and handling the main event loop
func manager(killChan chan bool, exitChan chan int) {
	// Channels to handle incoming requests and outgoing responses
	reqChan := make(chan net.Conn, 1000)
	resChan := make(chan net.Conn, 1000)

	// Channel to receive statistics from the application
	statChan := make(chan int, 1000)

	// Channel to handle errors from the application
	errChan := make(chan error, 1000)

	// Channel to trigger a graceful halt for various components
	haltChan := make(chan bool, 10)

	// Start the listener, to accept requests and feed them into the request channel
	go listener(reqChan, errChan, haltChan)

	// Loop and handle events
	for {
		select {
		// Stop the application
		case <-killChan:
			// Close all manager channels
			close(reqChan)
			close(resChan)
			close(statChan)
			close(errChan)
			close(haltChan)

			// Trigger graceful shutdown
			exitChan <- 0
			break
		// Handle incoming requests, send them to workers
		case request := <-reqChan:
			go worker(request, resChan, errChan)
			break
		// Handle outgoing responses, send them to transporters
		case response := <-resChan:
			go transporter(response, errChan)
			_ = response
			break
		// Handle application statistic processing
		case stat := <-statChan:
			go statworker(stat, errChan)
			break
		// Handle application errors
		case err := <-errChan:
			if err != nil {
				log.Println(err)
			}
			break
		}
	}
}
