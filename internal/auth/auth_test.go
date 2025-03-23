package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	headerNoAuth := http.Header{}
	headerNoAuth.Add("Content-Type", "application/json")

	_, err := GetBearerToken(headerNoAuth)
	if err.Error() != "authorization token missing from header" {
		t.Errorf("get bearer token should not detect missing authorization")
	}

	bearerToken := "this-is-my-token"
	headerWithAuth := http.Header{}
	headerWithAuth.Add("Authorization", "Bearer " + bearerToken)

	token, err := GetBearerToken(headerWithAuth)
	if err != nil {
		t.Errorf("bearer token not detected: %v", err)
	}

	if token != bearerToken {
		t.Errorf("token mismatch --> %s != %s", token, bearerToken)
	}
}

func TestMakeRefreshToken(t *testing.T) {
	refreshToken, err := MakeRefreshToken()
	
	if err != nil {
		t.Errorf("make refresh token shouldn't return an error")
	}

	if len(refreshToken) != 64 {
		t.Errorf("refresh token should be of lenght 64, but is %d", len(refreshToken))
	}
}