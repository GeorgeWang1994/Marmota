package main

import (
	"github.com/spf13/cobra"
	"log"
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alarm",
		Short: "alarm daemon",
		Run: func(cmd *cobra.Command, args []string) {
			err := initApp()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}
	return cmd
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
