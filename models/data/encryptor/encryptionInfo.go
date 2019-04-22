package modelsDataEncryptor

type EncryptionInfo struct {
	RandomKey         []byte
	InitialVector     []byte // a.k.a. Nonce
	AuthenticationTag []byte // a.k.a. MAC
}

func (i *EncryptionInfo) ValidAesGcmInfo() bool {
	return i.RandomKey != nil && len(i.RandomKey) > 0 &&
		i.InitialVector != nil && len(i.InitialVector) > 0 &&
		i.AuthenticationTag != nil && len(i.AuthenticationTag) > 0
}
