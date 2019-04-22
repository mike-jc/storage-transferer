package encryptors

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/spacemonkeygo/openssl"
	"github.com/youmark/pkcs8"
	"io"
	"service-recordingStorage/models/data/encryptor"
	"service-recordingStorage/system/crypto"
)

type HybridCrypto struct {
	options modelsDataEncryptor.HybridCryptoOptions

	encryptor *openssl.AuthenticatedEncryptionCipherCtx
	decryptor *openssl.AuthenticatedDecryptionCipherCtx

	AbstractEncryptor
}

func (e *HybridCrypto) EncryptorName() string {
	return "HybridCrypto encryptor"
}

func (e *HybridCrypto) Init(options interface{}) error {
	err := e.AbstractEncryptor.Init(options)
	if err != nil {
		return err
	}

	e.self = e

	if options == nil {
		e.options = modelsDataEncryptor.NewHybridCryptoOptionsDefault()
	} else if hcOptions, ok := options.(modelsDataEncryptor.HybridCryptoOptions); ok {
		if hcOptions.Hash == nil || hcOptions.Mgf1Hash == nil || hcOptions.SymmetricalMethod == 0 {
			return errors.New("Wrong values in HybridCrypto options")
		}
		e.options = hcOptions
	} else {
		return errors.New("Options can be of type HybridCryptoOptions or nil")
	}

	return err
}

func (e *HybridCrypto) Reset() {
	e.AbstractEncryptor.Reset()

	e.encryptor = nil
	e.decryptor = nil
}

func (e *HybridCrypto) Encryptor() (engine *openssl.AuthenticatedEncryptionCipherCtx, err error) {
	if e.encryptor == nil {
		params := modelsDataEncryptor.Methods[e.options.SymmetricalMethod]

		// random key
		key := make([]byte, params.KeyLength)
		if _, err = io.ReadFull(rand.Reader, key); err != nil {
			return
		}
		e.encryptionInfo.RandomKey = key

		// initial vector
		vector := make([]byte, params.IvSize)
		if _, err = io.ReadFull(rand.Reader, vector); err != nil {
			return nil, err
		}
		e.encryptionInfo.InitialVector = vector

		// engine
		var encryptor openssl.AuthenticatedEncryptionCipherCtx
		if encryptor, err = openssl.NewGCMEncryptionCipherCtx(params.BlockSize, nil, key, vector); err != nil {
			return nil, err
		}
		e.encryptor = &encryptor
	}
	return e.encryptor, nil
}

func (e *HybridCrypto) Decryptor() (engine *openssl.AuthenticatedDecryptionCipherCtx, err error) {
	if e.decryptor == nil {
		params := modelsDataEncryptor.Methods[e.options.SymmetricalMethod]

		// engine
		var decryptor openssl.AuthenticatedDecryptionCipherCtx
		if decryptor, err = openssl.NewGCMDecryptionCipherCtx(params.BlockSize, nil, e.encryptionInfo.RandomKey, e.encryptionInfo.InitialVector); err != nil {
			return nil, err
		}
		e.decryptor = &decryptor
	}
	return e.decryptor, nil
}

// @Description encrypt data with symmetrical method
func (e *HybridCrypto) Encrypt(data []byte) (encryptedData []byte, err error) {
	var encryptor *openssl.AuthenticatedEncryptionCipherCtx
	if encryptor, err = e.Encryptor(); err != nil {
		return
	}

	return (*encryptor).EncryptUpdate(data)
}

// @Description finalize encrypting data with symmetrical method
func (e *HybridCrypto) EncryptFinal() error {
	if encryptor, err := e.Encryptor(); err != nil {
		return err
	} else {
		if _, err := (*encryptor).EncryptFinal(); err != nil {
			return err
		}
		if tag, err := (*encryptor).GetTag(); err != nil {
			return err
		} else {
			e.encryptionInfo.AuthenticationTag = tag
		}
	}
	return nil
}

// @Description encrypt data with symmetrical method
func (e *HybridCrypto) Decrypt(data []byte, info modelsDataEncryptor.EncryptionInfo) (decryptedData []byte, err error) {
	e.encryptionInfo = info

	// try to decode data, if failed then use original data
	var decodedData []byte
	if decodedData, err = base64.StdEncoding.DecodeString(string(data)); err == nil {
		data = decodedData
	}

	// get decryptor
	var decryptor *openssl.AuthenticatedDecryptionCipherCtx
	if decryptor, err = e.Decryptor(); err != nil {
		return
	}

	// decrypt
	if err = (*decryptor).SetTag(info.AuthenticationTag); err != nil {
		return
	}
	if decryptedData, err = (*decryptor).DecryptUpdate(data); err != nil {
		return nil, err
	}
	if _, err = (*decryptor).DecryptFinal(); err != nil {
		return nil, err
	}

	return
}

// @Description encrypt data
func (e *HybridCrypto) EncryptWithKey(data []byte, pemKey []byte) (encryptedData []byte, err error) {
	block, _ := pem.Decode(pemKey)

	var result interface{}
	if result, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		err = errors.New("Can not parse public key: " + err.Error())
		return
	}
	publicKey := result.(*rsa.PublicKey)

	return systemCryptoRSA.EncryptOAEP(e.options.Hash, e.options.Mgf1Hash, rand.Reader, publicKey, data, []byte{})
}

// @Description encrypt data
func (e *HybridCrypto) DecryptWithKey(data []byte, pemKey []byte, password []byte) (decryptedData []byte, err error) {

	// try to decode data, if failed then use original data
	var decodedData []byte
	if decodedData, err = base64.StdEncoding.DecodeString(string(data)); err == nil {
		data = decodedData
	}

	// decrypt key
	block, _ := pem.Decode(pemKey)

	var result interface{}
	if result, err = pkcs8.ParsePKCS8PrivateKey(block.Bytes, password); err != nil {
		err = errors.New("Can not parse private key: " + err.Error())
		return
	}
	privateKey := result.(*rsa.PrivateKey)

	// decrypt data
	return systemCryptoRSA.DecryptOAEP(e.options.Hash, e.options.Mgf1Hash, rand.Reader, privateKey, data, []byte{})
}
