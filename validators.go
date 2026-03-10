package asky

import (
	"fmt"
	"math"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

// --- Input validators ------------------------------------

// ValidateInputChain combines multiple input validators into one.
// Validators are run in order — the first failure stops the chain.
//
//	asky.Input().WithValidator(asky.ValidateInputChain(
//	    asky.ValidateInputRequired(),
//	    asky.ValidateInputMinMaxLength(3, 50),
//	    asky.ValidateInputAlphanumeric(),
//	))
func ValidateInputChain(validators ...func(string) (string, bool)) func(string) (string, bool) {
	return func(s string) (string, bool) {
		for _, v := range validators {
			if msg, ok := v(s); !ok {
				return msg, false
			}
		}
		return "", true
	}
}

// ValidateInputRequired fails if the input is empty or whitespace only.
func ValidateInputRequired() func(string) (string, bool) {
	return func(s string) (string, bool) {
		if strings.TrimSpace(s) == "" {
			return "required", false
		}
		return "", true
	}
}

// ValidateInputMinLength fails if the input is shorter than n characters.
func ValidateInputMinLength(n int) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if len([]rune(s)) < n {
			return fmt.Sprintf("must be at least %d characters", n), false
		}
		return "", true
	}
}

// ValidateInputMaxLength fails if the input is longer than n characters.
func ValidateInputMaxLength(n int) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if len([]rune(s)) > n {
			return fmt.Sprintf("must be at most %d characters", n), false
		}
		return "", true
	}
}

// ValidateInputMinMaxLength fails if the input length is outside the range [min, max].
func ValidateInputMinMaxLength(min, max int) func(string) (string, bool) {
	return func(s string) (string, bool) {
		l := len([]rune(s))
		if l < min || l > max {
			return fmt.Sprintf("must be %d–%d characters", min, max), false
		}
		return "", true
	}
}

// ValidateInputEmail fails if the input is not a valid email address.
func ValidateInputEmail() func(string) (string, bool) {
	return func(s string) (string, bool) {
		if _, err := mail.ParseAddress(s); err != nil {
			return "must be a valid email address", false
		}
		return "", true
	}
}

// ValidateInputURL fails if the input is not a valid URL with a scheme and host.
func ValidateInputURL() func(string) (string, bool) {
	return func(s string) (string, bool) {
		u, err := url.ParseRequestURI(s)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return "must be a valid URL", false
		}
		return "", true
	}
}

// ValidateInputNumeric fails if the input contains any non-digit characters.
func ValidateInputNumeric() func(string) (string, bool) {
	return func(s string) (string, bool) {
		for _, r := range s {
			if !unicode.IsDigit(r) {
				return "must contain digits only", false
			}
		}
		return "", true
	}
}

// ValidateInputAlphanumeric fails if the input contains any non-alphanumeric characters.
func ValidateInputAlphanumeric() func(string) (string, bool) {
	return func(s string) (string, bool) {
		for _, r := range s {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return "must contain letters and digits only", false
			}
		}
		return "", true
	}
}

// ValidateInputRegex fails if the input does not match the given regular expression.
// msg is the error message shown to the user on failure.
func ValidateInputRegex(pattern, msg string) func(string) (string, bool) {
	re := regexp.MustCompile(pattern)
	return func(s string) (string, bool) {
		if !re.MatchString(s) {
			return msg, false
		}
		return "", true
	}
}

// ValidateInputMin fails if the input, parsed as a number, is less than n.
func ValidateInputMin(n float64) func(string) (string, bool) {
	return func(s string) (string, bool) {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return "must be a number", false
		}
		if v < n {
			return fmt.Sprintf("must be at least %s", formatFloat(n)), false
		}
		return "", true
	}
}

// ValidateInputMax fails if the input, parsed as a number, is greater than n.
func ValidateInputMax(n float64) func(string) (string, bool) {
	return func(s string) (string, bool) {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return "must be a number", false
		}
		if v > n {
			return fmt.Sprintf("must be at most %s", formatFloat(n)), false
		}
		return "", true
	}
}

// ValidateInputMinMax fails if the input, parsed as a number, is outside the range [min, max].
func ValidateInputMinMax(min, max float64) func(string) (string, bool) {
	return func(s string) (string, bool) {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return "must be a number", false
		}
		if v < min || v > max {
			return fmt.Sprintf("must be between %s and %s", formatFloat(min), formatFloat(max)), false
		}
		return "", true
	}
}

// ValidateInputIPAddr fails if the input is not a valid IPv4 or IPv6 address.
func ValidateInputIPAddr() func(string) (string, bool) {
	return func(s string) (string, bool) {
		if net.ParseIP(s) == nil {
			return "must be a valid IP address", false
		}
		return "", true
	}
}

// ValidateInputPortNumber fails if the input is not a valid port number (1–65535).
func ValidateInputPortNumber() func(string) (string, bool) {
	return func(s string) (string, bool) {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 65535 {
			return "must be a valid port number (1–65535)", false
		}
		return "", true
	}
}

// ValidateInputNoSpaces fails if the input contains any whitespace characters.
func ValidateInputNoSpaces() func(string) (string, bool) {
	return func(s string) (string, bool) {
		for _, r := range s {
			if unicode.IsSpace(r) {
				return "must not contain spaces", false
			}
		}
		return "", true
	}
}

// ValidateInputStartsWith fails if the input does not start with the given prefix.
func ValidateInputStartsWith(prefix string) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if !strings.HasPrefix(s, prefix) {
			return fmt.Sprintf("must start with %q", prefix), false
		}
		return "", true
	}
}

// ValidateInputEndsWith fails if the input does not end with the given suffix.
func ValidateInputEndsWith(suffix string) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if !strings.HasSuffix(s, suffix) {
			return fmt.Sprintf("must end with %q", suffix), false
		}
		return "", true
	}
}

// ValidateInputOneOf fails if the input does not exactly match one of the allowed values.
func ValidateInputOneOf(options ...string) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if slices.Contains(options, s) {
			return "", true
		}
		return fmt.Sprintf("must be one of: %s", strings.Join(options, ", ")), false
	}
}

// ValidateInputMatches fails if the input does not match the value pointed to by other.
// Useful for confirm-style prompts such as password confirmation.
//
//	pass, _ := asky.Input().WithLabel("Password").Render()
//	asky.Input().
//	    WithLabel("Confirm password").
//	    WithValidator(asky.ValidateInputMatches(&pass)).
//	    Render()
func ValidateInputMatches(other *string) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if s != *other {
			return "does not match", false
		}
		return "", true
	}
}

// --- Select validators ------------------------------------

// ValidateSelectRequired fails if no choice has been made (zero value Choice).
func ValidateSelectRequired() func(Choice) (string, bool) {
	return func(c Choice) (string, bool) {
		if c == (Choice{}) {
			return "required", false
		}
		return "", true
	}
}

// --- MultiSelect validators -------------------------------

// ValidateMultiSelectChain combines multiple MultiSelect validators into one.
// Validators are run in order — the first failure stops the chain.
//
//	asky.MultiSelect().WithValidator(asky.ValidateMultiSelectChain(
//	    asky.ValidateMultiSelectRequired(),
//	    asky.ValidateMultiSelectMinMax(1, 3),
//	))
func ValidateMultiSelectChain(validators ...func([]Choice) (string, bool)) func([]Choice) (string, bool) {
	return func(choices []Choice) (string, bool) {
		for _, v := range validators {
			if msg, ok := v(choices); !ok {
				return msg, false
			}
		}
		return "", true
	}
}

// ValidateMultiSelectRequired fails if no choices have been selected.
func ValidateMultiSelectRequired() func([]Choice) (string, bool) {
	return func(choices []Choice) (string, bool) {
		if len(choices) == 0 {
			return "required", false
		}
		return "", true
	}
}

// ValidateMultiSelectMin fails if fewer than n choices are selected.
func ValidateMultiSelectMin(n int) func([]Choice) (string, bool) {
	return func(choices []Choice) (string, bool) {
		if len(choices) < n {
			return fmt.Sprintf("select at least %d", n), false
		}
		return "", true
	}
}

// ValidateMultiSelectMax fails if more than n choices are selected.
func ValidateMultiSelectMax(n int) func([]Choice) (string, bool) {
	return func(choices []Choice) (string, bool) {
		if len(choices) > n {
			return fmt.Sprintf("select at most %d", n), false
		}
		return "", true
	}
}

// ValidateMultiSelectMinMax fails if the number of selected choices is outside [min, max].
func ValidateMultiSelectMinMax(min, max int) func([]Choice) (string, bool) {
	return func(choices []Choice) (string, bool) {
		n := len(choices)
		if n < min || n > max {
			return fmt.Sprintf("select between %d and %d", min, max), false
		}
		return "", true
	}
}

// formatFloat formats a float64 as an integer string if it has no fractional
// part, otherwise as a decimal string.
func formatFloat(f float64) string {
	if f == math.Trunc(f) {
		return strconv.Itoa(int(f))
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}
