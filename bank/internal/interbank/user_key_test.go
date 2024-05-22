package interbank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserKey(t *testing.T) {
	uk := NewUserKey(100, 3123123)
	assert.Equal(t, "100-3123123", uk.String())
}

func TestCreateUserKeyFromStr(t *testing.T) {
	uk, err := NewUserKeyFromStr("100-3123123")
	assert.Nil(t, err)
	assert.Equal(t, NewBankId(100), uk.BankId)
	assert.Equal(t, NewUserId(3123123), uk.UserId)
}

func TestCreateInvalidUserKeyFromStr(t *testing.T) {
	uk, err := NewUserKeyFromStr("100-312312312312323")
	assert.Nil(t, uk)
	assert.Equal(t, err.Error(), "invalid user id")
}

func TestCompareUserKey(t *testing.T) {
	expected := NewUserKey(100, 3123123)
	actual, err := NewUserKeyFromStr("100-3123123")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestIdsToStr(t *testing.T) {
	userId := NewUserId(100)
	bankId := NewBankId(3123)

	assert.Equal(t, "100", userId.String())
	assert.Equal(t, "3123", bankId.String())
}
