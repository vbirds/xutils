package xutils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"os"
	"time"
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

const (
	aeskey = "nlasfl2wrwfsnsfs131#$%fs"
)

type Licences struct {
	Address       string    // 发布证书服务地址
	ServeGuid     string    // 服务唯一Id
	ServeName     string    // 服务名称
	MaxNumber     int       // 最大接入设备数目
	CreatedAt     time.Time // 创建时间
	EffectiveTime int       // 有效时长
}

// WriteLicences 写入文件
func CreateLicences(fileName string, l *Licences) error {
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}
	xpass, err := AesEncrypt(data, []byte(aeskey))
	if err != nil {
		return err
	}
	os.Remove(fileName)
	fp, err := os.Create(fileName)
	if err != nil {
		return err
	}
	fp.Write(xpass)
	defer fp.Close()
	return nil
}

func ReadLicences(filename string) (*Licences, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	//获取文件内容
	info, _ := fp.Stat()
	buf := make([]byte, info.Size())
	fp.Read(buf)
	tpass, err := AesDecrypt(buf, []byte(aeskey))
	if err != nil {
		return nil, err
	}
	var lice Licences
	json.Unmarshal(tpass, &lice)
	return &lice, nil
}
