package utils

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

const letterBytes = "abcdef0123456789"
const (
	letterIdxBits = 4
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func GenerateID(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func MapKey[T any](m map[string]any, key string) (T, error) {
	val, ok := m[key]
	if !ok {
		return *new(T), fmt.Errorf("key %q not found", key)
	}

	typedVal, ok := val.(T)
	if !ok {
		return *new(T), fmt.Errorf("key %q is of type %T, expected type is %s", key, val, reflect.TypeOf((*T)(nil)).Elem())
	}

	return typedVal, nil
}
