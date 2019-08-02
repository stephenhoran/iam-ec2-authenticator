package authiam

import (
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"testing"
)

// MockIAMService used to stab the IAM service response for testing.
type MockIAMService struct {
	iamiface.IAMAPI
	iamdb mockIAMStore
}

// GetGroup is a stub used for testing the iam.GetGroup function from the AWS Go SDK package.
// Returned in the mock interface.
func (m MockIAMService) GetGroup(group *iam.GetGroupInput) (*iam.GetGroupOutput, error) {
	resp := &iam.GetGroupOutput{
		Group: &iam.Group{},
		Users: []*iam.User{},
	}
	for _, mockgroup := range m.iamdb.Groups {
		if mockgroup.GroupName == *group.GroupName {
			resp.Group.GroupName = group.GroupName

			var users []*iam.User
			for i := range mockgroup.Users {
				users = append(users, &iam.User{
					Arn:      &mockgroup.Users[i].Arn,
					UserId:   &mockgroup.Users[i].UserID,
					UserName: &mockgroup.Users[i].Username,
				})
			}

			resp.Users = users
		}
	}

	return resp, nil
}

// MockIAMStore contains the root of an Stubbed IAM store. It allows you to easy create tests.
type mockIAMStore struct {
	Groups []mockIAMGroup
}

// MockIAMGroup is a stub of a single group. Add fields here for additional testing of group.
type mockIAMGroup struct {
	GroupName string
	Users     []mockIAMUser
}

// MockIAMUser is a stub of a single user. Add fields here for additional testing of a group.
type mockIAMUser struct {
	Username string
	UserID   string
	Arn      string
}

func TestGetIAMUsers(t *testing.T) {
	// Our MockIAM Database
	iamdb := mockIAMStore{
		Groups: []mockIAMGroup{
			{
				GroupName: "developers",
				Users: []mockIAMUser{
					{
						Username: "Zed",
						UserID:   "123",
						Arn:      "arn:123",
					},
					{
						Username: "Dead",
						UserID:   "456",
						Arn:      "arn:456",
					},
				},
			},
		},
	}

	cases := []struct {
		Group    int    // The group you wish to test (this is their index position)
		Expected string // pass or fail (expect error)
	}{
		{
			Group:    0,
			Expected: "Pass",
		},
	}

	for _, c := range cases {
		svc := MockIAMService{iamdb: iamdb}
		out, err := GetIAMUsers(svc, iamdb.Groups[c.Group].GroupName)
		if err != nil {
			if c.Expected == "pass" {
				t.Fail()
			}
		}

		if *out.Group.GroupName != iamdb.Groups[c.Group].GroupName {
			t.Fatalf("Wrong group returned, Got: %s, Wanted: %s", *out.Group.GroupName, iamdb.Groups[c.Group].GroupName)
		}
	}

}
