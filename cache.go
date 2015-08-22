package main

import (
	"bytes"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getLatestVersion() string {
	if s3Connection == nil {
		return ""
	}

	out, err := s3Connection.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(config.Bucket),
		Key:    aws.String("latest"),
	})

	if err != nil {
		log.Printf("Unable to load latest version: %s", err)
		return ""
	}
	defer out.Body.Close()

	buf := bytes.NewBuffer([]byte{})
	io.Copy(buf, out.Body)

	return buf.String()
}

func getFromCache(version string) ([]byte, bool) {
	if s3Connection == nil {
		return []byte{}, false
	}

	out, err := s3Connection.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(config.Bucket),
		Key:    aws.String(version),
	})

	if err != nil {
		log.Printf("Could not fetch from store: %s", err)
		return []byte{}, false
	}
	defer out.Body.Close()

	buf := bytes.NewBuffer([]byte{})
	io.Copy(buf, out.Body)

	return buf.Bytes(), true
}

func saveToCache(version string, data []byte) bool {
	if s3Connection == nil {
		return false
	}

	buf := bytes.NewReader(data)
	_, err := s3Connection.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(config.Bucket),
		Key:    aws.String(version),
		Body:   buf,
	})

	if err != nil {
		log.Printf("Unable to save to version: %s", err)
		return false
	}

	_, err = s3Connection.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(config.Bucket),
		Key:    aws.String("latest"),
		Body:   bytes.NewReader([]byte(version)),
	})

	if err != nil {
		log.Printf("Unable to put version id: %s", err)
		return false
	}

	return true
}
