package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"crypto/rand"
)

/**
 * 对一个字符串进行MD5加密,不可解密
 */
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s)) //使用zhifeiya名字做散列值，设定后不要变
	return hex.EncodeToString(h.Sum(nil))
}

/**
 * 获取一个Guid值
 */
func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}


//xxtea 加密
func XXEncrypt(str, key string) string {
	encrypt_data:= Encrypt([]byte(str), []byte(key))
	return base64.StdEncoding.EncodeToString(encrypt_data)
}

//xxtea 解密
func XXDecrypt(str, key string) (string,error) {
	encrypt_data,err:=base64.StdEncoding.DecodeString(str)
	if err!=nil{
		return "",err
	}
	return  string(Decrypt(encrypt_data, []byte(key))),nil
}

