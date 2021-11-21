package services

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"os"
	"strings"
)

type VideoUpload struct {
	Paths []string
	VideoPath string
	OutputBucket string
	Errors []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (vu * VideoUpload) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error  {

	path := strings.Split(objectPath, os.Getenv("LOCAL_STORAGE_PATH"))

	f, err := os.Open(objectPath)

	if err != nil {
		return err
	}

	defer f.Close()

	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}