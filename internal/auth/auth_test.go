package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashingPasswords(t *testing.T) {
	cases := []struct {
		input string
	}{
		{
			input: "myPassword",
		},
		{
			input: "myOtherPassword",
		},
		{
			input: "yetAnotherPassword",
		},
	}

	for _, c := range cases {
		hash, err := HashPassword(c.input)
		if err != nil {
			t.Errorf("hashing broken: %v", err)
		}

		err = CheckPasswordHash(hash, c.input)
		if err != nil {
			t.Errorf("checking hashes broken: %v", err)
		}
	}
}

func TestJWTToken(t *testing.T) {
	testId, _ := uuid.Parse("b592343d-b059-4d87-a1db-69d3c8accccf")
	secret := "very-secret-secret"

	token, err := MakeJWT(testId, secret, time.Second * 3)
	if err != nil {
		t.Errorf("cannot create token: %v", err)
	}

	id, err := ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("cannot validate token: %v", err)
	}

	if id != testId {
		t.Errorf("validation failed, id's do not match: --> %v != %v <--", id, testId)
	}

	_, err = ValidateJWT(token, secret + "secret")
	if err == nil {
		t.Errorf("cannot validate token with wrong secret: %v", err)
	}
}

func TestJWTTokenExpiration(t *testing.T) {
	testId, _ := uuid.Parse("b592343d-b059-4d87-a1db-69d3c8accccf")
	secret := "very-secret-secret"

	token, err := MakeJWT(testId, secret, time.Millisecond * 2)
	if err != nil {
		t.Errorf("cannot create token: %v", err)
	}

	time.Sleep(time.Millisecond * 3)

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("token should be expired")
	}
}