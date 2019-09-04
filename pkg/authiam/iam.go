// IAM package is used to lookup users from a specific authiam group.
package authiam

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

// GetIAMUsers determines a which users reside in a given group on IAM. We are using the stub package to allow for testing
// which means a service client will need to be passed.
// It accepts a service client for IAM as well as a group name as a string.
// It returns a slice of sting of users.
func GetIAMUsers(svc iamiface.IAMAPI, groupName string) (*iam.GetGroupOutput, error) {

	users, err := svc.GetGroup(&iam.GetGroupInput{
		GroupName: aws.String(groupName),
		Marker:    nil,
		MaxItems:  nil,
	})
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserSSHKey(svc iamiface.IAMAPI, userName string) (*iam.ListSSHPublicKeysOutput, error) {
	user, err := svc.ListSSHPublicKeys(&iam.ListSSHPublicKeysInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
