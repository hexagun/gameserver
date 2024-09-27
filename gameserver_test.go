package gameserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	result := Add(2, 3)

	// Assert equality using Testify
	assert.Equal(t, 5, result, "they should be equal")
}
