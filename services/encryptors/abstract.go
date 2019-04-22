package encryptors

import "service-recordingStorage/models/data/encryptor"

type AbstractEncryptor struct {
	encryptionInfo modelsDataEncryptor.EncryptionInfo
	self           AbstractEncryptorContract

	AbstractEncryptorContract
}

type AbstractEncryptorContract interface {
	EncryptorName() string
	Status() string
	Init(options interface{}) error
	Reset()
	Encrypt(data []byte) (encryptedData []byte, err error)
	EncryptFinal() error
	Decrypt(data []byte, info modelsDataEncryptor.EncryptionInfo) (decryptedData []byte, err error)
	EncryptWithKey(data []byte, pemKey []byte) (encryptedData []byte, err error)
	DecryptWithKey(data []byte, pemKey []byte, password []byte) (decryptedData []byte, err error)
	EncryptionInfo() modelsDataEncryptor.EncryptionInfo
}

func (e *AbstractEncryptor) EncryptorName() string {
	return "Abstract encryptor"
}

func (e *AbstractEncryptor) Status() string {
	return "data encrypting by " + e.self.EncryptorName()
}

func (e *AbstractEncryptor) Init(options interface{}) error {
	e.self = e
	return nil
}

func (e *AbstractEncryptor) Reset() {
	e.encryptionInfo = modelsDataEncryptor.EncryptionInfo{}
}

// @Description encrypt data with symmetrical method
func (e *AbstractEncryptor) Encrypt(data []byte) (encryptedData []byte, err error) {
	return
}

// @Description finalize encrypting with symmetrical method
func (e *AbstractEncryptor) EncryptFinal() error {
	return nil
}

// @Description decrypt data with symmetrical method
func (e *AbstractEncryptor) Decrypt(data []byte, info modelsDataEncryptor.EncryptionInfo) (decryptedData []byte, err error) {
	return
}

// @Description encrypt data with public key
func (e *AbstractEncryptor) EncryptWithKey(data []byte, pemKey []byte) (encryptedData []byte, err error) {
	return
}

// @Description decrypt data with private key
func (e *AbstractEncryptor) DecryptWithKey(data []byte, pemKey []byte, password []byte) (decryptedData []byte, err error) {
	return
}

func (e *AbstractEncryptor) EncryptionInfo() modelsDataEncryptor.EncryptionInfo {
	return e.encryptionInfo
}
