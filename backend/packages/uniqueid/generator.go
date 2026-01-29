package uniqueid

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const (
	idLength         = 14
	timePartLength   = 9 // Length of a Unix millisecond timestamp encoded in Base36.
	randomPartLength = 5 // 14 - 9 = 5
)

// New creates a new 14-character, uppercase, alphanumeric, time-ordered unique ID.
func New() string {
	// 1. Time Component
	// Get the current time in milliseconds since the epoch.
	now := time.Now().UnixMilli()
	// Convert the timestamp to a Base36 string and uppercase it.
	timePart := strings.ToUpper(strconv.FormatInt(now, 36))

	// Pad the time part if necessary (for future-proofing).
	if len(timePart) < timePartLength {
		timePart = strings.Repeat("0", timePartLength-len(timePart)) + timePart
	}

	// 2. Random Component
	// Calculate the maximum value for the random part (36^5).
	maxRandom := new(big.Int)
	maxRandom.Exp(big.NewInt(36), big.NewInt(randomPartLength), nil)

	// Generate a cryptographically secure random number within the calculated range.
	randomInt, err := rand.Int(rand.Reader, maxRandom)
	if err != nil {
		// This is a critical failure. The OS cannot provide random data.
		panic(fmt.Sprintf("failed to generate random number for unique ID: %v", err))
	}

	// Convert the random number to a Base36 string and uppercase it.
	randomPart := strings.ToUpper(randomInt.Text(36))

	// 3. Padding
	// Pad the random part with leading zeros to ensure it's exactly `randomPartLength` characters.
	paddedRandomPart := strings.Repeat("0", randomPartLength-len(randomPart)) + randomPart

	// 4. Concatenation
	// Combine the time and random parts to create the final 14-character ID.
	return timePart + paddedRandomPart
}
