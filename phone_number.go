package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golangpub/log"
	"github.com/nyaruka/phonenumbers"
)

type PhoneNumber struct {
	Code      int    `json:"code"`
	Number    int64  `json:"number"`
	Extension string `json:"extension,omitempty"`
}

func (n *PhoneNumber) String() string {
	if len(n.Extension) == 0 {
		return fmt.Sprintf("+%d%d", n.Code, n.Number)
	}

	return fmt.Sprintf("+%d%d-%s", n.Code, n.Number, n.Extension)
}

func (n *PhoneNumber) InternationalFormat() string {
	pn, err := phonenumbers.Parse(n.String(), "")
	if err != nil {
		return ""
	}
	return phonenumbers.Format(pn, phonenumbers.INTERNATIONAL)
}

func (n *PhoneNumber) MaskString() string {
	nnBytes := []byte(fmt.Sprint(n.Number))
	maskLen := (len(nnBytes) + 2) / 3
	start := len(nnBytes) - 2*maskLen
	for i := 0; i < maskLen; i++ {
		nnBytes[start+i] = '*'
	}

	nn := string(nnBytes)

	if len(n.Extension) == 0 {
		return fmt.Sprintf("+%d%s", n.Code, nn)
	}

	return fmt.Sprintf("+%d%s-%s", n.Code, nn, n.Extension)
}

func NewPhoneNumber(s string) (*PhoneNumber, error) {
	s = trimPhoneNumberString(s)
	pn, err := phonenumbers.Parse(s, "")
	if err != nil {
		return nil, err
	}

	if !phonenumbers.IsValidNumber(pn) {
		return nil, errors.New("invalid phone number")
	}

	return &PhoneNumber{
		Code:      int(pn.GetCountryCode()),
		Number:    int64(pn.GetNationalNumber()),
		Extension: pn.GetExtension(),
	}, nil
}

func NewPhoneNumberV2(s string, code int) (*PhoneNumber, error) {
	s = trimPhoneNumberString(s)
	pn, err := NewPhoneNumber(s)
	if err == nil {
		return pn, nil
	}
	if s[0] != '+' && code != 0 {
		s = fmt.Sprintf("+%d%s", code, s)
		return NewPhoneNumber(s)

	}
	return nil, err
}

func trimPhoneNumberString(s string) string {
	s = strings.Replace(s, "-", "", -1)
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	return s
}

func IsPhoneNumber(s string) bool {
	if s == "" {
		return false
	}
	n, err := phonenumbers.Parse(s, "")
	if err != nil {
		log.Errorf("Parse %s: %v", s, err)
		return false
	}
	return phonenumbers.IsValidNumber(n)
}
