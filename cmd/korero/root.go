package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:                   "korero [options]",
		Long:                  "Send, receive and manage messages on multiple platforms",
		SilenceUsage:          true,
		SilenceErrors:         true,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
)

func Execute() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(125)
	}
	os.Exit(0)
}
