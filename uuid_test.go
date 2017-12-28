package uuid_test

import (
	"github.com/4xoc/uuid"
	"testing"
)

func TestMain(t *testing.T) {
	var (
		myScopes    [64]string
		mySetScopes [64]string
		myUUID      *uuid.UUID
		myUUID2     *uuid.UUID
		err         error
	)

	//trying uninitialized package (no scopes set)
	_, err = uuid.New("foo")
	if err == nil {
		t.Error("There are no scopes defined thus there should be no new uuid")
	}

	//setting scopes
	myScopes = [64]string{"one", "two", "three", "four", "five", "six", "seven", "eight"}
	err = uuid.SetScopes(myScopes)

	if err != nil {
		t.Error("scopes should have been set")
	}

	//getting scopes
	mySetScopes = uuid.Scopes()
	if len(mySetScopes) != len(myScopes) {
		t.Error("Expected ", len(myScopes), "# of scopes but got ", mySetScopes)
	}

	//getting scope on nil ptr
	if myUUID.Scope() != "" {
		t.Error("Scope wasn't empty string as expected")
	}

	//getting hex on nil ptr
	if myUUID.Hex() != "" {
		t.Error("Hex wasn't empty string as expected")
	}

	//getting bin on nil ptr
	if myUUID.Bin() != [16]byte{} {
		t.Error("Bin data should be an empty byte arra but wasn't.")
	}

	//creating new uuid with known scope
	myUUID, err = uuid.New("five")

	if err != nil {
		t.Error("Expected UUID to be generated but failed with error ", err.Error())
	}

	if myUUID.Scope() != "five" {
		t.Error("UUID does not match the scope defined on creation time.")
	}

	//read good UUID
	myUUID2, err = uuid.Read(myUUID.Hex())

	if err != nil {
		t.Error("Expected UUID to be generated but failed with error ", err.Error())
	}

	if myUUID.Hex() != myUUID2.Hex() ||
		myUUID.Bin() != myUUID2.Bin() ||
		myUUID.Scope() != myUUID2.Scope() {
		t.Error("UUIDs should be identical but aren't")
	}

	if !myUUID2.ScopeMatches(myScopes[:]) {
		t.Error("Scope should match but did not.")
	}

	if myUUID2.ScopeMatches([]string{"ten"}) {
		t.Error("Scope should match but did not.")
	}

	//now to the bad things

	//reading a bad UUID
	_, err = uuid.Read(myUUID.Hex()[:1])

	if err == nil {
		t.Error("UUID shouldn't have been generated")
	}

	//now setting new set of scopes
	myScopes = [64]string{"one"}
	err = uuid.SetScopes(myScopes)

	if err == nil {
		t.Error("setting new scopes should not have been possible")
	}

	//reading a previously valid UUID which now has an unknown scope
	_, err = uuid.Read("ff8cb1d0-84f3-9d8d-76cc-682d1ca34dae")

	if err == nil {
		t.Error("UUID shouldn't have been generated")
	}
}
