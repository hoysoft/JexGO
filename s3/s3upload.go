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
	fmt.Println("start s3upload:",filename)
	b,err:=s.Bucket(s3client)
	if err != nil {
		return err
	}

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

func (s *S3Config)Bucket(client *s3.S3) (*s3.Bucket,error){
	b := client.Bucket(s.BucketName)
	resp, err := client.ListBuckets()
	if err != nil {
		return b,err
	}
	if bucketOfName(resp.Buckets,s.BucketName)==nil{
		err=b.PutBucket(s3.ACL(s.Acl))
		if err!=nil{
			return b,err
		}
	}
	return b,nil
}

func bucketOfName(buckets []s3.Bucket,name string) *s3.Bucket {
	for _,b:=range buckets{
		if b.Name==name {
			return &b
		}
	}
	return nil
}