package main

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
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
	privateKey, err := mnemonic.ToPrivateKey(config.AcctPassphrase)
	if err != nil {
		fmt.Printf("Issue with mnemonic conversion: %s\n", err)
		return
	}

	var myAddress types.Address
	publicKey := privateKey.Public()
	cpk := publicKey.(ed25519.PublicKey)
	copy(myAddress[:], cpk[:])
	fmt.Printf("My address: %s\n", myAddress.String())

	accountInfo, err := client.AccountInformation(myAddress.String()).Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting account info: %s\n", err)
		return
	}

	var microBalance types.MicroAlgos = types.MicroAlgos(accountInfo.Amount)
	var algoBalance float64 = microBalance.ToAlgos()

	fmt.Printf("Account balance: %d microAlgos\n", microBalance)
	fmt.Printf("Account balance: %f Algos\n", algoBalance)
	}