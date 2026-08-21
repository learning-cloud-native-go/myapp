package model

import (
	"errors"
	"fmt"
)

type BookStatus uint8

var ErrBookStatusInvalid = errors.New("invalid book status")

const (
	BookStatusPending  BookStatus = iota // 0
	BookStatusVerified                   // 1
)

const (
	BookStatusPendingStr  = "pending"
	BookStatusVerifiedStr = "verified"
)

var mapBookStatusStrToInt = map[string]BookStatus{
	BookStatusPendingStr:  BookStatusPending,
	BookStatusVerifiedStr: BookStatusVerified,
}

var mapBookStatusIntToStr = map[BookStatus]string{
	BookStatusPending:  BookStatusPendingStr,
	BookStatusVerified: BookStatusVerifiedStr,
}

func ParseBookStatus(str string) (*BookStatus, error) {
	v, ok := mapBookStatusStrToInt[str]
	if !ok {
		return nil, ErrBookStatusInvalid
	}

	return &v, nil
}

func (s BookStatus) String() string {
	return mapBookStatusIntToStr[s]
}

func (s BookStatus) MarshalJSON() ([]byte, error) {
	v, ok := mapBookStatusIntToStr[s]
	if !ok {
		return nil, ErrBookStatusInvalid
	}

	return []byte(fmt.Sprintf("%q", v)), nil
}
