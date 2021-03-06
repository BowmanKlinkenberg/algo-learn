package main

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ApiToken       string `mapstructure:"ALGO_API_TOKEN"`
	Url            string `mapstructure:"ALGO_URL"`
	AcctAddress    string `mapstructure:"ALGO_ADDRESS"`
	AcctPassphrase string `mapstructure:"ALGO_PASSPHRASE"`
}

func (config Config) pubAddress() (address types.Address) {
	privateKey, _ := mnemonic.ToPrivateKey(config.AcctPassphrase)
	var myAddress types.Address
	publicKey := privateKey.Public()
	cpk := publicKey.(ed25519.PublicKey)
	copy(myAddress[:], cpk[:])
	return myAddress
}

func (config Config) priKey() (priKey ed25519.PrivateKey) {
	privateKey, _ := mnemonic.ToPrivateKey(config.AcctPassphrase)
	return privateKey
}

// loadConfig reads configuration from file or environment variables.
func loadConfig() (config Config, err error) {
	viper.AutomaticEnv()
	err = viper.BindEnv("ALGO_API_TOKEN")
	if err != nil {
		return Config{}, err
	}
	//viper.SetDefault("ALGO_API_TOKEN", "")
	err = viper.BindEnv("ALGO_URL")
	if err != nil {
		return Config{}, err
	}
	viper.SetDefault("ALGO_URL", "https://testnet-algorand.api.purestake.io/ps2/")
	err = viper.BindEnv("ALGO_PASSPHRASE")
	if err != nil {
		return Config{}, err
	}
	err = viper.BindEnv("ALGO_ADDRESS")
	if err != nil {
		return Config{}, err
	}

	err = viper.UnmarshalExact(&config)
	if err != nil {
		log.Fatal("unable to decode into struct, %v", err)
	}
	return config, nil
}

// newClient creates a new Algod client using standard settings
func newClient(config Config) (client *algod.Client) {

	commonClient, err := common.MakeClient(config.Url, "X-API-key", config.ApiToken)
	if err != nil {
		fmt.Printf("Issue with creating algod client: %s\n", err)
		return
	}
	// TODO: add checks for catchup-time and validation against block explorer
	return (*algod.Client)(commonClient)
}

// getStatus returns a JSON-formatted status for the client connected node
func getStatus(client *algod.Client) (s string) {
	status, err := client.Status().Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting status: %s\n", err)
		return
	}
	statusJSON, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		fmt.Printf("Can not marshall status data: %s\n", err)
	}
	return string(statusJSON)
}

// getTxParams returns JSON-formatted txparams
func getTxParams(client *algod.Client) (s string) {
	txParams, err := client.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("Error Algorand suggested parameters: %s\n", err)
		return
	}
	JSON, err := json.MarshalIndent(txParams, "", "\t")
	if err != nil {
		fmt.Printf("Can not marshall suggested parameters data: %s\n", err)
	}
	return string(JSON)
}

// prettyPrint prints Go structs
func prettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}
