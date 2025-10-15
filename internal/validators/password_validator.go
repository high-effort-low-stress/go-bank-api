package validators

import "regexp"

// Complexity criteria:
// - Between 8 and 64 characters long
// - At least one lowercase letter
// - At least one uppercase letter
// - At least one number
// - At least one special character

func ValidatePasswordPattern(plaintextPassword string) bool {

	if len(plaintextPassword) < 8 || len(plaintextPassword) > 64 {
		return false
	}

	hasLowercase, _ := regexp.MatchString(`[a-z]`, plaintextPassword)
	hasUppercase, _ := regexp.MatchString(`[A-Z]`, plaintextPassword)
	hasNumber, _ := regexp.MatchString(`\d`, plaintextPassword)
	hasSpecial, _ := regexp.MatchString(`[!@#$%^&*()-_=+]`, plaintextPassword)

	onlyAllowedChars, _ := regexp.MatchString(`^[a-zA-Z0-9!@#$%^&*()-_=+]+$`, plaintextPassword)

	if !hasLowercase || !hasUppercase || !hasNumber || !hasSpecial || !onlyAllowedChars {
		return false
	}

	return true
}
