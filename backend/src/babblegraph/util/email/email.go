package email

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func ValidateEmailAddress(emailAddress string) error {
	if len(emailAddress) < 3 || len(emailAddress) > 254 {
		return fmt.Errorf("Email Address is too long or too short")
	}
	if !emailRegex.MatchString(emailAddress) {
		return fmt.Errorf("Invalid email address formatting")
	}
	emailDomain := strings.Split(emailAddress, "@")[1]
	mx, err := net.LookupMX(emailDomain)
	if err != nil {
		return err
	}
	if len(mx) == 0 {
		return fmt.Errorf("No mx record")
	}
	return nil
}

func FormatEmailAddress(emailAddress string) string {
	return strings.ToLower(strings.Trim(emailAddress, " "))
}
