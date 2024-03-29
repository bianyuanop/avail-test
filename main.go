package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/bianyuanop/avail-test/client"
)

func main() {
	conf := client.AvailClientConfig{
		ApiURL: "wss://goldberg.avail.tools:443/ws",
		Seed:   "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice",
		AppID:  0,
	}

	cli, err := client.NewAvailClient(&conf)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	// Submitting single
	data := make([]byte, 200)
	_, err = rand.Read(data)
	if err != nil {
		fmt.Println("unable to gen random data")
	}
	ctx, cancle := context.WithTimeout(ctx, 70*time.Second)
	defer cancle()
	h, err := cli.SubmitData(ctx, data)
	if err != nil {
		fmt.Printf("unable to submit data, reason: %s\n", err.Error())
	}
	fmt.Printf("hash of tx: %s\n", h.Hex())

	var wg sync.WaitGroup
	wg.Add(2)

	// Submitting simultaneously
	go func() {
		data := make([]byte, 200)
		_, err := rand.Read(data)
		if err != nil {
			fmt.Println("unable to gen random data1")
		}
		ctx, cancle := context.WithTimeout(ctx, 70*time.Second)
		defer cancle()
		h, err := cli.SubmitData(ctx, data)
		if err != nil {
			fmt.Printf("unable to submit data1, reason: %s\n", err.Error())
		}
		fmt.Printf("hash of tx1: %s\n", h.Hex())
		wg.Done()
	}()
	go func() {
		data := make([]byte, 200)
		_, err = rand.Read(data)
		if err != nil {
			fmt.Println("unable to gen random data2")
		}
		ctx, cancle := context.WithTimeout(ctx, 70*time.Second)
		defer cancle()
		h, err := cli.SubmitData(ctx, data)
		if err != nil {
			fmt.Printf("unable to submit data2, reason: %s\n", err.Error())
		}
		fmt.Printf("hash of tx2: %s\n", h.Hex())
		wg.Done()
	}()

	wg.Wait()
}
