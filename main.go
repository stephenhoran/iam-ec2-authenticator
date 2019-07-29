package main

import (
	log "github.com/sirupsen/logrus"
	"iam-ec2-authenticator/cmd"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	cmd.Execute()
}
