package main

import ()

// manager is reponsible for coordinating the application and handling the main event loop
func manager(killChan chan bool, exitChan chan int) {
	// Loop and handle events
	for {
		select {
		// Stop the application
		case <-killChan:
			exitChan <- 0
			break
		}
	}
}
