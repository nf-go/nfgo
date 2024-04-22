// Copyright 2022 The nfgo Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ncrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/nf-go/nfgo/nerrors"
)

const (
	aesGCMDefaultNonceSizeBytes = 12
)

func NewAESKey(keySizeBits int) ([]byte, error) {
	key, err := newRandomBytes(keySizeBits / 8)
	if err != nil {
		return nil, nerrors.Errorf("ncrypto failed to generate aes key: %v", err)
	}
	return key, nil
}

func NewAESKeyString(keySizeBits int) (string, error) {
	key, err := NewAESKey(keySizeBits)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func AESEncrypt(planData []byte, key []byte) ([]byte, error) {
	gcm, err := newBlockCipher(key, aesGCMDefaultNonceSizeBytes)
	if err != nil {
		return nil, err
	}

	nonce, err := newRandomBytes(gcm.NonceSize())
	if err != nil {
		return nil, nerrors.Errorf("ncrypt failed to generate nonce: %v", err)
	}

	encrypted := gcm.Seal(nil, nonce, planData, nil)
	return append(nonce, encrypted...), nil
}

func AESDecrypt(encryptedData []byte, key []byte) ([]byte, error) {
	gcm, err := newBlockCipher(key, aesGCMDefaultNonceSizeBytes)
	if err != nil {
		return nil, err
	}
	plainData, err := gcm.Open(nil, encryptedData[:aesGCMDefaultNonceSizeBytes], encryptedData[aesGCMDefaultNonceSizeBytes:], nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt value: %v", err)
	}
	return plainData, nil
}

func AESEncryptString(plainText string, keyText string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(keyText)
	if err != nil {
		return "", nerrors.WithStack(err)
	}
	plainData := []byte(plainText)
	encryptedData, err := AESEncrypt(plainData, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func AESDecryptString(encryptedText string, keyText string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(keyText)
	if err != nil {
		return "", nerrors.WithStack(err)
	}
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", nerrors.WithStack(err)
	}
	decryptData, err := AESDecrypt(encryptedData, key)
	if err != nil {
		return "", err
	}
	return string(decryptData), nil
}

func newRandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return nil, nerrors.Wrap(err, "ncrypto failed to generate random bytes")
	}
	return bytes, nil
}

func newBlockCipher(key []byte, nonceSizeBytes int) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nerrors.Wrap(err, "ncrypto failed to construct AES cipher")
	}
	gcm, err := cipher.NewGCMWithNonceSize(block, nonceSizeBytes)
	if err != nil {
		return nil, nerrors.Wrap(err, "ncrypto failed to construct block cipher: %v")
	}
	return gcm, nil
}
