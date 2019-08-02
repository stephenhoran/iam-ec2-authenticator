package cmd

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"iam-ec2-authenticator/pkg/authiam"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "authiam-ec2-authenticator",
	Short: "authiam-ec2-authenticator adds users from IAM to EC2 servers from a group",
	Long: `authiam-ec2-authenticator creates and manages user ssh-keys from IAM. The service will query a specified IAM group
for changes and apply those changes to the server this runs on. No IAM credentials will be added to the server, only ssh-key
and user creation takes place via this tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		svc := iam.New(sess)
		authiam.GetIAMUsers(svc, "developers")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
