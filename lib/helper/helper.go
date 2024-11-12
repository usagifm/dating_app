package helper

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Contains(arr []int, target int) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

func GetTodayStartAndEnd() (time.Time, time.Time) {
	now := time.Now()

	// Start of the day (00:00 AM)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// End of the day (23:59 PM, 59 seconds, 999999999 nanoseconds)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	return startOfDay, endOfDay
}

// HashPassword hashes the given password string
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// ComparePassword compares a hashed password with a plaintext password
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil // Returns true if the password is correct
}
