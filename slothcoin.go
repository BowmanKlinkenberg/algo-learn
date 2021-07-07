package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
)

// printAssetHolding utility to print asset holding for account
func printAssetHolding(assetID uint64, account string, client *algod.Client) {

	act, err := client.AccountInformation(account).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to get account information: %s\n", err)
		return
	}
	for _, assetholding := range act.Assets {
		if assetID == assetholding.AssetId {
			prettyPrint(assetholding)
			break
		}
	}
}

// printCreatedAsset utility to print created assert for account
func printCreatedAsset(assetID uint64, account string, client *algod.Client) {

	act, err := client.AccountInformation(account).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to get account information: %s\n", err)
		return
	}
	for _, asset := range act.CreatedAssets {
		if assetID == asset.Index {
			prettyPrint(asset)
			break
		}
	}
}

func newAsset(client *algod.Client, config Config) (err error) {
	// CREATE ASSET
	txParams, err := client.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting suggested tx params: %s\n", err)
		return
	}
	creator := config.pubAddress().String()
	assetName := "slothcoin"
	unitName := "SLOTH"
	assetURL := ""
	assetMetadataHash := "hookhands"
	defaultFrozen := false
	decimals := uint32(10)
	totalIssuance := uint64(1000)
	manager := config.pubAddress().String()
	reserve := config.pubAddress().String()
	freeze := config.pubAddress().String()
	clawback := config.pubAddress().String()
	note := []byte(nil)

	txn, err := future.MakeAssetCreateTxn(creator, note, txParams, totalIssuance, decimals, defaultFrozen, manager, reserve, freeze, clawback, unitName, assetName, assetURL, assetMetadataHash)
	if err != nil {
		fmt.Printf("Failed to make asset: %s\n", err)
		return
	}
	fmt.Printf("Asset created AssetName: %s\n", txn.AssetConfigTxnFields.AssetParams.AssetName)
	// sign the transaction
	txid, stx, err := crypto.SignTransaction(config.priKey(), txn)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %s\n", err)
		return
	}
	fmt.Printf("Transaction ID: %s\n", txid)
	// Broadcast the transaction to the network
	sendResponse, err := client.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to send transaction: %s\n", err)
		return
	}
	fmt.Printf("Submitted transaction %s\n", sendResponse)
	// Wait for confirmation
	confirmedTxn, err := waitForConfirmation(txid, client, 5)
	if err != nil {
		fmt.Printf("Error wating for confirmation on txID: %s\n", txid)
		return
	}
	// print tx info
	txnJSON, err := json.MarshalIndent(confirmedTxn.Transaction.Txn, "", "\t")
	if err != nil {
		fmt.Printf("Can not marshall txn data: %s\n", err)
	}
	fmt.Printf("Transaction information: %s\n", txnJSON)

	fmt.Printf("Decoded note: %s\n", string(confirmedTxn.Transaction.Txn.Note))

	// Retrieve asset ID by grabbing the max asset ID
	// from the creator account's holdings.
	act, err := client.AccountInformation(config.pubAddress().String()).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to get account information: %s\n", err)
		return
	}

	assetID := uint64(0)
	//  find newest (highest) asset for this account
	for _, asset := range act.CreatedAssets {
		if asset.Index > assetID {
			assetID = asset.Index
		}
	}

	// print created asset and asset holding info for this asset
	fmt.Printf("Asset ID: %d\n", assetID)
	printCreatedAsset(assetID, config.pubAddress().String(), client)
	printAssetHolding(assetID, config.pubAddress().String(), client)
	return nil
}
