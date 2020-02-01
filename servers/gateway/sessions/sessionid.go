package sessions

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

//InvalidSessionID represents an empty, invalid session ID
const InvalidSessionID SessionID = ""

//idLength is the length of the ID portion
const idLength = 32

//signedLength is the full length of the signed session ID
//(ID portion plus signature)
const signedLength = idLength + sha256.Size

//SessionID represents a valid, digitally-signed session ID.
//This is a base64 URL encoded string created from a byte slice
//where the first `idLength` bytes are crytographically random
//bytes representing the unique session ID, and the remaining bytes
//are an HMAC hash of those ID bytes (i.e., a digital signature).
//The byte slice layout is like so:
//+-----------------------------------------------------+
//|...32 crypto random bytes...|HMAC hash of those bytes|
//+-----------------------------------------------------+
type SessionID string

//ErrInvalidID is returned when an invalid session id is passed to ValidateID()
var ErrInvalidID = errors.New("Invalid Session ID")

// GenRandomBytes generates random byte string of given size
func GenRandomBytes(size int) (byteS []byte, err error) {
	byteS = make([]byte, size)
	_, err = rand.Read(byteS)
	return
}

//NewSessionID creates and returns a new digitally-signed session ID,
//using `signingKey` as the HMAC signing key. An error is returned only
//if there was an error generating random bytes for the session ID
func NewSessionID(signingKey string) (SessionID, error) {
	if len(signingKey) == 0 {
		return InvalidSessionID, fmt.Errorf("Signing key may not be empty")
	}

	byteString, err := GenRandomBytes(32)
	hasher := hmac.New(sha256.New, []byte(signingKey))
	hasher.Write(byteString)
	signature := hasher.Sum(nil)
	byteSession := append(byteString[:], signature[:]...)
	sessionID := base64.URLEncoding.EncodeToString(byteSession)
	sessionObjID := SessionID(sessionID)

	return sessionObjID, err
}

//ValidateID validates the string in the `id` parameter
//using the `signingKey` as the HMAC signing key
//and returns an error if invalid, or a SessionID if valid
func ValidateID(id string, signingKey string) (SessionID, error) {
	decode, _ := base64.URLEncoding.DecodeString(id)
	idPortion := decode[0:32]

	// Get the hmac hash from the given id
	hasher := hmac.New(sha256.New, []byte(signingKey))
	hasher.Write(idPortion)
	signature := hasher.Sum(nil)

	existingSignature := decode[32:]

	res := bytes.Compare(signature, existingSignature)
	if res == 0 {
		return SessionID(id), nil
	}
	return InvalidSessionID, ErrInvalidID
}

//String returns a string representation of the sessionID
func (sid SessionID) String() string {
	return string(sid)
}
