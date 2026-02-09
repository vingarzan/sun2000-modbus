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
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/goburrow/modbus"
)

func initModbus(ip string, port uint16, timeout uint, slaveID byte) (handler *modbus.TCPClientHandler, client modbus.Client, err error) {
	// Modbus TCP
	modbus_target := fmt.Sprintf("%s:%d", ip, port)
	handler = modbus.NewTCPClientHandler(modbus_target)
	handler.Timeout = time.Duration(timeout) * time.Second
	handler.SlaveId = slaveID
	handler.Logger = lModBus
	// Connect manually so that multiple requests are handled in one connection session
	err = handler.Connect()
	if err != nil {
		return nil, nil, fmt.Errorf("modbus failed to connect to %s:%d: %v", ip, port, err)
	}

	client = modbus.NewClient(handler)

	return handler, client, nil
}

func readModbusFromTo(what string, from uint16, to uint16) (results []byte, err error) {
	lDebug.Printf("   >>   Reading %s from modbus %d..%d\n", what, from, to)

	size := (to - from)
	return clientModbus.ReadHoldingRegisters(from, size)
}

func handleReadModbusResults(results []byte, err error) (ok bool) {
	if err != nil {
		lWarning.Printf("Error reading modbus: %v\n", err)
		errorCount++
		totalErrorCount++
		if errorCount > 10 {
			handlerModbus.Close()
			handlerModbus = nil

			lWarning.Printf("Too many errors, sleeping for 3 minute\n")
			time.Sleep(3 * time.Minute)
			handlerModbus, clientModbus, err = initModbus(cfg.modbusIP, cfg.modbusPort, uint(cfg.modbusTimeout), cfg.modbusSlaveID)
			if err != nil {
				log.Fatal(err)
			}
			return false
		}
		if strings.Contains(err.Error(), "modbus: response transaction id") {
			lWarning.Printf("modbus: sun2000 has a known bug where it fucks up transaction ids... we must close the connection, then reopen.")
			handlerModbus.Close()
			handlerModbus = nil
			lInfo.Printf("Sleeping 30 seconds...")
			time.Sleep(30 * time.Second)
			handlerModbus, clientModbus, err = initModbus(cfg.modbusIP, cfg.modbusPort, uint(cfg.modbusTimeout), cfg.modbusSlaveID)
			if err != nil {
				log.Fatal(err)
			}
		}
		// TODO - how to detect if the connection was dropped and to re-open it?
		return false
	} else {
		errorCount = 0
		totalSuccessCount++
		lastSuccessTime = time.Now()
	}
	// lDebug.Printf("Results: %q\n", results)
	return true
}
