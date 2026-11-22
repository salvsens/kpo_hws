package domain

import "errors"

// ErrInvalidWeight indicates that the weight is outside the permitted range.
var ErrInvalidWeight = errors.New("invalid weight")

// Weight is a whole percent from 1 to 100.
// The zero value is invalid.
type Weight struct {
	percent int
}

// NewWeight returns a weight with the specified percent.
// It returns ErrInvalidWeight if the percent is outside the range from 1 to 100.
func NewWeight(percent int) (Weight, error) {
	if percent <= 0 || percent > 100 {
		return Weight{}, ErrInvalidWeight
	}
	return Weight{percent: percent}, nil
}

// Percent returns the weight in percent.
func (w Weight) Percent() int {
	return w.percent
}
