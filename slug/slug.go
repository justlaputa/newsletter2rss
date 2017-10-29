package slug

import (
	"math/rand"
	"strings"
	"time"
)

// New generate a new slug string based on input name
// check the exists function, make sure the generated slug always return false
// while passed to the exists function
func New(title string, exists func(string) bool) string {
	title = strings.TrimSpace(title)
	title = strings.ToLower(title)
	title = strings.Join(strings.Fields(title), "-")

	result := title
	for exists(result) {
		result = title + "-" + randomSuffix()
	}

	return result
}

func randomSuffix() string {
	rand.Seed(time.Now().UTC().UnixNano())

	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	length := 6
	buffer := make([]byte, length)

	for i := 0; i < length; i++ {
		buffer[i] = chars[rand.Intn(len(chars))]
	}

	return string(buffer)
}
