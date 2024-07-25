package cmd

import (
	"context"
	"net/http"

	"github.com/spf13/cobra"
)

func Execute(ctx context.Context) int {
	rootCmd := &cobra.Command{
		Use:   "yujin",
		Short: "Minimalist, Opinionated League of Legends.",
	}

	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	err := rootCmd.Execute()
	if err != nil {
		return 1
	}

	return 0
}
