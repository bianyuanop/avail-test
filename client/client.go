package client

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

const AvailMaxBlobSize = 2 * 1024 * 1024

var ErrSubmissionTimeout = errors.New("da submission timeout")

type AvailClient struct {
	Api         *gsrpc.SubstrateAPI
	AppID       int
	MaxBlobSize int

	kp *signature.KeyringPair
	// this is used for consecutive data submission
	nonce uint32
}

type AvailClientConfig struct {
	ApiURL string `json:"apiURL"`
	AppID  int    `json:"appID"`
	Seed   string `json:"seed"`
}

func NewAvailClient(config *AvailClientConfig) (*AvailClient, error) {
	api, err := gsrpc.NewSubstrateAPI(config.ApiURL)
	if err != nil {
		return nil, err
	}
	keyringPair, err := signature.KeyringPairFromSecret(config.Seed, 42)
	if err != nil {
		panic(fmt.Sprintf("cannot create KeyPair:%v", err))
	}

	return &AvailClient{
		Api:         api,
		AppID:       config.AppID,
		MaxBlobSize: AvailMaxBlobSize,

		kp:    &keyringPair,
		nonce: 0,
	}, nil
}

// The following example shows how submit data blob and track transaction status
func (ac *AvailClient) SubmitData(ctx context.Context, data []byte) (*types.Hash, error) {
	// var configJSON string
	// var config config.Config
	// flag.StringVar(&configJSON, "config", "", "config json file")
	// flag.Parse()

	// if configJSON == "" {
	// 	log.Println("No config file provided. Exiting...")
	// 	os.Exit(0)
	// }

	// err := config.GetConfig(configJSON)
	// if err != nil {
	// 	panic(fmt.Sprintf("cannot get config:%v", err))
	// }

	api := ac.Api
	// , err := gsrpc.NewSubstrateAPI(ac.ApiURL)
	// if err != nil {
	// 	panic(fmt.Sprintf("cannot create api:%v", err))
	// }

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(fmt.Sprintf("cannot get metadata:%v", err))
	}

	// Set data and appID according to need
	subData := hex.EncodeToString(data)
	appID := 0

	// if app id is greater than 0 then it must be created before submitting data
	if ac.AppID != 0 {
		appID = ac.AppID
	}

	newCall, err := types.NewCall(meta, "DataAvailability.submit_data", types.NewBytes([]byte(subData)))
	if err != nil {
		panic(fmt.Sprintf("cannot create new call:%v", err))
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(newCall)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, fmt.Errorf("cannot get block hash:%v", err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return nil, fmt.Errorf("cannot get latest runtime version:%v", err)
	}

	keyringPair := ac.kp

	key, err := types.CreateStorageKey(meta, "System", "Account", keyringPair.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create storage key:%w", err)
	}

	var accountInfo types.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		panic(fmt.Sprintf("cannot get latest storage:%v", err))
	}

	var nonce uint32
	if ac.nonce == 0 {
		nonce = uint32(accountInfo.Nonce)
		ac.nonce = nonce
	} else {
		nonce = ac.nonce + 1
		ac.nonce++
	}
	fmt.Printf("nonce: %d\n", nonce)

	options := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(100),
		AppID:              types.NewUCompactFromUInt(uint64(appID)),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction using Alice's default account
	err = ext.Sign(*keyringPair, options)
	if err != nil {
		return nil, fmt.Errorf("cannot sign:%v", err)
	}

	// Send the extrinsic
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return nil, fmt.Errorf("cannot submit extrinsic:%v", err)
	}

	defer sub.Unsubscribe()
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				fmt.Printf("Txn inside block %v\n", status.AsInBlock.Hex())
			} else if status.IsFinalized {
				fmt.Printf("Txn inside finalized block\n")
				hash := status.AsFinalized
				err := getData(hash, api, subData)
				if err != nil {
					return nil, err
				}
				return &hash, nil
			}
		case <-ctx.Done():
			return nil, ErrSubmissionTimeout
		}
	}
}

// getData extracts data from the block and compares it
func getData(hash types.Hash, api *gsrpc.SubstrateAPI, data string) error {
	block, err := api.RPC.Chain.GetBlock(hash)
	if err != nil {
		return fmt.Errorf("cannot get block by hash:%w", err)
	}
	for _, ext := range block.Block.Extrinsics {
		// these values below are specific indexes only for data submission, differs with each extrinsic
		if ext.Method.CallIndex.SectionIndex == 29 && ext.Method.CallIndex.MethodIndex == 1 {
			arg := ext.Method.Args
			str := string(arg)
			slice := str[2:]
			// fmt.Println("string value", slice)
			// fmt.Println("data", data)
			if slice == data {
				fmt.Println("Data found in block")
			}
		}
	}
	return nil
}
