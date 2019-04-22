package storages

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"io/ioutil"
	"regexp"
	"service-recordingStorage/errors/storage"
	"service-recordingStorage/models/data/sqs"
	"service-recordingStorage/models/data/storage"
	"service-recordingStorage/services/amazon"
)

type S3 struct {
	engine *s3.S3

	bucketName string
	key        string

	amazon.AbstractService
	AbstractStorage
}

func (s *S3) StorageName() string {
	return "S3 storage"
}

func (s *S3) Init(message modelsDataSqs.FileTransfererMessage, lg *logger.Logger) errorsStorage.ErrorContract {
	err := s.AbstractStorage.Init(message, lg)
	if err != nil {
		return err
	}

	s.self = s

	if bucketName, mErr := message.S3BucketName(); mErr == nil {
		s.SetBucketName(bucketName)
	} else {
		s.Logger().Log(logger.CreateError("Could not get S3 bucket name: " + mErr.Error()).SetErrorCode(15009))
		return errorsStorage.NewError(mErr.Error(), errorsStorage.Message, errorsStorage.FileError)
	}

	return err
}

func (s *S3) SetBucketName(name string) {
	s.bucketName = name
}

func (s *S3) Engine() *s3.S3 {
	if s.engine == nil {
		s.engine = s3.New(s.Session(), aws.NewConfig() /*.WithLogLevel(aws.LogDebugWithHTTPBody)*/)
	}
	return s.engine
}

func (s *S3) SourceFiles(message modelsDataSqs.FileTransfererMessage) (files []interface{}, err errorsStorage.ErrorContract) {
	if keys, mErr := message.S3BucketKeys(); mErr != nil {
		err = errorsStorage.NewError("Could not parse S3 bucket keys from queue message: "+mErr.Error(), errorsStorage.Message, errorsStorage.FileError)
	} else {
		files = make([]interface{}, len(keys))
		for i, key := range keys {
			files[i] = key
		}
	}
	return
}

func (s *S3) SetSourceFile(file interface{}) errorsStorage.ErrorContract {
	if err := s.AbstractStorage.SetSourceFile(file); err != nil {
		return err
	}

	s.key = file.(string)
	return nil
}

// @Description check if source file exists
// (if we're gonna download from the storage)
// and set file info
func (s *S3) CheckSourceFile() errorsStorage.ErrorContract {
	err := s.AbstractStorage.CheckSourceFile()
	if err != nil {
		return err
	}

	var output *s3.HeadObjectOutput
	var cErr error
	output, cErr = s.Engine().HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s.key),
	})
	if cErr == nil {
		s.fileInfo.Path = s.key
		s.fileInfo.Name = modelsDataStorage.NameFromPath(s.key)
		s.fileInfo.Extension = modelsDataStorage.ExtensionFromPath(s.key)

		if output.ContentLength != nil && *output.ContentLength > 0 {
			s.state.Range.Limit = *output.ContentLength - 1
			s.fileInfo.Size = int(*output.ContentLength)
		}
		if output.ContentType != nil && len(*output.ContentType) > 0 {
			ext := modelsDataStorage.ExtensionFromMime(*output.ContentType)
			if len(ext) > 0 {
				s.fileInfo.Extension = ext
			}
		}
		if len(s.fileInfo.Extension) == 0 {
			if matched, mErr := regexp.MatchString(`\.mp3$`, s.key); matched && mErr == nil {
				s.fileInfo.Extension = "mp3"
			} else {
				s.fileInfo.Extension = "mp4"
			}
		}
	} else {
		if rErr, ok := cErr.(s3.RequestFailure); ok && rErr.StatusCode() == 404 {
			return errorsStorage.NewError(cErr.Error(), errorsStorage.Source, errorsStorage.NotFoundError)
		}
		return errorsStorage.NewError(cErr.Error(), errorsStorage.Source, errorsStorage.GeneralError)
	}

	return nil
}

// @Description download data from the storage (probably chunk by chunk)
func (s *S3) Download() (data []byte, err errorsStorage.ErrorContract) {
	data, err = s.AbstractStorage.Download()
	if err != nil {
		return
	}

	// if the end exists, set start as the next byte after previous end
	if s.state.Range.End > 0 {
		s.state.Range.Start = s.state.Range.End + 1
	}

	// if the limit exists, check if we reach the limit (file size)
	if s.state.Range.Limit > 0 && s.state.Range.Start > s.state.Range.Limit {
		return nil, nil
	}

	// move range to the next chunk
	s.state.Range.End = s.state.Range.Start + s.ChunkSize()

	// download
	var output *s3.GetObjectOutput
	var rErr error
	output, rErr = s.Engine().GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s.key),
		Range:  aws.String(s.state.Range.RangeHeader()),
	})
	if output != nil && output.Body != nil {
		defer output.Body.Close()
	}

	if rErr == nil {
		if data, rErr = ioutil.ReadAll(output.Body); rErr == nil {
			s.state.Range.End = s.state.Range.Start + int64(len(data)) - 1 // adjust range
		}
	} else {
		if aErr, ok := err.(awserr.RequestFailure); ok && aErr.StatusCode() == 416 {
			rErr = nil // just return empty data to signal that downloading's completed
		}
	}
	if rErr != nil {
		err = errorsStorage.NewError(rErr.Error(), errorsStorage.Source, errorsStorage.DownloadError)
	}
	return
}

// @Description do some actions after completing transferring the file to all upload targets
// (e.g., remove source file, etc)
func (s *S3) Close() errorsStorage.ErrorContract {
	if err := s.AbstractStorage.Close(); err != nil {
		return err
	}

	if s.removeAfterTransferring {
		_, cErr := s.Engine().DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(s.key),
		})
		if cErr != nil {
			return errorsStorage.NewError(cErr.Error(), errorsStorage.Source, errorsStorage.DeletionError)
		}
	}

	return nil
}

// @Description check if service works fine
func (s *S3) HealthCheck() errorsStorage.ErrorContract {
	if err := s.AbstractStorage.HealthCheck(); err != nil {
		return err
	}

	// check if bucket exists and we have access to it
	_, cErr := s.Engine().HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(s.bucketName),
	})
	if cErr != nil {
		return errorsStorage.NewError(cErr.Error(), errorsStorage.Source, errorsStorage.ApiError)
	}
	return nil
}
