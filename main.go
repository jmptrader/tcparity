package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mdlayher/goset"
)

// app is the name of the application, as printed in logs
const app = "tcparity"

// servers is a flag which lists comma-separated servers to be used by tcparity
var servers = flag.String("servers", "", "Comma-separated list of servers to be used by tcparity.")

// test is a flag which causes tcparity to start, and exit shortly after
var test = flag.Bool("test", false, "Make tcparity start, and exit shortly after. Used for testing.")

func main() {
	// Set up command line options and logging
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// If test mode, trigger quit shortly after startup
	// Used for CI tests, so that we ensure tcparity starts up and is able to stop gracefully
	if *test {
		go func() {
			fmt.Println(app, ": launched in test mode")
			time.Sleep(5 * time.Second)

			fmt.Println(app, ": test mode triggering graceful shutdown")
			err := syscall.Kill(os.Getpid(), syscall.SIGTERM)
			if err != nil {
				fmt.Println(app, ": failed to invoke graceful shutdown, halting")
				os.Exit(1)
			}
		}()
	}

	// Create a set of servers to balance between
	serverSet := set.New()
	for _, s := range strings.Split(*servers, ",") {
		// Attempt to dial server to check connectivity
		if _, err := net.DialTimeout("tcp", s, time.Duration(5*time.Second)); err != nil {
			fmt.Println(app, ": failed to dial server:", err)
			os.Exit(1)
		}

		// Add server to set
		serverSet.Add(s)
	}

	// Launch manager via goroutine
	killChan := make(chan bool, 1)
	exitChan := make(chan int, 1)
	go manager(serverSet, killChan, exitChan)

	// Gracefully handle termination via UNIX signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)
	for sig := range sigChan {
		// Trigger manager shutdown
		fmt.Println(app, ": caught signal:", sig)
		killChan <- true
		break
	}

	// Force terminate if signaled twice
	go func() {
		for sig := range sigChan {
			fmt.Println(app, ": caught signal:", sig, ", force halting now!")
			os.Exit(1)
		}
	}()

	// Graceful exit
	code := <-exitChan
	fmt.Println(app, ": graceful shutdown complete")
	os.Exit(code)
}
