package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/liquidmetal-dev/fl/internal/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()

	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Fatalf("failed executing root command: %s", err)
	}
}
