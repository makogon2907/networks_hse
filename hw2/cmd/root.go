package cmd

import (
	"fmt"
	mtu "mtu/src"
	"os"

	"github.com/spf13/cobra"
)

var hostname string

var rootCmd = &cobra.Command{
	Use:   "mtu checker",
	Short: "",
	Long:  "",
	Run:   func(cmd *cobra.Command, args []string) { mtu.FindMTUWith(hostname) },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().StringVar(&hostname, "host", "", "host to check mtu with (REQUIRED)")
	rootCmd.MarkPersistentFlagRequired("host")
}
