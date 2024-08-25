package types

import (
	"database/sql/driver"
	"fmt"
)

type StatusType string

const (
	All      StatusType = ""
	Pending  StatusType = "pending"
	Accepted StatusType = "accepted"
	Rejected StatusType = "rejected"
)

func (s *StatusType) Scan(value interface{}) error {
	bytesValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StatusType: expected []byte, got %T", value)
	}

	strValue := string(bytesValue)
	switch StatusType(strValue) {
	case All, Pending, Accepted, Rejected:
		*s = StatusType(strValue)
	default:
		return fmt.Errorf("invalid StatusType value: %s", strValue)
	}
	return nil
}

func (s *StatusType) Value() (driver.Value, error) {
	return string(*s), nil
}
