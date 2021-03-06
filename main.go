package main

import (
	"fmt"
	"log"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	algodClient := newClient(config)

	// print status
	fmt.Println("Connected node status:")
	fmt.Println(getStatus(algodClient))

	// suggested transaction params
	fmt.Println("Suggested transaction parameters:")
	fmt.Println(getTxParams(algodClient))

	// creating and printing a new standalone account. leave commented as it will create a new account on every run
	//newStandaloneAccount()

	// print account balance
	acctInfo(algodClient, config)

	// send a test payment to faucet address
	sendFaucetPayment(100000, algodClient, config)

	// create a new slothcoin
	err = newAsset(algodClient, config)
	if err != nil {
		log.Fatal("failed to create new asset", err)
	}
}
