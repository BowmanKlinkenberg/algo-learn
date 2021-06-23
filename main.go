package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"os"
)

const algodAddress = "https://testnet-algorand.api.purestake.io/ps2/"

func main() {
	apiKey := os.Getenv("ALGO_KEY")

	commonClient, err := common.MakeClient(algodAddress, "X-API-Key", apiKey)
	if err != nil {
		fmt.Printf("Issue with creating algod client: %s\n", err)
		return
	}

	algodClient := (*algod.Client)(commonClient)

	//print status
	status, err := algodClient.Status().Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting status: %s\n", err)
		return
	}
	statusJSON, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		fmt.Printf("Can not marshall status data: %s\n", err)
	}
	fmt.Printf("%s\n", statusJSON)
}
