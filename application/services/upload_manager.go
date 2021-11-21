package services

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

func (vu *VideoUpload) loadPaths() error {
	err := filepath.Walk(vu.VideoPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error  {

	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	err := vu.loadPaths()
	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	//multiple process
	for process := 0; process < concurrency; process ++ {
		go vu.uploadWorker(in, returnChannel, uploadClient, ctx)
	}

	go func() {
		for x:=0; x<len(vu.Paths); x++ {
			in <- x
		}
	}()

	//if one of uploads got error or success message, all the others uploads stop immediatly
	for r := range returnChannel {
		if r != "" {
			doneUpload <- r
			break
		}
	}

	return nil
}

func (vu *VideoUpload) uploadWorker(in chan int, returnChan chan string, uploadClient *storage.Client, ctx context.Context)  {
	for x := range in {
		err := vu.UploadObject(vu.Paths[x], uploadClient, ctx)

		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[x])
			log.Printf("error during the upload: %v. Error: %v", vu.Paths[x], err)
			returnChan <- err.Error()
		}

		returnChan <- ""
	}

	returnChan <- "upload completed"
}

func getClientUpload() (*storage.Client, context.Context, error)  {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil{
		return nil, nil, err
	}

	return client, ctx, nil
}