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

}
