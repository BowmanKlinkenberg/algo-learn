package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Key  string `mapstructure:"ALGO_KEY"`
	Url string `mapstructure:"ALGO_URL"`
}

// loadConfig reads configuration from file or environment variables.
func loadConfig() (config Config, err error) {
	viper.AutomaticEnv()
	err = viper.BindEnv("ALGO_KEY")
	if err != nil {
		return Config{}, err
	}
	viper.SetDefault("ALGO_KEY", "")
	err = viper.BindEnv("ALGO_URL")
	if err != nil {
		return Config{}, err
	}
	viper.SetDefault("ALGO_URL", "https://testnet-algorand.api.purestake.io/ps2/")
	err = viper.UnmarshalExact(&config)
	if err != nil {
		log.Fatal("unable to decode into struct, %v", err)
	}
	return config, nil
}

// newClient creates a new Algod client using standard settings
func newClient(config Config) (client* algod.Client) {

	commonClient, err := common.MakeClient(config.Url, "X-API-Key", config.Key)
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
func getTxParams(client *algod.Client) (s string){
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

}
