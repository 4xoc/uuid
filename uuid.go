// Package uuid provides functions to create and manage custom UUIDs.
//
// It uses random data to generate the ID but uses the first 6 bits to set its type.
package uuid

import (
	crand "crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	mrand "math/rand"
	"regexp"
	"strings"
)

// Type UUID holds the ID's information like the Scope as well as a hex string and binary representation.
type UUID struct {
	// scope describes the UUIDs scope as a string given on creation.
	scope string
	// Hex contains the canonical writing of the UUID
	hex string
	// bin contains the actual binary UUID.
	bin [16]byte
}

const (
	// Errors
	ErrorMissingScope      string = "The provided scope is not known."
	ErrorBadScope          string = "The provided scope is not supported."
	ErrorBadString         string = "The provided string is not a UUID."
	ErrorOutOfScopes       string = "Limit of scopes exeeded."
	ErrorMalformattedHex   string = "The Hex representation of the UUID is malformatted."
	ErrorUninitializedUUID string = "The provided pointer refers to an uninitialized struct."
	ErrorScopesAlreadySet  string = "Scopes can only be set once."
)

var (
	// setScopes holds the mapping between existing scopes (identified map index 'string')
	// and a pointer to the byte set in `scopes`.
	setScopes map[string]*byte

	// scopes holds a list of all available bytes that can be used to set the binary scope.
	scopes = [64]byte{
		0x00, 0x04, 0x08, 0x0c,
		0x10, 0x14, 0x18, 0x1c,
		0x20, 0x24, 0x28, 0x2c,
		0x30, 0x34, 0x38, 0x3c,
		0x40, 0x44, 0x48, 0x4c,
		0x50, 0x54, 0x58, 0x5c,
		0x60, 0x64, 0x68, 0x6c,
		0x70, 0x74, 0x78, 0x7c,
		0x80, 0x84, 0x88, 0x8c,
		0x90, 0x94, 0x98, 0x9c,
		0xa0, 0xa4, 0xa8, 0xac,
		0xb0, 0xb4, 0xb8, 0xbc,
		0xc0, 0xc4, 0xc8, 0xcc,
		0xd0, 0xd4, 0xd8, 0xdc,
		0xe0, 0xe4, 0xe8, 0xec,
		0xf0, 0xf4, 0xf8, 0xfc}
)

// Scope returns the scope of a UUID as a string.
//
// This function is needed to have a private 'scope' variable in the
// struct. Errors where a scope has been manually changed should
// be prevented by this.
func (uuid *UUID) Scope() string {
	if uuid == nil {
		return ""
	}

	return uuid.scope
}

// Bin returns the binary representation of a given UUID.
func (uuid *UUID) Bin() [16]byte {
	var (
		panicHealer [16]byte
	)

	if uuid == nil {
		return panicHealer
	}

	return uuid.bin
}

// Hex returns the hex-string representation of a given UUID.
func (uuid *UUID) Hex() string {
	if uuid == nil {
		return ""
	}

	return uuid.hex
}

// ScopeMatches checks a given slice of strings to check
// a for a matching scope. If any of the given scopes matches,
// the function returns true.
//
// If the UUID is not initialized, false is returned.
func (uuid *UUID) ScopeMatches(scopes []string) bool {
	var (
		index int
	)

	for index = range scopes {
		if uuid.scope == scopes[index] {
			return true
		}
	}

	return false
}

// readScope is a function that checks the binary data of the uuid and
// defines the scope as sting for that uuid.
func (uuid *UUID) readScope() error {
	var (
		tmpBytes []byte
		tmpByte  byte
		scope    string
		err      error
	)

	//first we set bin from hex
	tmpBytes, err = hex.DecodeString(strings.Replace(uuid.hex, "-", "", -1))
	if err != nil {
		return errors.New(ErrorBadString)
	}

	copy(uuid.bin[:], tmpBytes)

	//reading first byte and clearing last two bits
	tmpByte = uuid.bin[0] &^ 0x03

	if setScopes == nil {
		return errors.New(ErrorMissingScope)
	}

	for scope = range setScopes {
		if tmpByte == *setScopes[scope] {
			uuid.scope = scope
			break
		}
	}

	if uuid.scope == "" {
		return errors.New(ErrorBadScope)
	}

	return nil
}

// Value provides a database/sql/driver interface to read the struct's value and pass it to a DB connection.
func (uuid UUID) Value() (driver.Value, error) {
	if len(uuid.hex) != 36 {
		return nil, errors.New(ErrorMalformattedHex)
	}

	return uuid.hex, nil
}

// Scan provides a database/sql/driver interface to read the data coming from a DB connection into a struct.
func (uuid *UUID) Scan(src interface{}) error {
	var (
		ok      bool
		tmpByte []byte
	)

	if tmpByte, ok = src.([]byte); !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	uuid.hex = fmt.Sprintf("%x-%x-%x-%x-%x",
		tmpByte[0:4],
		tmpByte[4:6],
		tmpByte[6:8],
		tmpByte[8:10],
		tmpByte[10:16])

	//returns nil if uuid is good or error if the is a problem
	return uuid.readScope()
}

// New generates a new UUID and sets its scope to the one provided as an argument.
// If the scope doesn't exist yet, it will return an error (see SetScopes function).
func New(scope string) (*UUID, error) {
	var (
		uuid UUID
		err  error
	)

	if setScopes[scope] == nil {
		return nil, errors.New(ErrorMissingScope)
	}

	_, err = crand.Read(uuid.bin[:])

	if err != nil {
		return nil, errors.New("Error generating new UUID: " + err.Error())
	}

	//set scope
	uuid.bin[0] = *setScopes[scope] | byte(mrand.Intn(4))
	uuid.scope = scope

	//formatting as canonical string
	uuid.hex = fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid.bin[0:4],
		uuid.bin[4:6],
		uuid.bin[6:8],
		uuid.bin[8:10],
		uuid.bin[10:16])

	return &uuid, nil
}

// Read uses a given string and parses it into a UUID struct.
func Read(input string) (*UUID, error) {
	var (
		uuid UUID
		err  error
	)

	if !regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$").MatchString(input) {
		return nil, errors.New(ErrorBadString)
	}

	uuid.hex = input

	err = uuid.readScope()
	if err != nil {
		return nil, errors.New(ErrorBadScope)
	}

	return &uuid, nil
}

// Scopes provides a list of all currently set scopes in a [64]string. The order is not the same as set with
// SetScopes function.
func Scopes() [64]string {
	var (
		scope  string
		scopes [64]string
		index  int
	)

	if setScopes != nil {
		for scope = range setScopes {
			scopes[index] = scope
			index++
		}
	}

	return scopes
}

// setScopes defines the scopes used within this package and its binary representation. This function can
// only set scopes when there aren't any configured yet. A dynamic update is not supported for the sake
// of preventing concurrency issues without compromising performance.
func SetScopes(newScopes [64]string) error {
	var (
		index  int
		tmpMap map[string]*byte
	)

	if setScopes != nil {
		return errors.New(ErrorScopesAlreadySet)
	}

	tmpMap = make(map[string]*byte)

	for index = 0; index < 64; index++ {
		tmpMap[newScopes[index]] = &scopes[index]
	}

	setScopes = tmpMap
	return nil
}
