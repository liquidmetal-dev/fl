package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/weaveworks-experiments/fl/internal/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()

	cobra.OnInitialize(initConfig)

	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Fatal("failed executing root command")
	}
}

func initConfig() {
	viper.SetEnvPrefix("FL")
	viper.AutomaticEnv()
}
