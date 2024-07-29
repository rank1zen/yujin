package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func Execute(ctx context.Context) int {
	rootCmd := &cobra.Command{
		Use:   "yujin",
		Short: "Minimalist, Opinionated League of Legends.",
	}

	rootCmd.AddCommand(UiCmd(ctx))

	err := rootCmd.Execute()
	if err != nil {
		return 1
	}

	return 0
}
