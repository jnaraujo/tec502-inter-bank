package interbank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIBK(t *testing.T) {
	uk := NewIBK(100, 3123123)
	assert.Equal(t, "100-3123123", uk.String())
}

func TestCreateIBKFromStr(t *testing.T) {
	uk, err := NewIBKFromStr("100-3123123")
	assert.Nil(t, err)
	assert.Equal(t, NewBankId(100), uk.BankId)
	assert.Equal(t, NewUserId(3123123), uk.UserId)
}

func TestCreateInvalidIBKFromStr(t *testing.T) {
	uk, err := NewIBKFromStr("100-312312312312323")
	assert.Nil(t, uk)
	assert.Equal(t, err.Error(), "invalid user id")
}

func TestCompareIBK(t *testing.T) {
	expected := NewIBK(100, 3123123)
	actual, err := NewIBKFromStr("100-3123123")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestIdsToStr(t *testing.T) {
	userId := NewUserId(100)
	bankId := NewBankId(3123)

	assert.Equal(t, "100", userId.String())
	assert.Equal(t, "3123", bankId.String())
}
