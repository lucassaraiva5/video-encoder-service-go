package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func init()  {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "convite.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}

	return video, repo
}

func TestVideoService_Download(t *testing.T) {

	video, repo := prepare()

	videoService := NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("encode-go-bucket")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)
	
}
