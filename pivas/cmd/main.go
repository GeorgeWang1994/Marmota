package main

import (
	"github.com/spf13/cobra"
	"log"
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pivas",
		Short: "pivas daemon",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return cmd
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}

