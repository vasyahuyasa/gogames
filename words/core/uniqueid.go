package core

import (
	"math/rand"
)

const alphavite = "abcdefghijklmnopqrstuvwxyz123456789"
const alphaviteLen = len(alphavite)

// uniqueID возвращает уникальный идентификатор
func uniqueID(length int) string {
	id := make([]byte, length)
	for i := range id {
		id[i] = alphavite[rand.Int()%alphaviteLen]
	}
	return string(id)
}
