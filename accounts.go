package main

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

// newStandaloneAccount creates a standalone account and prints the address and mnemonic
func newStandaloneAccount() {
	account := crypto.GenerateAccount()
	passphrase, err := mnemonic.FromPrivateKey(account.PrivateKey)

	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
	} else {
		fmt.Printf("My address: %s\n", account.Address)
		fmt.Printf("My passphrase: %s\n", passphrase)
	}
}

// acctInfo returns information about the configured account
func acctInfo(client *algod.Client, config Config) {
	fmt.Printf("My address: %s\n", config.pubAddress().String())

	accountInfo, err := client.AccountInformation(config.pubAddress().String()).Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting account info: %s\n", err)
		return
	}

	// NOTE: go sdk has a type for MicroAlgos including conversions to Algo
	var microBalance types.MicroAlgos = types.MicroAlgos(accountInfo.Amount)
	var algoBalance float64 = microBalance.ToAlgos()

	fmt.Printf("Account balance: %d microAlgos\n", microBalance)
	fmt.Printf("Account balance: %f Algos\n", algoBalance)
}

// sendFaucetPayment sends an example transaction from the configured address to the Algo faucet address
func sendFaucetPayment(amount uint64, client *algod.Client, config Config) {
	txParams, err := client.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting suggested tx params: %s\n", err)
		return
	}
	// comment out the next two (2) lines to use suggested fees
	txParams.FlatFee = true
	txParams.Fee = 1000

	fromAddr := config.pubAddress().String()
	toAddr := "GD64YIY3TWGDMCNPP553DZPPR6LDUSFQOIJVFDPPXWEG3FVOJCCDBBHU5A"
	var minFee uint64 = 1000
	note := []byte("Hello World")
	genID := txParams.GenesisID
	genHash := txParams.GenesisHash
	firstValidRound := uint64(txParams.FirstRoundValid)
	lastValidRound := uint64(txParams.LastRoundValid)

	txn, err := transaction.MakePaymentTxnWithFlatFee(fromAddr, toAddr, minFee, amount, firstValidRound, lastValidRound, note, "", genID, genHash)
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
		return
	}
	fmt.Println("Transaction id created")

	txID, signedTxn, err := crypto.SignTransaction(config.priKey(), txn)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %s\n", err)
		return
	}
	fmt.Printf("Signed txid: %s\n", txID)

	sendResponse, err := client.SendRawTransaction(signedTxn).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to send transaction: %s\n", err)
		return
	}
	fmt.Printf("Submitted transaction %s\n", sendResponse)
}
