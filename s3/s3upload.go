package s3

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"os"
	"path/filepath"
	"github.com/hoysoft/JexGO/reader"
	"fmt"
)

type S3Config struct {
	AccessKey string
    SecretKey string
	BucketName string
	S3Endpoint string
	File      string
	Acl       string
}


func (s *S3Config)Upload(fileDir string,progressFunc reader.ProgressReaderCallbackFunc)error {
	if len(s.Acl)==0{
		s.Acl=string(s3.PublicReadWrite)
	}
	auth, err := aws.GetAuth(s.AccessKey,s.SecretKey)
	s3client := s3.New(auth, aws.Region{Name: "us-east-1", S3Endpoint: s.S3Endpoint })

	filename :=filepath.Join(fileDir,s.File)
	b := s3client.Bucket(s.BucketName)
	f, err := os.Stat(filename)
	if err != nil {
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	progressR := &reader.Reader{
		Reader: file,
		Size:   f.Size(),
		DrawFunc: progressFunc,
	}

	err = b.PutReader(s.File, progressR, f.Size(),"application/octet-stream", s3.ACL(s.Acl))
	//err = b.Put("zoujtw2015-12-16.mkv", file, "content-type", s3.PublicReadWrite)
	if err!=nil{
		return err
	}
	fmt.Println("s3 upload file succeed!!!",file.Name())

	return nil
}
