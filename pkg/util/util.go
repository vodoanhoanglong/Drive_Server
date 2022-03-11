package util

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	alphaDigits   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var webBrowserRegex = regexp.MustCompile(`(?i)(opera|chrome|safari|firefox|msie|trident)[/\s]([\d.]+)`)

var src = rand.NewSource(time.Now().UnixNano())

func JSONDecode(r io.Reader, val interface{}) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(val)
}

func IsEmail(email string, isStrict bool) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	if !emailRegex.MatchString(email) {
		return false
	}

	// validate valid domain if isStrict = true
	if isStrict {
		parts := strings.Split(email, "@")
		mx, err := net.LookupMX(parts[1])
		if err != nil || len(mx) == 0 {
			return false
		}
	}
	return true

}

func EncodePostgresUUIDs(inputs []uuid.UUID) string {
	length := len(inputs)
	if length == 0 {
		return "{}"
	}

	sInputs := make([]string, length)
	for i, u := range inputs {
		sInputs[i] = u.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(sInputs, ","))
}

// HasUUID check if the array has input uuid item
func HasUUID(list []uuid.UUID, item uuid.UUID) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}

	return false
}

// HasString check if the array has input string item
func HasString(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}

	return false
}

// GenRandomString generate random string with fixed length
func GenRandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(alphaDigits) {
			sb.WriteByte(alphaDigits[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

// GenID generate random unique id
func GenID() string {
	now := time.Now()
	return strconv.FormatInt(int64(now.Year()), 32) +
		string(alphaDigits[int(now.Month())]) +
		string(alphaDigits[now.Day()]) +
		string(alphaDigits[now.Hour()]) +
		string(alphaDigits[now.Minute()]) +
		string(alphaDigits[now.Second()]) +
		GenRandomString(8)
}

// EncodePostgresArray encode array string to postgres array
func EncodePostgresArray(input []string) string {
	return fmt.Sprintf("{%s}", strings.Join(input, ","))
}

// DecodePostgresArray decode postgres array string
func DecodePostgresArray(input string) ([]string, error) {
	if input == "{}" {
		return []string{}, nil
	}

	if len(input) < 3 || input[0] != '{' || input[len(input)-1] != '}' {
		return nil, fmt.Errorf("invalid postgres array: %s", input)
	}

	return strings.Split(input[1:len(input)-1], ","), nil
}

// DurationToMilliseconds convert duration to milliseconds
func DurationToMilliseconds(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}

// GetRequestIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetRequestIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// GetRequestHost gets a requests host by reading off the forwarded-host
// header (for proxies) and falls back to use the remote address.
func GetRequestHost(r *http.Request) string {
	host := r.Header.Get("X-Forwarded-Host")
	port := r.Header.Get("X-Forwarded-Port")

	if host != "" && (port == "" || port == "80" || port == "443") {
		return host
	}

	return fmt.Sprintf("%s:%s", host, port)
}

// GetRequestOrigin get request origin from x-forwarded-origin header
func GetRequestOrigin(r *http.Request) string {
	return r.Header.Get("X-Forwarded-Origin")
}

func intersectionStrings2(src []string, target []string) []string {
	result := []string{}
	for _, s := range src {
		if HasString(target, s) && !HasString(result, s) {
			result = append(result, s)
		}
	}

	return result
}

// IntersectionStrings return common values between 2 array strings
func IntersectionStrings(src []string, target []string, others ...[]string) []string {
	res := intersectionStrings2(src, target)
	if len(others) > 0 {
		for _, ss := range others {
			res = intersectionStrings2(res, ss)
		}
	}

	return res
}

// GetFirstStringInMap get first string value in string map
func GetFirstStringInMap(input map[string]string) string {
	for _, v := range input {
		return v
	}
	return ""
}

// IsWebBrowserAgent checks if the user agent is from web browser
func IsWebBrowserAgent(userAgent string) bool {
	if userAgent == "" {
		return false
	}

	return webBrowserRegex.MatchString(userAgent)
}
