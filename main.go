package main

import (
	"fmt"
	"minecraft_searcher/api"
	"minecraft_searcher/api/rest"
	"minecraft_searcher/scanner"
	"os"
	"os/signal"

	"github.com/spf13/viper"
)

func main() {

	viper.AddConfigPath(".")

	viper.SetDefault("threads", 100)
	viper.SetDefault("timeout", 10000)
	viper.SetDefault("input", "minecraft.bin")
	viper.SetDefault("redis.addr", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")

	viper.SafeWriteConfigAs("config.toml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("error while reading config: %s\n", err.Error())
	}

	file, err := os.Open(viper.GetString("input"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := scanner.RunNewScanner(viper.GetInt("timeout"), 0, viper.GetInt("threads"), file)

	//initialize API
	go api.Init()
	go rest.Init(scanner.Telemetry, scanner.LatestErrors)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("\nGoodbye")
}
