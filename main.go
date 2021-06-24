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

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (c Config, err error) {
	viper.AutomaticEnv()

	err = viper.BindEnv("key", "ALGO_KEY")
	if err != nil {
		return
	}
	err = viper.BindEnv("url", "ALGO_URL")
	if err != nil {
		return
	}
	fmt.Println(viper.AllSettings())

	err = viper.UnmarshalExact(&c)
	if err != nil {
		log.Fatal("unable to decode into struct, %v", err)
	}
	return c, nil
}

func main() {

	//config, err := LoadConfig()
	//if err != nil {
	//	log.Fatal("cannot load config:", err)
	//}
	config := new(Config)
	viper.AutomaticEnv()
	viper.BindEnv("ALGO_KEY")
	viper.SetDefault("ALGO_KEY", "foo")
	viper.BindEnv("ALGO_URL")
	viper.SetDefault("ALGO_URL", "https://testnet-algorand.api.purestake.io/ps2/")
	viper.ReadInConfig()
	fmt.Println(viper.AllSettings())
	err := viper.UnmarshalExact(&config)
	if err != nil {
		log.Fatal("unable to decode into struct, %v", err)
	}


	commonClient, err := common.MakeClient(config.Url, "X-API-Key", config.Key)
	if err != nil {
		fmt.Printf("Issue with creating algod client: %s\n", err)
		return
	}

	algodClient := (*algod.Client)(commonClient)

	// print node status
	status, err := algodClient.Status().Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting status: %s\n", err)
		return
	}
	statusJSON, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		fmt.Printf("Can not marshall status data: %s\n", err)
	}
	fmt.Println("Connected node status:")
	fmt.Printf("%s\n", statusJSON)

	// suggested transaction params
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("Error Algorand suggested parameters: %s\n", err)
		return
	}
	JSON, err := json.MarshalIndent(txParams, "", "\t")
	if err != nil {
		fmt.Printf("Can not marshall suggested parameters data: %s\n", err)
	}
	fmt.Println("Suggested transaction parameters:")
	fmt.Printf("%s\n", JSON)

}
