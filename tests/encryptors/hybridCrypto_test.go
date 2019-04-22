package encryptors_test

import (
	"bytes"
	"service-recordingStorage/services"
	"service-recordingStorage/services/encryptors"
	"service-recordingStorage/tests"
	"testing"
)

var encryptor *encryptors.HybridCrypto

func init() {
	tests.Init()

	encryptor = new(encryptors.HybridCrypto)
}

func TestEncryptingOK(t *testing.T) {
	if err := encryptor.Init(nil); err != nil {
		t.Fatalf("Could not initialize encryptor: %s", err.Error())
	}

	data := []byte("some test data for encryption")
	if encryptedData, err := encryptor.Encrypt(data); err != nil {
		t.Fatalf("Encryption failed, should be successful: %s", err.Error())
	} else {
		if err := encryptor.EncryptFinal(); err != nil {
			t.Fatalf("Can not finalize encryption: %s", err.Error())
		}
		if info := encryptor.EncryptionInfo(); !info.ValidAesGcmInfo() {
			t.Fatalf("Encryption info is wrong: %+v", info)
		} else {
			if decryptedData, err := encryptor.Decrypt(encryptedData, info); err != nil {
				t.Fatalf("Decryption failed, should be successful: %s", err.Error())
			} else if !bytes.Equal(data, decryptedData) {
				t.Fatalf("Decryption failed, decrypted data does not equal to original data. Original: %s, decrypted: %s ", data, decryptedData)
			}
		}
	}
}

func TestEncryptingWithKeyOK(t *testing.T) {
	pemPrivateKey := []byte("-----BEGIN ENCRYPTED PRIVATE KEY-----\r\n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\r\n-----END ENCRYPTED PRIVATE KEY-----\r\n")
	pemPublicKey := []byte("-----BEGIN PUBLIC KEY-----\r\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjdLi8iXMEi0z1G25iz+G\r\nUfakrS69dJNRdt/YvLI64XOEfdljx9qsoGUJvqSN5SYYQrnMYrKfUL12ZeNnM1kw\r\njHwjHinll9dF6uwL+ITMzXbmiSXQ0hV1CZOcjZRkMXejDtmShtrMpvd3EPCBRSwt\r\np/j8W1CMHTsSLhgRrfb/k858gf/MjEOQY2QnQrXuGgd3jquJ3XyX+wxNAvrzgWif\r\nK9al9+S0c4gxrLy15quFzW8CP7G4ajyEz1J8/oHJu3SJRV2s1iA4ohdEA9cQxcXv\r\nflaDp16BUylHIrlIONfBVv7qD0BwhDi7yxq1cYN0VePz4qbUt3/d2OHAetK/eC8c\r\nPwIDAQAB\r\n-----END PUBLIC KEY-----\r\n")
	data := []byte("some test data for encryption")

	if password, err := services.DracoonEncryptionPassword(services.MessageDomain(nil)); err != nil {
		t.Fatalf("Can not get Dracoon encryption password %s", err.Error())
	} else {
		if encryptedData, err := encryptor.EncryptWithKey(data, pemPublicKey); err != nil {
			t.Fatalf("Encryption with public key failed, should be successful: %s", err.Error())
		} else {
			if decryptedData, err := encryptor.DecryptWithKey(encryptedData, pemPrivateKey, []byte(password)); err != nil {
				t.Fatalf("Decryption with private key failed, should be successful: %s", err.Error())
			} else if !bytes.Equal(data, decryptedData) {
				t.Fatalf("Decryption with private key failed, decrypted data does not equal to original data. Original: %s, decrypted: %s ", data, decryptedData)
			}
		}
	}
}
