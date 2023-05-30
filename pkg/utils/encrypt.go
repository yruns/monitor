package utils

var Encrypt *Encryption

// AES对称加密
type Encryption struct {
	Key string
}

func init() {
	Encrypt = &Encryption{}
}

func (e *Encryption) SetKey(key string) {
	e.Key = key
}

func (e *Encryption) AesEncoding(s string) string {
	return s
}
