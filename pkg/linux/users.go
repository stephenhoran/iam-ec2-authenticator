package linux

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const passwdFile = "/etc/passwd"

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
