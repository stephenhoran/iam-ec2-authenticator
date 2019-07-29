package linux

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"testing"
)

const mockFunctionalPasswdFile = `
# commenting test

# another comment with some white spaces above and below

nobody:*:-2:-2:Unprivileged User:/var/empty:/usr/bin/false
root:*:0:0:System Administrator:/var/root:/bin/sh
daemon:*:1:1:System Services:/var/root:/usr/bin/false
# adding a trailing white space

`

const mockErrorLongUserPasswdFile = `
# commenting test

# another comment with some white spaces above and below

nobody:*:-2:-2:Unprivileged User:/var/empty:/usr/bin/false
root:*:0:0:System Administrator:/var/root:/bin/sh
daemon:*:1:1:System Services:/var/root:/usr/bin/false:failure
# adding a trailing white space

`

const mockErrorFailedUIDPasswdFile = `
nobody:*:-2:-2:Unprivileged User:/var/empty:/usr/bin/false
root:*:abc:0:System Administrator:/var/root:/bin/sh
daemon:*:1:1:System Services:/var/root:/usr/bin/false:failure
`

const mockErrorFailedGIDPasswdFile = `
nobody:*:-2:-2:Unprivileged User:/var/empty:/usr/bin/false
root:*:0:abc:System Administrator:/var/root:/bin/sh
daemon:*:1:1:System Services:/var/root:/usr/bin/false:failure
`

// TestGetLinuxUsers creates a temp /etc/passwd file in which GetLinuxUsers will need to validate it can read and parse
// a properly formatted /etc/passwd file.
func TestGetLinuxUsers(t *testing.T) {
	// table tests as follows:
	// testFile: The mock test file you would like to use.
	// pass: If this test should not return an error or not.
	// length: The length of the slice returned by GetLinuxUsers
	tables := []struct {
		testFile string
		pass     bool
		length   int
	}{
		{
			mockFunctionalPasswdFile,
			true,
			3,
		},
		{
			mockErrorLongUserPasswdFile,
			false,
			0,
		},
		{
			mockErrorFailedUIDPasswdFile,
			false,
			0,
		},
		{
			mockErrorFailedGIDPasswdFile,
			false,
			0,
		},
	}

	for testIndex, test := range tables {
		// Creating a temp /etc/passwd file from mockPasswdFile
		tempDir, err := ioutil.TempDir("", "etc")
		if err != nil {
			t.Logf("Unable to create temp dir: %s", err)
			t.Fail()
		}

		tempPwFile := path.Join(tempDir, "passwd")

		err = ioutil.WriteFile(tempPwFile, []byte(test.testFile), 0644)
		if err != nil {
			t.Fatalf("Unable to write mock passwd file: %s", err)
		}

		users, err := GetLinuxUsers(tempPwFile)
		if err != nil && test.pass == true {
			t.Fatal(err)
		}

		if len(users) != test.length {
			t.Errorf("Test %v failed. Expected length: %v, Received: %v", testIndex, test.length, len(users))
		}

		err = os.Remove(tempPwFile)
		if err != nil {
			t.Log("unable to clean up temp file")
		}
	}
}

// TestGetLinuxUsersCannotOpenFile validates that GetLinuxUsers will error in the event it cannot find or access file
// that is provided to the the function.
func TestGetLinuxUsersCannotOpenFile(t *testing.T) {
	// Create a temp directory etc and then randomly generate a filename and attempt to read it.
	tempDir, err := ioutil.TempDir("", "etc")
	if err != nil {
		t.Logf("Unable to create temp dir: %s", err)
		t.Fail()
	}
	random := []rune("abcdefg1234567")
	r := make([]rune, 10)
	for i := range r {
		r[i] = random[rand.Intn(len(random))]
	}

	tempPwFile := path.Join(tempDir, string(r))

	_, err = GetLinuxUsers(tempPwFile)
	if err == nil {
		t.Fatalf("accessed a file that should not exist: %v", tempPwFile)
	}
}
