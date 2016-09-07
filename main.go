package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gracedown "github.com/shogo82148/go-gracedown"
)

func waitHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	fmt.Fprintf(w, "hello")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/wait", waitHandler)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP)
	go func() {
		for {
			s := <-signalChan
			log.Printf("%+v", s)
			if s == syscall.SIGHUP {
				log.Println("received SIGHUP")
				gracedown.Close() // trigger graceful shutdown
			}
		}
	}()
	log.Printf("Server PID: %d", os.Getpid())
	if err := gracedown.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
