// Copyright: 2024 Dragos Vingarzan vingarzan -at- gmail -dot- com
// License: AGPL-3.0
//
// This file is part of sun2000-modbus.
//
// sun2000-modbus is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General
// Public License Version 3 (AGPL-3.0) as published by the Free Software Foundation.
//
// sun2000-modbus is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied
// warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more
// details.
//
// You should have received a copy of the AGPL-3.0 along with sun2000-modbus. If not,
// see <https://www.gnu.org/licenses/>.
package main

import (
	"encoding/binary"
	"fmt"
)

func getSTR(data []byte, idx, size uint) (out string, idxOut uint, err error) {
	if len(data) < int(idx+2*size) {
		lError.Printf("data length %d < %d", len(data), idx+2*size)
		return "", idx, fmt.Errorf("data length %d < %d", len(data), idx+2*size)
	}
	in := data[idx : idx+2*size]
	first := 0
	for i := 0; i < len(in); i++ {
		if in[i] != 0 {
			first = i
			break
		}
	}
	last := 0
	for i := len(in) - 1; i >= 0; i-- {
		if in[i] != 0 {
			last = i + 1
			break
		}
	}
	x := in[first:last]
	for i := 0; i < len(x); i++ {
		if x[i] == 0 {
			x[i] = '.' // replace nulls with dots
		}
	}
	return string(x), idx + 2*size, nil
}

func getU16(data []byte, idx uint) (out uint16, idxOut uint, err error) {
	var size uint = 1
	if len(data) < int(idx+2*size) {
		lError.Printf("data length %d < %d", len(data), idx+2*size)
		return 0, idx, fmt.Errorf("data length %d < %d", len(data), idx+2*size)
	}
	out = binary.BigEndian.Uint16(data[idx : idx+2*size])
	return out, idx + 2*size, nil
}

func getU32(data []byte, idx uint) (out uint32, idxOut uint, err error) {
	var size uint = 2
	if len(data) < int(idx+2*size) {
		lError.Printf("data length %d < %d", len(data), idx+2*size)
		return 0, idx, fmt.Errorf("data length %d < %d", len(data), idx+2*size)
	}
	out = binary.BigEndian.Uint32(data[idx : idx+2*size])
	return out, idx + 2*size, nil
}

func getI16(data []byte, idx uint) (out int16, idxOut uint, err error) {
	var size uint = 1
	if len(data) < int(idx+2*size) {
		lError.Printf("data length %d < %d", len(data), idx+2*size)
		return 0, idx, fmt.Errorf("data length %d < %d", len(data), idx+2*size)
	}
	out = int16(binary.BigEndian.Uint16(data[idx : idx+2*size]))
	return out, idx + 2*size, nil
}

func getI32(data []byte, idx uint) (out int32, idxOut uint, err error) {
	var size uint = 2
	if len(data) < int(idx+2*size) {
		lError.Printf("data length %d < %d", len(data), idx+2*size)
		return 0, idx, fmt.Errorf("data length %d < %d", len(data), idx+2*size)
	}
	out = int32(binary.BigEndian.Uint32(data[idx : idx+2*size]))
	return out, idx + 2*size, nil
}

func skipRecords(data []byte, idx, size uint) (idxOut uint, err error) {
	if len(data) < int(idx+2*size) {
		lError.Printf("data length %d < %d", len(data), idx+2*size)
		return idx, fmt.Errorf("data length %d < %d", len(data), idx+2*size)
	}
	return idx + 2*size, nil
}
