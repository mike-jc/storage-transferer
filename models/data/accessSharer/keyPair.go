package modelsDataAccessSharer

type KeyPair struct {
	PrivateKeyContainer PrivateKeyContainer
	PublicKeyContainer  PublicKeyContainer
}

type PrivateKeyContainer struct {
	Version    string `json:"version"`
	PrivateKey string `json:"privateKey"`
}

type PublicKeyContainer struct {
	Version   string `json:"version"`
	PublicKey string `json:"publicKey"`
}
