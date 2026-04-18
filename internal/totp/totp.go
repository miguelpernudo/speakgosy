// Package totp provides utilities to generate and validate the passwords.
package totp

import (
	"time"
	ptotp "github.com/pquerna/otp/totp"
)

// Generate creates the current TOTP code for the given base32 secret using the system's current time.
func Generate(secret string) (string, error) {
	return ptotp.GenerateCode(secret, time.Now())
}

// Validate verifies a TOTP code against a secret, allowing a skew of one 30s window to account for clock drift.
func Validate(code, secret string) (bool, error) {
	opts := ptotp.ValidateOpts{
		Skew: 1,
	}
	return pquerna.ValidateCustom(code, secret, time.Now(), opts)
}
