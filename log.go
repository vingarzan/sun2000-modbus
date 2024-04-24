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
)

var (
	lDebug   *log.Logger
	lInfo    *log.Logger
	lWarning *log.Logger
	lError   *log.Logger

	lModBus *log.Logger
)

const (
	ANSI_RESET  = "\033[0m"
	ANSI_GRAY   = "\033[38;5;8m"
	ANSI_GREEN  = "\033[38;5;2m"
	ANSI_YELLOW = "\033[38;5;3m"
	ANSI_RED    = "\033[38;5;1m"
	ANSI_BLUE   = "\033[38;5;4m"
)

func init() {

	// we could switch to log/slog for structured logging... that has built-in levels and attributes, context,
	// formatting, etc... but doesn't have the simpler Printf.

	logFormat := log.Ldate | log.Lmicroseconds | log.Ltime | log.Lshortfile
	// logFormat |= log.Lmsgprefix

	prefix := "sun2000-modbus"

	lDebug = log.New(os.Stderr, ANSI_GRAY+prefix+"-DBG  "+ANSI_RESET+"| ", logFormat)
	lInfo = log.New(os.Stderr, ANSI_GREEN+prefix+"-INFO "+ANSI_RESET+"| ", logFormat)
	lWarning = log.New(os.Stderr, ANSI_GREEN+prefix+"-WARN "+ANSI_RESET+"| ", logFormat)
	lError = log.New(os.Stderr, ANSI_GREEN+prefix+"-ERR  "+ANSI_RESET+"| ", logFormat)

	lModBus = log.New(os.Stderr, ANSI_BLUE+prefix+"-MBUS "+ANSI_RESET+"| ", logFormat)

	// setting also the default logger, just in case something uses just log.stuff
	log.SetFlags(logFormat)
	log.SetPrefix(prefix + "      | ")
}
