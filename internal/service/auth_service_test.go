package service

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordCompare(t *testing.T) {
	hash := "$2a$12$rUE6SmL3T13tyZu/fZTIHut.LNOa2vKjG9WDiGjR5kWQ.tYXKbcVa"
	password := "12345"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		t.Error("❌ Password does NOT match")
	} else {
		t.Log("✅ Password MATCHES")
	}
}
