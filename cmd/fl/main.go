package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/weaveworks-experiments/fl/internal/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()

	app := cmd.NewApp()
	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Fatalf("failed executing root command: %s", err)
	}
}
