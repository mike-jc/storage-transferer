package modelsDataWorkerOptions

type QueueReaderOptions struct {
	VisibilityTimeout   int64 // in seconds
	WaitTimeSeconds     int64 // in seconds
	MaxNumberOfMessages int64
	RequestDelay        int64 // in seconds
}

func DefaultQueueReaderOptions() QueueReaderOptions {
	return QueueReaderOptions{
		VisibilityTimeout:   5,
		WaitTimeSeconds:     3,
		MaxNumberOfMessages: 10,
		RequestDelay:        3,
	}
}
