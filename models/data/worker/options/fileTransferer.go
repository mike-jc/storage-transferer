package modelsDataWorkerOptions

type FileTransfererOptions struct {
	QueueReaderOptions QueueReaderOptions
}

func DefaultFileTransfererOptions() FileTransfererOptions {
	return FileTransfererOptions{
		QueueReaderOptions: DefaultQueueReaderOptions(),
	}
}
