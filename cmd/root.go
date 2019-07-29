package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"iam-ec2-authenticator/pkg/linux"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "iam-ec2-authenticator",
	Short: "iam-ec2-authenticator adds users from IAM to EC2 servers from a group",
	Long: `iam-ec2-authenticator creates and manages user ssh-keys from IAM. The service will query a specified IAM group
for changes and apply those changes to the server this runs on. No IAM credentials will be added to the server, only ssh-key
and user creation takes place via this tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		linux.GetLinuxUsers()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
