# UUID - Custom IDs

This package provides functionality to generate random UUIDs with unique identification of its type within. Any UUID is randomly generated (version 4 UUID) but has certain bits set to identify exactly the 'type' is is refering to. This package does **NOT** generate UUIDs following RFC 4122. However, it nevertheless allows for exactly the same number of possible combinations (ignoring the type bits) which is 2^122 UUIDs **for each type**.
UUID also has no 3rd party dependencies meaning that out-of-the-box just golang is needed.

## Open Issues
- no safe sync in place for concurrent read/write of scopes map

## Types and UUID Structure
Generally referred to as 'types' is the combination of the first 6 (most-significant) bits that basically set a group of UUIDs each belongs to. If you for example want to generate UUIDs for users, then the type could be `0x00` or `000000` in binary. Each user UUID will therefore always start with `00` in the hex-string identifying it easily as a user UUID.

## But Why?
Using specific identifying bits allows one to continue using UUIDs (which are generally used across many systems and are widely used) while keeping track of what ID is supposed to identify. Imagine a database with users, posts and comments: One would rely on a UUID being passed in any way to be of the correct type. Detecting that an ID is incorrect would at least require one DB query while having set types, an ID will tell itself that it represents a very specific type. 64 different types can be used with 2^122 UUIDs for each type. This is the same space RFC compliant UUIDs have too.

## Usage
There are a couple of very important points to get started. Once you're through that the rest is a piece of cake.

1. Import the package (obviously)
```
import (
    "github.com/4xoc/uuid"
)
```

2. You must initialize a string array with size 64 that defines the scopes of UUIDs within your project. This is necessary to ensure that the scopes and its binary representation never changes when re-running the program. Slice and map are not very safe. Only scopes defines in this array can be used when creating or reading UUIDs.
```
// declaring scopes; it can also be a constant
const MY_CONST_SCOPE string = "four"
myScopes := [64]string{
    "one", "two", "three", MY_CONST_SCOPE,
}

uuid.SetScopes(myScopes)
```

3. Now we can create a new UUID
```
myUUID, err := uuid.New("one")

if err != nil {
    //basically only happens when the scope it unknown
    //use the packages constants which hold the error strings for detecting the reason of the failure.
    fmt.Printf("Failed to generate UUID with error: %s\n", err.Error())
}

//now we can access the information of the uuid
fmt.Printf("Generated UUID %s with scope %s and binary %08b\n", myUUID.Hex(), myUUID.Scope(), myUUID.Bin())

```

4. And then we try reading one
```
myCopy, _ = uuid.Read(myUUID.Hex())

fmt.Printf("Copied UUID %s with scope %s and binary %08b\n", myCopy.Hex(), myCopy.Scope(), myCopy.Bin())

```

5. Checking if a uuid matches any of the given scopes. This can be helpful to validate UUIDs before processing them further.
```
//one or more allowed scopes can be set, if any matches it'll be a true response
myScopes := []string{"one","two"}

if myCopy.ScopeMatches(myScope) {
    fmt.Println("Yay, this uuid is of the allowed scopes.")
}
```

## Database
This package implements the `database/sql/driver` interfaces to give Golang's built-in sql package access to read and write data from and into the struct transparently for the developer. Simply use the UUID type in your structs and the interfaces do the magic themselves.

## FAQ
**Dude, why do I always need to call a function to just get a value?**  
All fields of the struct are not directly accessable to prevent problems with manual changes bin/scope/hex data that would either cause a panic or at least become unpredictable in its workings. Therefore only interfaces allow the access to actual values so that a change of any data always also updates the other (if necessary).

## Contribution
Anyone feel happy to get involved and respond to issues or simply create a PR.
