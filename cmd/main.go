package main

import (
	"context"
	"os"

	"compose-gen/internal/interface/cli"
)

func main() {
	ctx := context.Background()
	
	cliApp := cli.NewCLI()
	rootCmd := cliApp.RootCommand()
	
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}