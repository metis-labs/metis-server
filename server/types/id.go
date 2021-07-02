/*
 * Copyright 2021-present NAVER Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
