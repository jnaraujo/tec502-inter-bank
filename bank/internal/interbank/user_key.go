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

	binary.LittleEndian.PutUint16(bank[:], bankId)
	binary.LittleEndian.PutUint32(user[:], userId)

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
	bankId := binary.LittleEndian.Uint16(uk.BankId[:])
	userId := binary.LittleEndian.Uint32(uk.UserId[:])
	return fmt.Sprintf("%d-%d", bankId, userId)
}

func NewBankId(id uint16) BankId {
	bi := BankId{}
	binary.LittleEndian.PutUint16(bi[:], id)
	return bi
}

func NewUserId(id uint32) UserId {
	ui := UserId{}
	binary.LittleEndian.PutUint32(ui[:], id)
	return ui
}

func (ui UserId) String() string {
	return strconv.Itoa(int(binary.LittleEndian.Uint16(ui[:])))
}

func (bi BankId) String() string {
	return strconv.Itoa(int(binary.LittleEndian.Uint16(bi[:])))
}

func (bi BankId) MarshalText() (text []byte, err error) {
	return []byte(bi.String()), nil
}

func (bi *BankId) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return errors.New("empty bank id")
	}

	id, err := strconv.Atoi(string(text))
	if err != nil {
		return err
	}

	binary.LittleEndian.PutUint16(bi[:], uint16(id))
	return nil
}
