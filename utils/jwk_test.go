package utils

import (
	"fmt"
	"testing"
)

func TestJwk(t *testing.T) {
	if err := SetJWKSFile("default", "hello", nil); err != nil {
		fmt.Println(err)
		return
	}
}
