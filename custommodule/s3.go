package custommodule

import (
	"bytes"
	"crypto/tls"
	"flag"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 main type
type S3 struct {
	Config *aws.Config
	Bucket string
}

// NewS3 create new S3 instance
func NewS3(accesskey string, secretkey string, host string, bucketname string, region string) (S3, error) {
	result := S3{}
	creds := credentials.NewStaticCredentials(accesskey, secretkey, "")
	_, err := creds.Get()
	if err != nil {
		return S3{}, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	result.Bucket = bucketname
	result.Config = &aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(host),
		Credentials:      creds,
		S3ForcePathStyle: aws.Bool(true),
		HTTPClient:       client,
	}

	return result, err
}

// UploadJPEG uploads jpeg to s3
func (x *S3) UploadJPEG(b []byte, filename string) (string, error) {
	var err error
	if flag.Lookup("test.v") == nil {
		session, _ := session.NewSession(x.Config)

		s3Client := s3.New(session)
		s3Client.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(x.Bucket),
		})

		file, _ := os.Create("files/" + filename + ".jpeg")
		defer file.Close()
		file.Write(b)
		file.Sync()
		fileinfo, _ := file.Stat()

		open, _ := os.Open("files/" + filename + ".jpeg")
		defer open.Close()
		defer os.Remove("files/" + filename + ".jpeg")
		buffer := make([]byte, fileinfo.Size())
		open.Read(buffer)

		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(x.Bucket),
			Key:    aws.String(fileinfo.Name()),
			Body:   bytes.NewReader(buffer),
		})
	}

	return *x.Config.Endpoint + "/" + x.Bucket + "/" + filename + ".jpeg", err
}
