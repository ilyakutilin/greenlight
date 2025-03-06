package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

// Constants for the token scope.
const (
	ScopeActivation = "activation"
)

// Token struct to hold the data for an individual token. This includes the plaintext
// and hashed versions of the token, associated user ID, expiry time and scope.
type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

// Generates a new token using a cryptographically secure random number generator
// (CSPRNG). Accepts a user ID, expiry duration (ttl / time-to-live), and a token scope
// as parameters.
func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	// Create a Token instance containing the user ID, expiry, and scope information.
	// We add the provided ttl (time-to-live) duration parameter to the current time
	// to get the expiry time.
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	// Initialize a zero-valued byte slice with a length of 16 bytes.
	// The resulting plaintext token strings (below) will not be 16 characters long -
	// but rather they have an underlying entropy of 16 bytes of randomness.
	randomBytes := make([]byte, 16)

	// Use the Read() function from the crypto/rand package to fill the byte slice with
	// random bytes from your operating system's CSPRNG. This will return an error if
	// the CSPRNG fails to function correctly.
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Encode the byte slice to a base-32-encoded string and assign it to the token
	// Plaintext field. This will be the token string that we send to the user in their
	// welcome email. They will look similar to this:
	//
	// Y3QMGX3PJ3WLRL2YRTQGQ6KRHU
	//
	// The length of the plaintext token string itself depends on how 16 random bytes
	// (above) are encoded to create a string. In our case we encode the random bytes to
	// a base-32 string, which results in a string with 26 characters. In contrast, if
	// we encoded the random bytes using hexadecimal (base-16) the string would be 32
	// characters long instead.
	// Note that by default base-32 strings may be padded at the end with the =
	// character. We don't need this padding character for the purpose of our tokens, so
	// we use the WithPadding(base32.NoPadding) method in the line below to omit them.
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).
		EncodeToString(randomBytes)

	// Generate a SHA-256 hash of the plaintext token string. This will be the value
	// that we store in the `hash` field of our database table. Note that the
	// sha256.Sum256() function returns an *array* of length 32, so to make it easier to
	// work with we convert it to a slice using the [:] operator before storing it.
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}
