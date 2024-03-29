package cmd

/*
Copyright © 2022  <>

Licensed under the HLT License, Version 0.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"fmt"
	"github.com/spf13/cobra"
)

// whitelistCmd represents the whitelistCmd command
var whitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "A wrapper for all whitelisting subcommands",
	Long:  `A wrapper for all whitelisting subcommands`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use `whitelist help` for help")
	},
}

func init() {
	rootCmd.AddCommand(whitelistCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whitelistCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whitelistCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
