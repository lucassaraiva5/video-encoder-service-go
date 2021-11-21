package services_test

import (
	"encoder/application/services"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func init()  {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestVideoService_Upload(t *testing.T) {

	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("encode-go-bucket")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "encode-go-bucket"

	videoUpload.VideoPath = os.Getenv("LOCAL_STORAGE_PATH") + video.ID

	doneUpload := make(chan string)
	videoUpload.ProcessUpload(5, doneUpload)

	result := <-doneUpload
	require.Equal(t, result, "upload completed")
}