package types

import (
	"encoding/hex"
)

// ID represents ID of entity.
type ID string

// String returns a string representation of this ID.
func (id ID) String() string {
	return string(id)
}

// Bytes returns bytes of decoded hexadecimal string representation of this ID.
func (id ID) Bytes() []byte {
	decoded, err := hex.DecodeString(id.String())
	if err != nil {
		return nil
	}
	return decoded
}

// IDFromBytes returns ID represented by the encoded hexadecimal string from bytes.
func IDFromBytes(bytes []byte) ID {
	return ID(hex.EncodeToString(bytes))
}
