package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"iam-ec2-authenticator/pkg/authiam"
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
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		// Lets grab a list of current users from a group in IAM.
		svc := iam.New(sess)
		groups, err := authiam.GetIAMUsers(svc, "developers")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Determine the current users on this linux host.
		linuxUsers, err := linux.GetLinuxUsers(linux.PasswdFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Figure out the work required between these two data point's.
		linux.UserTrueUp(linuxUsers, *groups, svc)

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
