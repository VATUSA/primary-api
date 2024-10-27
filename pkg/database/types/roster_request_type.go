package types

import (
	"database/sql/driver"
	"errors"
)

type RequestType string

const (
	Visiting     RequestType = "visiting"
	Transferring RequestType = "transferring"
)

func (s *RequestType) Scan(value interface{}) error {
	if byteSlice, ok := value.([]byte); ok {
		*s = RequestType(byteSlice)
	} else {
		return errors.New("failed to scan request_type")
	}

	return nil
}

func (s *RequestType) Value() (driver.Value, error) {
	return string(*s), nil
}
