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
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/shammishailaj/scli/pkg/schemas"
	"github.com/spf13/cobra"
)

// cleanCmd represents the cleanCmd command
var imapDemoCmd = &cobra.Command{
	Use:   "demo",
	Short: "A demo of the imap function",
	Long:  `A demo of the imap function`,
	Run: func(cmd *cobra.Command, args []string) {
		serviceProviders := make(map[string]*schemas.ImapServer)

		// Server Details from: https://support.microsoft.com/en-us/office/pop-imap-and-smtp-settings-8361e398-8af4-4e97-b147-6c6c4ac95353
		serviceProviders["gmail"] = schemas.NewImapServer("imap.gmail.com", 993, true)
		serviceProviders["gmail"] = schemas.NewImapServer("imap.gmail.com", 993, true)
		serviceProviders["outlook"] = schemas.NewImapServer("outlook.office365.com", 993, true)
		serviceProviders["live"] = schemas.NewImapServer("outlook.office365.com", 993, true)
		serviceProviders["hotmail"] = schemas.NewImapServer("outlook.office365.com", 993, true)
		serviceProviders["microsoft365"] = schemas.NewImapServer("outlook.office365.com", 993, true)
		serviceProviders["ms365"] = schemas.NewImapServer("outlook.office365.com", 993, true)
		serviceProviders["msn"] = schemas.NewImapServer("imap-mail.outlook.com", 993, true)

		serviceProvider, serviceProviderErr := cmd.Flags().GetString("provider")
		if serviceProviderErr != nil {
			u.Log.Fatalf("Illegal or unsupported service provider name. Valid service provider values are: gmail, outlook, live, hotmail, microsoft365 or ms365, msn. %s", serviceProviderErr.Error())
		}

		username, usernameErr := cmd.Flags().GetString("username")
		if usernameErr != nil {
			u.Log.Fatalf("Illegal or empty username. %s", usernameErr.Error())
		}
		password, passwordErr := cmd.Flags().GetString("password")
		if passwordErr != nil {
			u.Log.Fatalf("Illegal or empty password. %s", passwordErr.Error())
		}

		serverAddress := serviceProviders[serviceProvider].ServerPort()

		u.Log.Printf("Connecting to server...")

		// Connect to server
		c, err := client.DialTLS(serverAddress, nil)
		if err != nil {
			u.Log.Fatal(err)
		}
		fmt.Printf(".......................Connected\n")

		// Don't forget to logout
		defer func() {
			logoutErr := c.Logout()
			if logoutErr != nil {
				u.Log.Errorf("Error logging out. %s", logoutErr.Error())
			}
		}()

		// Login
		if err := c.Login(username, password); err != nil {
			u.Log.Fatalf("Failed to login to server: %s with username %s and password %s. %s", serverAddress, username, password, err.Error())
		}
		u.Log.Println("Logged in")

		// List mailboxes
		mailboxes := make(chan *imap.MailboxInfo, 10)
		done := make(chan error, 1)
		go func() {
			done <- c.List("", "*", mailboxes)
		}()

		u.Log.Println("Mailboxes:")
		for m := range mailboxes {
			u.Log.Println("* " + m.Name)
		}

		if err := <-done; err != nil {
			u.Log.Fatal(err)
		}

		// Select INBOX
		mbox, err := c.Select("INBOX", false)
		if err != nil {
			u.Log.Fatal(err)
		}
		u.Log.Println("Flags for INBOX:", mbox.Flags)

		// Get the last 4 messages
		from := uint32(1)
		to := mbox.Messages
		if mbox.Messages > 3 {
			// We're using unsigned integers here, only subtract if the result is > 0
			from = mbox.Messages - 3
		}
		seqset := new(imap.SeqSet)
		seqset.AddRange(from, to)

		messages := make(chan *imap.Message, 10)
		done = make(chan error, 1)
		go func() {
			done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		}()

		u.Log.Println("Last 4 messages:")
		for msg := range messages {
			u.Log.Println("* " + msg.Envelope.Subject)
		}

		if err := <-done; err != nil {
			u.Log.Fatal(err)
		}

		u.Log.Println("Done!")
	},
}

func init() {
	imapCmd.AddCommand(imapDemoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	imapDemoCmd.Flags().StringP("provider", "s", "microsoft365", "IMAP service provider name. Support providers - gmail, outlook, live, hotmail, microsoft365 or ms365, msn")
	imapDemoCmd.Flags().StringP("username", "u", "", "Username to be used while logging-in")
	imapDemoCmd.Flags().StringP("password", "p", "", "Password to be used while logging-in")
}
