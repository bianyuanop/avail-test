package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/rollkit/avail-da"
	"github.com/rollkit/go-da"
)

func main() {
	conf := avail.Config{
		AppID: 0,
		LcURL: "http://localhost:7000/v2",
	}

	ctx := context.Background()

	cli := avail.NewAvailDA(conf, ctx)

	blob := make([]byte, 100)

	ids, err := cli.Submit(ctx, []da.Blob{blob}, -1, nil)
	if err != nil {
		fmt.Printf("unable to submit data, reason: %s\n", err.Error())
		return
	}

	if len(ids) != 1 {
		fmt.Printf("length of submitted data doesnt' match, len: %d\n", len(ids))
		return
	}

	fmt.Println(hex.EncodeToString(ids[0]))
}
