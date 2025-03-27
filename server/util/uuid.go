package util

import (
	"fmt"

	nanoid "github.com/matoous/go-nanoid/v2"
)

// 生成uuid
func Uuid() (string, error) {
	id, err := nanoid.New()
	if err != nil {
		fmt.Println("Error generating NanoID:", err)
		return "", err
	}
	return id, nil
}
