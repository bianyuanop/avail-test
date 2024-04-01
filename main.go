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

	var wg sync.WaitGroup
	wg.Add(3)

	// Submitting simultaneously
	fmt.Println("submitting data1")

	// go func() {
	// 	defer wg.Done()
	// 	size1 := 0.5 * 1024 * 1024
	// 	err = submitWithSize(ctx, cli, int(size1))
	// 	if err != nil {
	// 		fmt.Printf("unable to submit with size %d, err: %s\n", int(size1), err)
	// 	}
	// }()

	// go func() {
	// 	defer wg.Done()
	// 	size1 := 1 * 1024 * 1024
	// 	err = submitWithSize(ctx, cli, int(size1))
	// 	if err != nil {
	// 		fmt.Printf("unable to submit with size %d, err: %s\n", int(size1), err)
	// 	}
	// }()

	// time.Sleep(5 * time.Second)

	go func() {
		defer wg.Done()
		size1 := 0.15 * 1024 * 1024
		err = submitWithSize(ctx, cli, int(size1))
		if err != nil {
			fmt.Printf("unable to submit with size %d, err: %s\n", int(size1), err)
		}
	}()
	time.Sleep(3 * time.Second)

	go func() {
		defer wg.Done()
		size1 := 1 * 1024 * 1024
		err = submitWithSize(ctx, cli, int(size1))
		if err != nil {
			fmt.Printf("unable to submit with size %d, err: %s\n", int(size1), err)
		}
	}()

	time.Sleep(3 * time.Second)
	go func() {
		defer wg.Done()
		size1 := 0.1 * 1024 * 1024
		err = submitWithSize(ctx, cli, int(size1))
		if err != nil {
			fmt.Printf("unable to submit with size %d, err: %s\n", int(size1), err)
		}
	}()
	wg.Wait()
	// data := make([]byte, int(size))
	// _, err = rand.Read(data)
	// if err != nil {
	// 	fmt.Println("unable to gen random data1")
	// }
	// ctx, cancle := context.WithTimeout(ctx, 70*time.Second)
	// defer cancle()
	// h, err := cli.SubmitData(ctx, data)
	// if err != nil {
	// 	fmt.Printf("unable to submit data1, reason: %s\n", err.Error())
	// }
	// fmt.Printf("hash of tx1: %s\n", h.Hex())
}

func submitWithSize(ctx context.Context, cli *client.AvailClient, size int) error {
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		return fmt.Errorf("unable to gen random data")
	}
	ctx, cancle := context.WithTimeout(ctx, 70*time.Second)
	defer cancle()
	h, err := cli.SubmitData(ctx, data)
	if err != nil {
		return fmt.Errorf("unable to submit data, reason: %s\n", err.Error())
	}
	fmt.Printf("hash of tx: %s\n", h.Hex())

	return nil
}
