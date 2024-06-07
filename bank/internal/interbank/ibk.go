package interbank

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
)

type BankId uint16
type UserId uint32

type IBK struct {
	BankId BankId
	UserId UserId
}

func NewIBKFromStr(value string) (*IBK, error) {
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

	return NewIBK(uint16(bankId), uint32(userId)), nil
}

func NewIBK(bankId uint16, userId uint32) *IBK {
	return &IBK{
		BankId: BankId(bankId),
		UserId: UserId(userId),
	}
}

func (uk IBK) String() string {
	return fmt.Sprintf("%d-%d", uk.BankId, uk.UserId)
}

func (uk *IBK) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	ptr, err := NewIBKFromStr(s)
	if err != nil {
		return err
	}

	*uk = *ptr

	return err
}

func (uk IBK) MarshalJSON() ([]byte, error) {
	return json.Marshal(uk.String())
}

func NewBankId(id uint16) BankId {
	return BankId(id)
}

func NewUserId(id uint32) UserId {
	return UserId(id)
}

func (ui UserId) String() string {
	return strconv.Itoa(int(ui))
}

func (bi BankId) String() string {
	return strconv.Itoa(int(bi))
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

	*bi = BankId(id)
	return nil
}
