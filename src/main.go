package main

import (
	"github.com/joleques/go-sqs-poc/src/api"
	"github.com/joleques/go-sqs-poc/src/sqs"
	"log"
	"sync"
)

var wg sync.WaitGroup

func main() {
	log.Println("Start API 1.0")
	wg.Add(2)
	go api.Start()
	go sqs.Receive()
	wg.Wait()
}
