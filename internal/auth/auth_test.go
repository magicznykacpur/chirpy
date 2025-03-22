package auth

import (
	"testing"
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
