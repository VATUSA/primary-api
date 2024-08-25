package types

import (
	"database/sql/driver"
	"fmt"
)

type FeedbackRating string

const (
	Unsatisfactory FeedbackRating = "unsatisfactory"
	Poor           FeedbackRating = "poor"
	Fair           FeedbackRating = "fair"
	Good           FeedbackRating = "good"
	Excellent      FeedbackRating = "excellent"
)

func (s *FeedbackRating) Scan(value interface{}) error {
	bytesValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan FeedbackRating: expected []byte, got %T", value)
	}

	strValue := string(bytesValue)
	switch FeedbackRating(strValue) {
	case Unsatisfactory, Poor, Fair, Good, Excellent:
		*s = FeedbackRating(strValue)
	default:
		return fmt.Errorf("invalid FeedbackRating value: %s", strValue)
	}
	return nil
}

func (s *FeedbackRating) Value() (driver.Value, error) {
	return string(*s), nil
}
