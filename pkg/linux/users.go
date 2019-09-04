package linux

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/iam"
	"iam-ec2-authenticator/pkg/authiam"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

const PasswdFile = "/etc/passwd"
const DefaultShell = "/bin/bash"

// EtcPasswdEntry contains a single line entry representing a parse user from /etc/passwd in linux. An example line
// form this file looks like:
//
// root:*:0:0:System Administrator:/var/root:/bin/sh
type EtcPasswdEntry struct {
	username string
	password string
	uid      int
	gid      int
	geckos   string
	homedir  string
	shell    string
}

// GetLinuxUsers in the entry point function for the linux package. GetLinuxUsers will return a list of all of the users
// if the file contents of f. This must be in /etc/passwd format.
// It returns a slice of EtcPasswdEntry or and error if the file cannot be opened or if it is not is the correct format.
func GetLinuxUsers(file string) ([]EtcPasswdEntry, error) {
	var EtcPasswdEntries []EtcPasswdEntry

	passwd, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(passwd)), "\n")
	for index, l := range lines {
		// Checking to make sure this is not an empty line or a commented line. If so we skip it.
		if len(l) == 0 || strings.HasPrefix(l, "#") {
			continue
		}

		formatPasswd := strings.Split(l, ":")
		if len(formatPasswd) != 7 {
			return nil, fmt.Errorf("error in line number: %v. Verify this content is correct", index+1)
		}

		uid, err := strconv.Atoi(formatPasswd[2])
		if err != nil {
			return nil, fmt.Errorf("uid field badly formatted at line: %v", index+1)
		}

		gid, err := strconv.Atoi(formatPasswd[3])
		if err != nil {
			return nil, fmt.Errorf("gid field badly formatted at line: %v", index+1)
		}

		EtcPasswdEntries = append(EtcPasswdEntries, EtcPasswdEntry{
			username: formatPasswd[0],
			password: formatPasswd[1],
			uid:      uid,
			gid:      gid,
			geckos:   formatPasswd[4],
			homedir:  formatPasswd[5],
			shell:    formatPasswd[6],
		})
	}

	return EtcPasswdEntries, nil
}

//UserTrueUp determines if any work needs to be done depending on the the results from IAM.
// It accepts a slice of EtcPasswdEntries as well as the IAM Group Output from the AWS SDK Package.
func UserTrueUp(linuxUsers []EtcPasswdEntry, iamUsers iam.GetGroupOutput, svc *iam.IAM) {
	var newUser, existingUser []string

	for _, iamUser := range iamUsers.Users {
		foundUser := 0
		for _, linuxUser := range linuxUsers {
			// User is found. Add it to slice of managed users.
			if linuxUser.username == *iamUser.UserName {
				foundUser = 1
				existingUser = append(existingUser, *iamUser.UserName)
			}
		}
		// User not found. Add it to slice to be added.
		if foundUser == 0 {
			newUser = append(newUser, *iamUser.UserName)
		}
	}
	//AddUserLinux(newUser)
	ValidateSSHKey(newUser, svc)
}

func AddUserLinux(newUsers []string) {
	for _, user := range newUsers {
		err := exec.Command("/usr/sbin/useradd", user, "-m", "--shell", DefaultShell).Run()
		if err != nil {
			panic(err)
		}
	}
}

func ValidateSSHKey(users []string, svc *iam.IAM) {
	for _, user := range users {
		authiam.GetUserSSHKey(svc, user)
	}

}
