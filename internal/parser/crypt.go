package parser

import (
	"bytes"
	"crypto/aes"
)

const blockSize = 16

func EncryptAESWithECB(data, key []byte) ([]byte, error) {
	// Make a copy of the input data to avoid modifying the original slice
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	// Calculate the required padding length
	paddingLen := blockSize - (len(dataCopy) % blockSize)

	// Append padding bytes
	padding := bytes.Repeat([]byte{byte(paddingLen)}, paddingLen)
	dataCopy = append(dataCopy, padding...)

	// Create a new AES cipher with the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Encrypt the padded data using ECB mode
	encrypted := make([]byte, len(dataCopy))
	for i := 0; i < len(dataCopy); i += blockSize {
		block.Encrypt(encrypted[i:i+blockSize], dataCopy[i:i+blockSize])
	}

	return encrypted, nil
}

func DecryptAESWithECB(encrypted, key []byte) ([]byte, error) {
	// Create a new AES cipher with the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Check if the encrypted data length is a multiple of the block size
	if len(encrypted)%blockSize != 0 {
		// Calculate the padding length needed to align the last block
		paddingLen := blockSize - (len(encrypted) % blockSize)

		// Append padding to the encrypted data to align the last block
		padded := append(encrypted, bytes.Repeat([]byte{byte(0)}, paddingLen)...)

		// Perform decryption of the padded data
		decrypted := make([]byte, len(padded))
		for i := 0; i < len(padded); i += blockSize {
			block.Decrypt(decrypted[i:i+blockSize], padded[i:i+blockSize])
		}

		return decrypted, nil
	}

	// Perform decryption of the encrypted data directly
	decrypted := make([]byte, len(encrypted))
	for i := 0; i < len(encrypted); i += blockSize {
		block.Decrypt(decrypted[i:i+blockSize], encrypted[i:i+blockSize])
	}

	return decrypted, nil
}
