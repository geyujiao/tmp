package proxy

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

//Base64_encode(AES(Json_encode(消息体)))
//AES加密方式：CBC128位/PKCS7(JAVA PKCS5)/iv(1234567890123456)

const ivStr = "1234567890123456"

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, []byte(ivStr))
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
