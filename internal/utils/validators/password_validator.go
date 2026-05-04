// Package validators provides functions for validating user input data, such as CPF and password complexity.
package validators

import "regexp"

// Complexity criteria:
// - Between 8 and 64 characters long
// - At least one lowercase letter
// - At least one uppercase letter
// - At least one number
// - At least one special character

var (
	reLowercase    = regexp.MustCompile(`[a-z]`)
	reUppercase    = regexp.MustCompile(`[A-Z]`)
	reNumber       = regexp.MustCompile(`\d`)
	reSpecial      = regexp.MustCompile(`[!@#$%^&*()_=+ .-]`)
	reAllowedChars = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_=+ .-]+$`)
)

func ValidatePasswordPattern(plaintextPassword string) bool {
	if len(plaintextPassword) < 8 || len(plaintextPassword) > 64 {
		return false
	}

	hasLowercase := reLowercase.MatchString(plaintextPassword)
	hasUppercase := reUppercase.MatchString(plaintextPassword)
	hasNumber := reNumber.MatchString(plaintextPassword)
	hasSpecial := reSpecial.MatchString(plaintextPassword)

	onlyAllowedChars := reAllowedChars.MatchString(plaintextPassword)

	if !hasLowercase || !hasUppercase || !hasNumber || !hasSpecial || !onlyAllowedChars {
		return false
	}

	return true
}
