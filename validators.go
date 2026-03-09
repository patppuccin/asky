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

// ValidatorChain combines multiple validators into one.
// Validators are run in order — the first failure stops the chain.
//
//	asky.Input().WithValidator(asky.ValidatorChain(
//	    asky.ValidateRequired(),
//	    asky.ValidateMinMaxLength(3, 50),
//	    asky.ValidateAlphanumeric(),
//	))
func ValidatorChain(validators ...func(string) (string, bool)) func(string) (string, bool) {
	return func(s string) (string, bool) {
		for _, v := range validators {
			if msg, ok := v(s); !ok {
				return msg, false
			}
		}
		return "", true
	}
}

// ValidateRequired fails if the input is empty or whitespace only.
func ValidateRequired() func(string) (string, bool) {
	return func(s string) (string, bool) {
		if strings.TrimSpace(s) == "" {
			return "required", false
		}
		return "", true
	}
}

// ValidateMinLength fails if the input is shorter than n characters.
func ValidateMinLength(n int) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if len([]rune(s)) < n {
			return fmt.Sprintf("must be at least %d characters", n), false
		}
		return "", true
	}
}

// ValidateMaxLength fails if the input is longer than n characters.
func ValidateMaxLength(n int) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if len([]rune(s)) > n {
			return fmt.Sprintf("must be at most %d characters", n), false
		}
		return "", true
	}
}

// ValidateMinMaxLength fails if the input length is outside the range [min, max].
func ValidateMinMaxLength(min, max int) func(string) (string, bool) {
	return func(s string) (string, bool) {
		l := len([]rune(s))
		if l < min || l > max {
			return fmt.Sprintf("must be %d–%d characters", min, max), false
		}
		return "", true
	}
}

// ValidateEmail fails if the input is not a valid email address.
func ValidateEmail() func(string) (string, bool) {
	return func(s string) (string, bool) {
		if _, err := mail.ParseAddress(s); err != nil {
			return "must be a valid email address", false
		}
		return "", true
	}
}

// ValidateURL fails if the input is not a valid URL with a scheme and host.
func ValidateURL() func(string) (string, bool) {
	return func(s string) (string, bool) {
		u, err := url.ParseRequestURI(s)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return "must be a valid URL", false
		}
		return "", true
	}
}

// ValidateNumeric fails if the input contains any non-digit characters.
func ValidateNumeric() func(string) (string, bool) {
	return func(s string) (string, bool) {
		for _, r := range s {
			if !unicode.IsDigit(r) {
				return "must contain digits only", false
			}
		}
		return "", true
	}
}

// ValidateAlphanumeric fails if the input contains any non-alphanumeric characters.
func ValidateAlphanumeric() func(string) (string, bool) {
	return func(s string) (string, bool) {
		for _, r := range s {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return "must contain letters and digits only", false
			}
		}
		return "", true
	}
}

// ValidateRegex fails if the input does not match the given regular expression.
// msg is the error message shown to the user on failure.
func ValidateRegex(pattern, msg string) func(string) (string, bool) {
	re := regexp.MustCompile(pattern)
	return func(s string) (string, bool) {
		if !re.MatchString(s) {
			return msg, false
		}
		return "", true
	}
}

// ValidateMin fails if the input, parsed as a number, is less than n.
func ValidateMin(n float64) func(string) (string, bool) {
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

// ValidateMax fails if the input, parsed as a number, is greater than n.
func ValidateMax(n float64) func(string) (string, bool) {
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

// ValidateMinMax fails if the input, parsed as a number, is outside the range [min, max].
func ValidateMinMax(min, max float64) func(string) (string, bool) {
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

// formatFloat formats a float64 as an integer string if it has no fractional
// part, otherwise as a decimal string.
func formatFloat(f float64) string {
	if f == math.Trunc(f) {
		return strconv.Itoa(int(f))
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// ValidateIPAddr fails if the input is not a valid IPv4 or IPv6 address.
func ValidateIPAddr() func(string) (string, bool) {
	return func(s string) (string, bool) {
		if net.ParseIP(s) == nil {
			return "must be a valid IP address", false
		}
		return "", true
	}
}

// ValidatePortNumber fails if the input is not a valid port number (1–65535).
func ValidatePortNumber() func(string) (string, bool) {
	return func(s string) (string, bool) {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 65535 {
			return "must be a valid port number (1–65535)", false
		}
		return "", true
	}
}

// ValidateNoSpaces fails if the input contains any whitespace characters.
func ValidateNoSpaces() func(string) (string, bool) {
	return func(s string) (string, bool) {
		for _, r := range s {
			if unicode.IsSpace(r) {
				return "must not contain spaces", false
			}
		}
		return "", true
	}
}

// ValidateStartsWith fails if the input does not start with the given prefix.
func ValidateStartsWith(prefix string) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if !strings.HasPrefix(s, prefix) {
			return fmt.Sprintf("must start with %q", prefix), false
		}
		return "", true
	}
}

// ValidateEndsWith fails if the input does not end with the given suffix.
func ValidateEndsWith(suffix string) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if !strings.HasSuffix(s, suffix) {
			return fmt.Sprintf("must end with %q", suffix), false
		}
		return "", true
	}
}

// ValidateOneOf fails if the input does not exactly match one of the allowed values.
func ValidateOneOf(options ...string) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if slices.Contains(options, s) {
			return "", true
		}
		return fmt.Sprintf("must be one of: %s", strings.Join(options, ", ")), false
	}
}

// ValidateMatches fails if the input does not match the value pointed to by other.
// Useful for confirm-style prompts such as password confirmation.
//
//	pass, _ := asky.InputSecret().WithLabel("Password").Render()
//	asky.InputSecret().
//	    WithLabel("Confirm password").
//	    WithValidator(asky.ValidateMatches(&pass)).
//	    Render()
func ValidateMatches(other *string) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if s != *other {
			return "does not match", false
		}
		return "", true
	}
}
