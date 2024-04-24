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
	"log"
	"os"
	"strconv"
)

type config struct {
	httpIP        string
	httpPort      string
	modbusIP      string
	modbusPort    uint16
	modbusTimeout uint
	modbusSleep   uint
}

func (c *config) setDefaults() {
	c.httpIP = "127.0.0.1"
	c.httpPort = "8080"
	c.modbusIP = ""
	c.modbusPort = 502
	c.modbusTimeout = 5
	c.modbusSleep = 5
}

func (c *config) getFromEnv() {
	x := os.Getenv("HTTP_IP")
	if len(x) > 0 {
		c.httpIP = x
	}
	x = os.Getenv("HTTP_PORT")
	if len(x) > 0 {
		c.httpPort = x
	}

	x = os.Getenv("MODBUS_IP")
	if len(x) > 0 {
		c.modbusIP = x
	} else {
		log.Fatal("MODBUS_IP is required! Please export it as an environment variable.")
	}
	x = os.Getenv("MODBUS_PORT")
	if len(x) > 0 {
		modbusPortUint, err := strconv.ParseUint(x, 10, 16)
		if err != nil {
			log.Fatal(err)
		}
		c.modbusPort = uint16(modbusPortUint)
	}

	x = os.Getenv("MODBUS_TIMEOUT")
	if len(x) > 0 {
		modbusTimeoutUint, err := strconv.ParseUint(x, 10, 16)
		if err != nil {
			log.Fatal(err)
		}
		c.modbusTimeout = uint(modbusTimeoutUint)
	}

	x = os.Getenv("MODBUS_SLEEP")
	if len(x) > 0 {
		modbusSleepUint, err := strconv.ParseUint(x, 10, 16)
		if err != nil {
			log.Fatal(err)
		}
		c.modbusSleep = uint(modbusSleepUint)
	}

}
