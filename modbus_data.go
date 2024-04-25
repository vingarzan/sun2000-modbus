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

import "time"

type modbusParsedData interface {
	isExpired() bool
	setLastRead(time.Time)
	setNextRead(time.Time)
	getNextRead() time.Time
	parse([]byte) error
	metricsString(*identificationData) string
}

type modbusInterval struct {
	name         string
	from, to     uint16
	target       modbusParsedData
	pullInterval time.Duration
}

// To optimize a bit the read process, we do bulk reads, in ranges of addresses.
// Then we parse the results and store them in the appropriate struct.
// Since some data is not updated very often, we can have different pull intervals for each address range.
// Seems like we could get at max 125 registers at once, or something around that.
var modbusAddrRanges = []modbusInterval{
	{
		name:         "Identification Data",
		from:         30000,
		to:           30087,
		target:       &parsedData.identification,
		pullInterval: 1 * time.Hour,
	},
	{
		name:         "Product Data",
		from:         30105,
		to:           30132,
		target:       &parsedData.product,
		pullInterval: 1 * time.Hour,
	},
	{
		name:         "Hardware Data Part 1",
		from:         30206,
		to:           30252,
		target:       &parsedData.hardware1,
		pullInterval: 1 * time.Hour,
	},
	{
		name:         "Hardware Data Part 2",
		from:         30300,
		to:           30327,
		target:       &parsedData.hardware2,
		pullInterval: 1 * time.Hour,
	},
	{
		name:         "Hardware Data Part 3",
		from:         30350,
		to:           30351,
		target:       &parsedData.hardware3,
		pullInterval: 1 * time.Hour,
	},
	// Does not work
	// {
	// 	name:         "Hardware Data Part 4",
	// 	from:         30364,
	// 	to:           30370,
	// 	target:       &parsedData.hardware4,
	// 	pullInterval: 5 * time.Minute,
	// },
	{
		name:         "Hardware Data Part 5",
		from:         31000,
		to:           31115,
		target:       &parsedData.hardware5,
		pullInterval: 1 * time.Hour,
	},
	// Does not work
	// {
	// 	name:         "Hardware Data Part 6",
	// 	from:         31130,
	// 	to:           31160,
	// 	target:       &parsedData.hardware6,
	// 	pullInterval: 1 * time.Hour,
	// },
	{
		name:         "Remote Signalling Data",
		from:         32000,
		to:           32003,
		target:       &parsedData.remoteSignalling,
		pullInterval: 1 * time.Hour,
	},
	{
		name:         "Alarm Data 1",
		from:         32008,
		to:           32011,
		target:       &parsedData.alarm1,
		pullInterval: 2 * time.Minute,
	},
	{
		name:         "PV Data",
		from:         32015,
		to:           32056,
		target:       &parsedData.pv,
		pullInterval: 10 * time.Second,
	},
	{
		name:         "Grid Data",
		from:         32064,
		to:           32097,
		target:       &parsedData.inverter,
		pullInterval: 10 * time.Second,
	},
	{
		name:         "Cumulative Data 1",
		from:         32106,
		to:           32120,
		target:       &parsedData.cumulative1,
		pullInterval: 2 * time.Minute,
	},
	{
		name:         "Cumulative Data 2",
		from:         32151,
		to:           32192,
		target:       &parsedData.cumulative2,
		pullInterval: 2 * time.Minute,
	},
	{
		name:         "Cumulative Data 3",
		from:         32190,
		to:           32192,
		target:       &parsedData.cumulative3,
		pullInterval: 2 * time.Minute,
	},
	{
		name:         "MPPT Data 1",
		from:         32212,
		to:           32232,
		target:       &parsedData.mppt1,
		pullInterval: 2 * time.Minute,
	},
	{
		name:         "Alarm Data 2",
		from:         32252,
		to:           32278,
		target:       &parsedData.alarm2,
		pullInterval: 2 * time.Minute,
	},
	// Does not work
	// {
	// 	name:         "String Access Data",
	// 	from:         32300,
	// 	to:           32318,
	// 	target:       &parsedData.stringAccess,
	// 	pullInterval: 1 * time.Hour,
	// },
	{
		name:         "MPPT Data 2",
		from:         32324,
		to:           32344,
		target:       &parsedData.mppt2,
		pullInterval: 2 * time.Minute,
	},
	{
		name:         "Internal Temperature Data",
		from:         35021,
		to:           35033,
		target:       &parsedData.internalTemperature,
		pullInterval: 2 * time.Minute,
	},
	{
		name:         "Meter Data",
		from:         37100,
		to:           37139,
		target:       &parsedData.meter,
		pullInterval: 10 * time.Second,
	},
	{
		name:         "ESU1 Data",
		from:         37000,
		to:           37070,
		target:       &parsedData.esu1,
		pullInterval: 30 * time.Second,
	},
	{
		name:         "ESU2 Data",
		from:         37700,
		to:           37757,
		target:       &parsedData.esu2,
		pullInterval: 30 * time.Second,
	},
	{
		name:         "ESU1-Pack1 Data",
		from:         38200,
		to:           38242,
		target:       &parsedData.esu1.pack[0],
		pullInterval: 30 * time.Second,
	},
	{
		name:         "ESU1-Pack2 Data",
		from:         38242,
		to:           38284,
		target:       &parsedData.esu1.pack[1],
		pullInterval: 30 * time.Second,
	},
	// // I don't have this one
	// {
	// 	name:         "ESU1-Pack3 Data",
	// 	from:         38284,
	// 	to:           38326,
	// 	target:       &parsedData.esu1.pack[2],
	// 	pullInterval: 30 * time.Second,
	// },
	// // I don't have this one
	// {
	// 	name:         "ESU2-Pack1 Data",
	// 	from:         38326,
	// 	to:           38368,
	// 	target:       &parsedData.esu2.pack[0],
	// 	pullInterval: 30 * time.Second,
	// },
	// // I don't have this one
	// {
	// 	name:         "ESU2-Pack2 Data",
	// 	from:         38368,
	// 	to:           38410,
	// 	target:       &parsedData.esu2.pack[1],
	// 	pullInterval: 30 * time.Second,
	// },
	// // I don't have this one
	// {
	// 	name:         "ESU2-Pack3 Data",
	// 	from:         38410,
	// 	to:           38452,
	// 	target:       &parsedData.esu2.pack[2],
	// 	pullInterval: 30 * time.Second,
	// },

	{
		name:         "ESU Temperatures",
		from:         38452,
		to:           38452 + 2*3*2,
		target:       &parsedData.esuTemperatures,
		pullInterval: 30 * time.Second,
	},
}

func readModbusLoop(pollInterval uint) {

	// dummy read, since the first call always seems to fail - inverted bug?
	readModbusFromTo("dummy read", 30000, 30015)

	for {

		for _, addrRange := range modbusAddrRanges {
			if !addrRange.target.isExpired() {
				continue
			}

			results, err := readModbusFromTo(addrRange.name, addrRange.from, addrRange.to)
			ok := handleReadModbusResults(results, err)
			if ok {
				// Interpret the results
				err := addrRange.target.parse(results)
				if err != nil {
					lError.Printf("Error parsing %s: %v", addrRange.name, err)
				}
				addrRange.target.setLastRead(time.Now())
				addrRange.target.setNextRead(time.Now().Add(addrRange.pullInterval))
				lInfo.Printf("Will read again the %s after %s", addrRange.name, addrRange.target.getNextRead().Format(time.RFC3339))
			}
		}

		lDebug.Printf("... sleeping %d seconds...", pollInterval)
		time.Sleep(time.Duration(time.Duration(pollInterval) * time.Second))
	}
}
