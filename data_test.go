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
	"reflect"
	"testing"
)

func TestGetAlarms(t *testing.T) {
	// Define the input data
	id := [alarmCount]uint16{0b0110000000000000, 0b000000000000001, 0b1000000000000000, 0b0010000000000000}

	// Call the function under test
	result := getAlarms(id)

	// Define the expected output
	expected := []sun2000Alarm{
		sun2000Alarms[1],
		sun2000Alarms[2],
		sun2000Alarms[31],
		sun2000Alarms[32],
		// Add more expected alarms here if needed
	}

	// Compare the result with the expected output
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("getAlarms() returned unexpected result, got: %v, want: %v", result, expected)
	}
}
