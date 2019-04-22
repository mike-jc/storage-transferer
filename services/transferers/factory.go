package transferers

import (
	"errors"
)

type FileTransfererFactory struct {
}

func (f *FileTransfererFactory) FileTransferByType(transfererType string) (transferer AbstractTransfererContract, err error) {

	switch transfererType {
	case "dracoon":
		transferer = new(AmazonToDracoon)
	case "fake":
		transferer = new(Fake)
	default:
		err = errors.New("wrong type of File Transferer in configuration")
	}
	return
}
