package cleaner

import "testing"

func TestCleanBody(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "this is a badword string",
			expected: "this is a **** string",
		},
		{
			input:    "this is a Badword string",
			expected: "this is a **** string",
		},
		{
			input:    "this is a BADWORD string",
			expected: "this is a **** string",
		},
	}

	for _, c := range cases {
		actual := CleanBodyBy(c.input, "badword")
		if actual != c.expected {
			t.Errorf("actual doesn't match expected --> %s != %s <--", actual, c.expected)
		}
	}

	casesMulti := []struct {
		input    string
		expected string
	}{
		{
			input:    "this is a badword string with another sillyword",
			expected: "this is a **** string with another ****",
		},
		{
			input:    "this is a Badword string with another Sillyword",
			expected: "this is a **** string with another ****",
		},
		{
			input:    "this is a BADWORD string with another SILLYWORD",
			expected: "this is a **** string with another ****",
		},
	}

	for _, c := range casesMulti {
		actual := CleanBodyBy(c.input, "badword")
		actual = CleanBodyBy(actual, "sillyword")
		if actual != c.expected {
			t.Errorf("actual doesn't match expected --> %s != %s <--", actual, c.expected)
		}
	}
}
