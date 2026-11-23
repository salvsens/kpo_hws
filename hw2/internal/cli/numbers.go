package cli

import (
	"fmt"
	"strconv"
	"strings"
)

// scoreDecimalPlaces is the number of digits after the decimal point in a score.
const scoreDecimalPlaces = 2

// hundredthsPerPoint is the number of hundredths in one point.
const hundredthsPerPoint = 100

// parseHundredths converts a decimal score to hundredths of a point.
func parseHundredths(text string) (int64, error) {
	sign, digits := "", text
	if strings.HasPrefix(digits, "-") || strings.HasPrefix(digits, "+") {
		sign, digits = digits[:1], digits[1:]
	}
	whole, fraction, hasPoint := strings.Cut(digits, ".")
	if !validDecimal(whole, fraction, hasPoint) {
		return 0, fmt.Errorf("число %q задано неверно. Введите число с точкой и не более чем двумя знаками после неё", text)
	}
	padding := strings.Repeat("0", scoreDecimalPlaces-len(fraction))
	number, err := strconv.ParseInt(sign+whole+fraction+padding, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("число %q слишком велико по модулю. Введите меньшее число по модулю", text)
	}
	return number, nil
}

// validDecimal checks the digits and the length of the fraction.
func validDecimal(whole, fraction string, hasPoint bool) bool {
	if whole == "" && fraction == "" {
		return false
	}
	if hasPoint && (fraction == "" || len(fraction) > scoreDecimalPlaces) {
		return false
	}
	return strings.IndexFunc(whole, isNotDigit) < 0 && strings.IndexFunc(fraction, isNotDigit) < 0
}

// isNotDigit reports if the character is outside the range from 0 to 9.
func isNotDigit(char rune) bool {
	return char < '0' || char > '9'
}

// formatScore returns a score with the necessary decimal digits.
func formatScore(scoreHundredths int64) string {
	return strings.TrimSuffix(strings.TrimRight(formatHundredths(scoreHundredths), "0"), ".")
}

// formatHundredths returns a score with two decimal digits.
func formatHundredths(scoreHundredths int64) string {
	whole := scoreHundredths / hundredthsPerPoint
	fraction := scoreHundredths % hundredthsPerPoint
	if scoreHundredths < 0 {
		return fmt.Sprintf("-%d.%02d", -whole, -fraction)
	}
	return fmt.Sprintf("%d.%02d", whole, fraction)
}
