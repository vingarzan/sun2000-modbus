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
	"net/http"
	"sync"
	"time"

	"github.com/goburrow/modbus"
)

// Global variables, because we are lazy
var (
	cfg config

	handlerModbus *modbus.TCPClientHandler
	clientModbus  modbus.Client

	lastSuccessTime   time.Time
	errorCount        uint
	totalErrorCount   uint
	totalSuccessCount uint

	parsedData sun2000DataStruct
)

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, parsedData.metricsString())
}

func main() {

	cfg.setDefaults()
	cfg.getFromEnv()

	var wg sync.WaitGroup
	var err error

	listenOn := fmt.Sprintf("%s:%s", cfg.httpIP, cfg.httpPort)

	http.HandleFunc("/metrics", handleMetrics)
	wg.Add(1)
	go func() {
		lError.Fatal(http.ListenAndServe(listenOn, nil))
	}()

	// Init the modbus client
	handlerModbus, clientModbus, err = initModbus(cfg.modbusIP, cfg.modbusPort, cfg.modbusTimeout)
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go readModbusLoop(uint(cfg.modbusSleep))

	wg.Wait()
	handlerModbus.Close()
}
