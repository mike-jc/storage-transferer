package storages_test

import (
	"bytes"
	"github.com/astaxie/beego"
	"service-recordingStorage/services/storages"
	"service-recordingStorage/tests"
	"testing"
)

const BucketName = "24s-public"
const TestFileSize = 646908
const ChunkSize = 1024 * 1024

var s3 *storages.S3

func init() {
	tests.Init()

	s3 = new(storages.S3)
	s3.SetRemoveAfterTransferring(beego.AppConfig.DefaultBool("s3.removeAfterTransferring", false))
	s3.SetChunkSize(beego.AppConfig.DefaultInt64("s3.chunkSize", ChunkSize))
}

func TestCheckSourceFileOK(t *testing.T) {
	s3.SetBucketName(BucketName)
	_ = s3.SetSourceFile("sandbox/test-video.mp4")

	if err := s3.CheckSourceFile(); err != nil {
		t.Fatalf("Source file checking failed, should be successful: %s", err.Error())
	}
}

func TestCheckSourceFileFailed(t *testing.T) {
	s3.SetBucketName(BucketName)
	_ = s3.SetSourceFile("sandbox/non-existing.mp4")

	if err := s3.CheckSourceFile(); err == nil {
		t.Fatalf("Source file checking was successful, should failed")
	}
}

func TestDownloadOK(t *testing.T) {
	s3.SetBucketName(BucketName)
	_ = s3.SetSourceFile("sandbox/test-video.mp4")

	var buf bytes.Buffer
	for {
		if data, err := s3.Download(); err != nil {
			t.Fatalf("File downloading failed, should be successful")
		} else {
			if data == nil {
				break
			} else {
				buf.Write(data)
			}
		}
	}
	if buf.Len() != TestFileSize {
		t.Fatalf("File size is %d, should be %d", buf.Len(), TestFileSize)
	}
}

func TestCloseS3OK(t *testing.T) {
	_ = s3.SetSourceFile("sandbox/test-deletion.mp4")
	if err := s3.Close(); err != nil {
		t.Fatalf("Closing of storage failed, should be successful: %s", err.Error())
	}
}

func TestS3HealthCheckOK(t *testing.T) {
	if err := s3.HealthCheck(); err != nil {
		t.Fatalf("Health check failed, should be successful: %s", err.Error())
	}
}
