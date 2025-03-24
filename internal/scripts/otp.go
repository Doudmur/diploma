package scripts

import (
	"math/rand"
	"time"
)

func GenerateOTP() string {
	// Create a new random number generator with a unique seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	digits := "0123456789"
	otp := ""
	for i := 0; i < 6; i++ {
		otp += string(digits[r.Intn(len(digits))])
	}
	return otp
}
