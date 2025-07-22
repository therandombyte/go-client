package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewIVCommand() *cobra.Command {
	fmt.Println("Init here")

	cmd := &cobra.Command {

	}
	return cmd

}
