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
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// cleanCmd represents the cleanCmd command
var installTerraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Used to install terraform to local machine",
	Long:  `Used to install terraform to local machine`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Running Pulumi.Run()...")
		commands := make(map[string][]string, 7)

		commands["install-prerequisites"] = []string{"sudo", "apt-get", "update", "&&", "sudo", "apt-get", "install", "-y", "gnupg", "software-properties-common"}
		commands["add-hashicorp-key"] = []string{"wget", "-O-", "https://apt.releases.hashicorp.com/gpg", "|", "gpg", "--dearmor", "|", "sudo", "tee", "/usr/share/keyrings/hashicorp-archive-keyring.gpg"}
		commands["verify-hashicorp-key"] = []string{"gpg", "--no-default-keyring", "--keyring", "/usr/share/keyrings/hashicorp-archive-keyring.gpg", "--fingerprint"}
		commands["add-hashicorp-ubuntu-repository"] = []string{"echo", "\"deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main\"", "|", "sudo", "tee", "/etc/apt/sources.list.d/hashicorp.list"}
		commands["update-ubuntu-repository-cache"] = []string{"sudo", "apt", "update"}
		commands["ubuntu-install-terraform"] = []string{"sudo", "apt-get", "install", "-y", "terraform"}
		commands["verify-terraform-installation"] = []string{"terraform", "-help"}
		commands["verify-terraform-installation-2"] = []string{"terraform", "-help", "plan"}
		commands["install-bash-completions"] = []string{"touch", "~/.bashrc", "&&", "terraform", "-install-autocomplete"}
		commands["install-zsh-completions"] = []string{"touch", "~/.zshrc", "&&", "terraform", "-install-autocomplete"}

		for commandName, command := range commands {
			log.Infof("Running Command for %s", commandName)
			output, err := u.ExecuteCommand(command)

			if err != nil {
				log.Errorf("Error running shell command %s. %s", commandName, err.Error())
			}

			log.Infof("Output:\n%s", output)
		}
	},
}

func init() {
	installCmd.AddCommand(installTerraformCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
