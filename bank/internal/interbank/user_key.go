package interbank

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
)

type BankId [2]byte
type UserId [4]byte

type UserKey struct {
	BankId BankId
	UserId UserId
}

func NewUserKeyFromStr(value string) (*UserKey, error) {
	if !bytes.Contains([]byte(value), []byte("-")) {
		return nil, errors.New("invalid user key format")
	}
	bankIdStr := bytes.Split([]byte(value), []byte("-"))[0]
	userIdStr := bytes.Split([]byte(value), []byte("-"))[1]

	bankId, err := strconv.Atoi(string(bankIdStr))
	if err != nil {
		return nil, err
	}
	if bankId > math.MaxUint16 {
		return nil, errors.New("invalid bank id")
	}

	userId, err := strconv.Atoi(string(userIdStr))
	if err != nil {
		return nil, err
	}
	if userId > math.MaxUint32 {
		return nil, errors.New("invalid user id")
	}

	return NewUserKey(uint16(bankId), uint32(userId)), nil
}

func NewUserKey(bankId uint16, userId uint32) *UserKey {
	bank := BankId{}
	user := UserId{}

	binary.BigEndian.PutUint16(bank[:], bankId)
	binary.BigEndian.PutUint32(user[:], userId)

	return &UserKey{
		BankId: bank,
		UserId: user,
	}
}

func (uk *UserKey) Bytes() []byte {
	b := make([]byte, 6)
	copy(b[:2], uk.BankId[:])
	copy(b[2:], uk.UserId[:])
	return b
}

func (uk *UserKey) String() string {
	bankId := binary.BigEndian.Uint16(uk.BankId[:])
	userId := binary.BigEndian.Uint32(uk.UserId[:])
	return fmt.Sprintf("%d-%d", bankId, userId)
}
