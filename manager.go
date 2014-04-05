package main

import ()

// manager is reponsible for coordinating the application and handling the main event loop
func manager(killChan chan bool, exitChan chan int) {
	// Channels to handle incoming requests and outgoing responses
	reqChan := make(chan []byte, 1000)
	resChan := make(chan []byte, 1000)

	// Loop and handle events
	for {
		select {
		// Stop the application
		case <-killChan:
			// Close manager channels
			close(reqChan)
			close(resChan)

			// Trigger graceful shutdown
			exitChan <- 0
			break
		// Handle incoming requests
		case request := <-reqChan:
			_ = request
			break
		// Handle outgoing responses
		case response := <-resChan:
			_ = response
			break
		}
	}
}
