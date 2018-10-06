// Package util implements utility methods
package util

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	// The region used to test
	AWSREGION = os.Getenv("AWSREGION")
	// The bucket to use
	AWSBUCKET = os.Getenv("AWSBUCKET")
	// The directory containing testdata
	TESTDATADIR = os.Getenv("TESTDATADIR")
)

// The test suite for the S3 code
type S3TestSuite struct {
	suite.Suite
	awsSession *session.Session
	testDir    string
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *S3TestSuite) SetupSuite() {
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
	suite.testDir = TESTDATADIR
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests are run.
func (suite *S3TestSuite) TearDownSuite() {
	os.RemoveAll(filepath.Join(suite.testDir))
	s3Session := s3.New(suite.awsSession)
	s3Session.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(AWSBUCKET), Key: aws.String("helloworld.txt")})
}

func (suite *S3TestSuite) TestUploadNotExistingFile() {
	err := UploadFile(suite.awsSession, filepath.Join(suite.testDir, "downloads"), "notexistingfile.txt", AWSBUCKET)
	assert.Error(suite.T(), err)
}

func (suite *S3TestSuite) TestUploadExistingFile() {
	err := UploadFile(suite.awsSession, suite.testDir, "helloworld.txt", AWSBUCKET)
	assert.NoError(suite.T(), err)

	s3Session := s3.New(suite.awsSession)
	output, err := s3Session.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(AWSBUCKET)})
	fmt.Println("---")
	for _, val := range output.Contents {
		fmt.Println(val.Key)
	}
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "helloworld.txt", *output.Contents[0].Key)
}

func (suite *S3TestSuite) TestDownloadNotExistingFile() {
	err := DownloadFile(suite.awsSession, filepath.Join(suite.testDir, "downloads"), "notexistingfile.txt", AWSBUCKET)
	assert.Error(suite.T(), err)
}

func (suite *S3TestSuite) TestDownloadExistingFile() {
	suite.TestUploadExistingFile()
	err := DownloadFile(suite.awsSession, filepath.Join(suite.testDir, "downloads"), "helloworld.txt", AWSBUCKET)
	assert.NoError(suite.T(), err)
}

func (suite *S3TestSuite) TestCopyNotExistingFile() {
	err := CopyFile(suite.awsSession, "notexistingfile.txt", AWSBUCKET)
	assert.Error(suite.T(), err)
}

func (suite *S3TestSuite) TestCopyExistingFile() {
	suite.TestUploadExistingFile()
	err := CopyFile(suite.awsSession, "helloworld.txt", AWSBUCKET)
	assert.NoError(suite.T(), err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestS3TestSuite(t *testing.T) {
	suite.Run(t, new(S3TestSuite))
}
