package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	memory      = 64 * 1024
	iterations  = 3
	parallelism = 2
	saltLength  = 16
	keyLength   = 32
)

// HashPassword generates an Argon2id hash in the standard encoded format.
func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	// Standard Argon2id encoded form:
	// $argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>
	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		memory,
		iterations,
		parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encoded, nil
}

// CheckPassword verifies a password against the encoded hash.
func CheckPassword(password, encoded string) error {
	version, p, salt, expected, err := decodeHash(encoded)
	if err != nil {
		return err
	}

	if version != argon2.Version {
		return fmt.Errorf("argon2 version mismatch: got %d want %d", version, argon2.Version)
	}

	actual := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	if subtle.ConstantTimeCompare(actual, expected) != 1 {
		return errors.New("password does not match")
	}

	return nil
}

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

// decodeHash parses the standard Argon2id encoded hash.
func decodeHash(encoded string) (version int, p *params, salt, hash []byte, err error) {
	parts := strings.Split(encoded, "$")
	// Expected:
	// ["", "argon2id", "v=19", "m=...,t=...,p=...", "<salt>", "<hash>"]
	if len(parts) != 6 {
		return 0, nil, nil, nil, fmt.Errorf("invalid encoded hash")
	}

	if parts[1] != "argon2id" {
		return 0, nil, nil, nil, fmt.Errorf("unsupported algorithm: %s", parts[1])
	}

	// version
	if _, err = fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return 0, nil, nil, nil, fmt.Errorf("invalid version field")
	}

	p = &params{}
	if _, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism); err != nil {
		return 0, nil, nil, nil, fmt.Errorf("invalid parameter field: %w", err)
	}

	// decode salt
	if salt, err = base64.RawStdEncoding.DecodeString(parts[4]); err != nil {
		return 0, nil, nil, nil, fmt.Errorf("invalid salt")
	}

	// decode hash
	if hash, err = base64.RawStdEncoding.DecodeString(parts[5]); err != nil {
		return 0, nil, nil, nil, fmt.Errorf("invalid hash")
	}

	p.keyLength = uint32(len(hash))

	return version, p, salt, hash, nil
}
