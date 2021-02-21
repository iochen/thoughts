package utils

import (
	"bufio"
	"crypto/rand"
	"strings"
)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStr(l uint) string {
	buf := bufio.NewReader(rand.Reader)
	count := uint(0)
	builder := strings.Builder{}
	for count < l {
		b, _ := buf.ReadByte()
		if isChar(b) {
			builder.WriteByte(b)
			count++
		}
	}
	return builder.String()
}

func isChar(b byte) bool {
	if (b >= 48 && b <= 57) ||
		(b >= 65 && b <= 90) ||
		(b >= 97 && b <= 122) {
		return true
	}
	return false
}
