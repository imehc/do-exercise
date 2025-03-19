package util

import "golang.org/x/crypto/bcrypt"

type Hash struct {
	Value string
}

// Hash 加密
func (h *Hash) Hash() (string, error) {
	// cost 参数等同于 js 中的 saltRounds，这里设置为10
	bytes, err := bcrypt.GenerateFromPassword([]byte(h.Value), 10)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Compare 验证
func (h *Hash) Compare(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h.Value), []byte(hash))
	return err == nil
}
