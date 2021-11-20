package domain

import (
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"time"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Job struct {
	ID string `valid:"uuid"`
	OutputBucketPath string `valid:"notnull"`
	Status string `valid:"notnull"`
	Video *Video `valid:"-"`
	VideoId string `valid:"-"`
	Error string `valid:"-"`
	CreatedAt time.Time `valid:"-"`
	UpdatedAt time.Time `valid:"-"`
}

func NewJob(output string, status string, video *Video) (*Job, error) {
	job := Job{
		OutputBucketPath: output,
		Status: status,
		Video: video,
	}

	job.prepare()

	err := job.Validate()

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (job *Job) prepare()  {
	job.ID = uuid.NewV4().String()
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
}

func (job *Job) Validate() error {
	_, err := govalidator.ValidateStruct(job)

	if err != nil {
		return err
	}

	return nil
}