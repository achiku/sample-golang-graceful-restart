package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	gracedown "github.com/shogo82148/go-gracedown"
)

type counter struct {
	m     *sync.Mutex
	count int
}

func (c *counter) Get() int {
	c.m.Lock()
	defer c.m.Unlock()

	c.count = c.count + 1
	return c.count
}

func newCounter() *counter {
	return &counter{
		m:     &sync.Mutex{},
		count: 0,
	}
}

func createCounterWaitHandler(c *counter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		i := c.Get()
		log.Printf("start # %d", i)
		time.Sleep(10 * time.Second)
		fmt.Fprintf(w, "hello")
		log.Printf("done # %d", i)
	}
}

func waitHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	fmt.Fprintf(w, "hello")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/wait", waitHandler)
	mux.HandleFunc("/count/wait", createCounterWaitHandler(newCounter()))

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
