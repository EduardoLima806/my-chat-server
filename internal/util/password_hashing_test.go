package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var passHasher PasswordHasher = &DefaultPasswordHasher{}

func Test_If_Password_Is_Hashed(t *testing.T) {
	hash, err := passHasher.HashPassword("P4$$w0rd")
	assert.NotEmpty(t, hash)
	assert.Nil(t, err)
}

func Test_If_Password_Matched_With_Hash(t *testing.T) {
	pass := "P4$$w0rd"
	hash, _ := passHasher.HashPassword(pass)
	assert.True(t, passHasher.VerifyPassword(pass, hash))
}
