package utils

import (
	"testing"
)

func TestGenerateRandomArray(t *testing.T) {
	l := GenerateRandomArray(1000)

	t.Log(l)
}
