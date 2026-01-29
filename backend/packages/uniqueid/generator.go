package uniqueid

import (
	"crypto/rand"
	"io"
)

// The character set for our unique IDs.
const idCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const idLength = 14

// New creates a new 14-character, uppercase, alphanumeric unique ID.
// It uses a cryptographically secure random number generator.
func New() string {
	// Create a byte slice to hold the random data.
	randomBytes := make([]byte, idLength)

	// Read random data into the byte slice.
	// We use io.ReadFull to ensure we get the exact number of bytes we need.
	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		// This is a critical error. If the OS can't provide random data,
		// we can't generate a safe ID. We panic to stop the application.
		panic(err)
	}

	// Create a string builder for efficiency.
	idChars := make([]byte, idLength)
	for i := 0; i < idLength; i++ {
		// Map the random byte to a character in our charset.
		// The modulo operator ensures the index is within the bounds of the charset.
		idChars[i] = idCharset[int(randomBytes[i])%len(idCharset)]
	}

	return string(idChars)
}
