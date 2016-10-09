package librevault

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"

	"github.com/jbenet/go-base58"
	"github.com/mastercactapus/librevault/keccak"
)

// Secret is a key, used to identify sync folders, perform encryption and peer discovery.
type Secret struct {
	Type    SecretType
	Payload []byte
}

// SecretType specifies the type of secret. It is the first byte of a Secret
type SecretType byte

// A secret can be one of the following types
const (
	SecretTypeOwner        SecretType = 'A'
	SecretTypeReadOnly     SecretType = 'C'
	SecretTypeDownloadOnly SecretType = 'D'
)

// MarshalJSON will encode data to a JSON string
func (s *Secret) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON will decode data from a JSON string
func (s *Secret) UnmarshalJSON(data []byte) error {
	str := new(string)
	err := json.Unmarshal(data, str)
	if err != nil {
		return err
	}
	sec, err := ParseSecret(*str)
	if err != nil {
		return err
	}
	*s = *sec
	return nil
}

// ParseSecret will parse and validate a secret string
func ParseSecret(s string) (*Secret, error) {
	if len(s) < 3 {
		return nil, io.ErrShortBuffer
	}
	if validateChecksum(s) {
		return nil, errors.New("invalid checksum")
	}
	if s[1] != '1' {
		return nil, errors.New("invalid param value")
	}

	payload := base58.Decode(s[2 : len(s)-1])
	switch SecretType(s[0]) {
	case SecretTypeOwner, SecretTypeDownloadOnly:
		if len(payload) != 32 {
			return nil, errors.New("invalid length for secret")
		}
	case SecretTypeReadOnly:
		if len(payload) != 65 {
			return nil, errors.New("invalid length for secret")
		}
	}

	return &Secret{Type: SecretType(s[0]), Payload: payload}, nil
}

// NewSecret will generate a new Owner-level secret
func NewSecret() (*Secret, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Secret{Type: SecretTypeOwner, Payload: key.D.Bytes()}, nil
}

func (s Secret) String() string {
	return string(s.Type) + "1" + appendChecksum(base58.Encode(s.Payload))
}

// Owner will return the Owner string if Secret is of type SecretTypeOwner, otherwise an empty string
func (s Secret) Owner() string {
	if s.Type != SecretTypeOwner {
		return ""
	}
	return s.String()
}

// ReadOnly will return the ReadOnly string if Secret is of type SecretTypeOwner, or SecretTypeReadOnly, otherwise an empty string
func (s Secret) ReadOnly() string {
	switch s.Type {
	case SecretTypeOwner:
		// s.Payload is the private key, aka the exponent. we can use it directly
		// to derive the X and Y coords (for public key)
		x, y := elliptic.P256().ScalarBaseMult(s.Payload)
		payload := make([]byte, 65)
		// the public key is point compressed (X coord with indicator for even/odd Y coord)
		// 0x02 for even, 0x03 for odd. point compression, something something because ASN.1 happened
		payload[0] = 2 + byte(y.Bit(0))
		copy(payload[1:], x.Bytes())

		// readOnly has the hash of the private key appended to the end
		hash := keccak.Sum256(s.Payload)
		copy(payload[33:], hash[:])
		return Secret{Type: SecretTypeReadOnly, Payload: payload}.String()
	case SecretTypeReadOnly:
		return s.String()
	}
	return ""
}

// DownloadOnly will return the DownloadOnly string for the Secret
func (s Secret) DownloadOnly() string {
	switch s.Type {
	case SecretTypeOwner:
		hash := keccak.Sum256(s.Payload)
		return Secret{Type: SecretTypeDownloadOnly, Payload: hash[:]}.String()
	case SecretTypeReadOnly:
		return Secret{Type: SecretTypeDownloadOnly, Payload: s.Payload[33:]}.String()
	case SecretTypeDownloadOnly:
		return s.String()
	}
	return ""
}
