package main

// 加密
import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// 公钥和私钥可以从文件中读取
var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAy+906EAEirQ5RUxWNDvQe33P5TsDg5mZAPQvQuBBBkH+yqsH
SleG/zFoqFlgK7dy7NfIsrlx+9ASISSWxiQLIizk3KDRxoSr5FV2DYAvso/mobyd
XLQAhm75S1o9VQ+J3bOC27JXDS5QL7hmjEZCyipcsffgYrcF7/XO6/R3UI3Li9sJ
/oPdheqkTcP6Pb8vocwFrVb2NHBjPkAdHPltumH3tWI6RiNh5E5ZZPOUBZ0aIUrM
oV3vPxmssxBs4CRZE4Z1tilBvVc0OkltpDrLyTuBs9wlX6DiviW6pKYy2aGT2kmE
wuDz3e4c7pthXBl/OOl3lofa1hi764cviTCYHwIDAQABAoIBAFf1yVvfONZGk6kj
Gs9euTZ6dm/tuz9IwaiaqcPTi9hSIL5zdCqJhA2P1w89tXBFqMkk7UjBGbu97APl
jy6ZH0A3UuMibjiMwsMyZT+/eVMwJA7AlrMEZHGXbeklW+zTTeiU4600x71Eq4tZ
osmACJDAIskUG/EX9fSg9gXppIjWSqAGntZmPjhhXJOO2g2Tw2Fizc/xxP3/dv69
+YaCfHftuF+66qkx4GkAvT5M9paE+0BjxFgodYI4v0Yn7X9ltyAGiqZgHGUpGJpS
b1Lc3Mqq5f5jnxfaZQCpTTWHLRPwY8yPb4/AXrnAlqCOZQ9U8IvBHf/OvDOjLpaK
N5KPaxkCgYEA7egJ2CAQQcklYFiO1TWk/IO8MUizAr0XnO8GpDxUy/hjqhfab4Rf
u+aktZ51/N681pqF8EAqgMRmgXxu4AUECoLK5fGkp+IZBJoV2JsmUAqwBYvLNvor
BJSI4IIKMCB+yLFt2eJD7pWM+PEW7RulrcRFh4k70d9f8n2hwW3URgUCgYEA23ID
ZW1qrp1SNk0mv5Oux+xioaSOUM330NvBItYejPXXEVTLikYMMcwPEf4YecSph0QD
bPCvTN+fh6LF2g9FhjztRpCSc1GNc11WR+O4ZBfO+iE5ELUEb77Hjf1dmMnN//7j
XUnnof7IlJj0MAXI0IQ6y0Pftc26iF6zgCiS+tMCgYAKEeZIxaKqhi8U0urIz4p6
PcE7fM5G8WYMeHmZfgxAzfS6AGR4j+vVcj/KiDiKSYtIsiW1M6IY7TdBh9jRlqTD
JSIddYr4qDNS5IrELl0CylEFCxPA8fncKcVZa2eu/dEgAZKaxF8HvEDJULsdsivj
HQmsYPytN31CMFsmatWvWQKBgDdDJwfL3inLBIEYPMHR9xnxtYTvY8eFlvrJ3IFh
WqA06Mw8hmVz7m477S+ixZckp2yg/BvbIMpDJnGJ1DltzxxXC4nRro/L4ctDng7M
kgri1AS5iR1j+JILgUWIoKFxcKcfETLVAbgR6YFCY3wUeNXJ9uRpW1T1Uhw1fQ6x
KRJxAoGBANDpn4APUNTjgtI81PDNgleef6aCK0YrMWVHLSyZOmtg9jgU7NCVwi0/
ykDfVU9e8rFRQFo+EW4vFDNInKv8xuBBkj6ETZEPdQsjlS4dX9PrS1RruRE0SWg+
nMYg+OMyLNV7r9rb+Mrfwk4bXaZBst6DYfCwmF0cILy3ZrR/h8NX
-----END RSA PRIVATE KEY-----
	`)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAy+906EAEirQ5RUxWNDvQ
e33P5TsDg5mZAPQvQuBBBkH+yqsHSleG/zFoqFlgK7dy7NfIsrlx+9ASISSWxiQL
Iizk3KDRxoSr5FV2DYAvso/mobydXLQAhm75S1o9VQ+J3bOC27JXDS5QL7hmjEZC
yipcsffgYrcF7/XO6/R3UI3Li9sJ/oPdheqkTcP6Pb8vocwFrVb2NHBjPkAdHPlt
umH3tWI6RiNh5E5ZZPOUBZ0aIUrMoV3vPxmssxBs4CRZE4Z1tilBvVc0OkltpDrL
yTuBs9wlX6DiviW6pKYy2aGT2kmEwuDz3e4c7pthXBl/OOl3lofa1hi764cviTCY
HwIDAQAB
-----END PUBLIC KEY-----`)

func RsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func main() {
	data, err := RsaEncrypt([]byte("polaris@studygolang.com"))
	if err != nil {
		panic(err)
	}
	fmt.Println("data after encrypt: ", b64.StdEncoding.EncodeToString([]byte(data)))
	origData, err := RsaDecrypt(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(origData))
}
