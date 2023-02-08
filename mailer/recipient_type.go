package mailer

import (
	"errors"
	"strings"
)

var (
	ErrUnknownRecipientType = errors.New("unknown recipient type")
)

type RecipientType string

const (
	RecipientTypeDefault         RecipientType = "R" // 수신자
	RecipientTypeCarbonCopy      RecipientType = "C" // 참조자
	RecipientTypeBlindCarbonCopy RecipientType = "B" // 숨은 참조자
)

func (rt RecipientType) String() string {
	s := strings.ToUpper(string(rt))
	switch s {
	case "R", "C", "B":
		return s
	default:
		return "UNKNOWN"
	}
}

func ParseRecipientType(recipientTypeStr string) (RecipientType, error) {
	switch strings.ToUpper(recipientTypeStr) {
	case "R":
		return RecipientTypeDefault, nil
	case "C":
		return RecipientTypeCarbonCopy, nil
	case "B":
		return RecipientTypeBlindCarbonCopy, nil
	default:
		return RecipientType("UNKNOWN"), ErrUnknownRecipientType
	}
}
