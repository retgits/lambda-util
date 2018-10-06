// Package util implements utility methods
package util

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// The test suite for the S3 code
type SSMTestSuite struct {
	suite.Suite
	awsSession *session.Session
	testDir    string
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *SSMTestSuite) SetupSuite() {
	// Create an AWS session
	awsAccessKeyID := os.Getenv("AWSACCESSKEYID")
	awsSecretAccessKey := os.Getenv("AWSSECRETACCESSKEY")
	if len(awsAccessKeyID) == 0 && len(awsSecretAccessKey) == 0 {
		suite.awsSession = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(AWSREGION),
		}))
	} else {
		awsCredentials := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")
		suite.awsSession = session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(AWSREGION),
			Credentials: awsCredentials,
		}))
	}
}

func (suite *SSMTestSuite) TestNotExistingParameter() {
	_, err := GetSSMParameter(suite.awsSession, "MyAwesomeParameter", true)
	assert.Error(suite.T(), err)
}

func (suite *SSMTestSuite) TestExistingParameter() {
	param, err := GetSSMParameter(suite.awsSession, "TestParamEncrypted", true)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), param, "HelloWorld")
	param, err = GetSSMParameter(suite.awsSession, "TestParamNotEncrypted", false)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), param, "HelloWorld")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSSMTestSuite(t *testing.T) {
	suite.Run(t, new(SSMTestSuite))
}
