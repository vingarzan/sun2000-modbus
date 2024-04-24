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
	"strings"
	"sync"
	"time"
)

type sun2000DataStruct struct {

	// Identification Data
	identification   identificationData
	product          productData
	hardware1        hardwareData1
	hardware2        hardwareData2
	hardware3        hardwareData3
	hardware4        hardwareData4
	hardware5        hardwareData5
	hardware6        hardwareData6
	remoteSignalling remoteSignallingData
	alarm            alarmData
	pv               pvData
	grid             gridData
}

func (x *sun2000DataStruct) metricsString() string {
	sb := strings.Builder{}
	sb.WriteString("# Huawei Sun2000 inverter scraped data from ModBus TCP\n#\n")
	sb.WriteString(fmt.Sprintf("#  - last read success at %s\n", lastSuccessTime.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("#  - consecutive read errors %d\n", errorCount))
	sb.WriteString(fmt.Sprintf("#  - total read errors %d\n", totalErrorCount))
	sb.WriteString(fmt.Sprintf("#  - total read successes %d\n", totalSuccessCount))
	sb.WriteString("\n")
	sb.WriteString(x.identification.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.product.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.hardware1.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.hardware2.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.hardware3.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.hardware4.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.hardware5.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.hardware6.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.remoteSignalling.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.alarm.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.pv.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.grid.metricsString(&x.identification))

	return sb.String()
}

type genericData struct {
	// RW mutex to protect the data
	sync.RWMutex
	lastRead time.Time
	nextRead time.Time
}

func (x *genericData) isExpired() bool {
	x.RLock()
	defer x.RUnlock()
	return x.nextRead.Before(time.Now())
}

func (x *genericData) setLastRead(t time.Time) {
	x.Lock()
	defer x.Unlock()
	x.lastRead = t
}

func (x *genericData) setNextRead(t time.Time) {
	x.Lock()
	defer x.Unlock()
	x.nextRead = t
}

func (x *genericData) getNextRead() time.Time {
	x.RLock()
	defer x.RUnlock()
	return x.nextRead
}

func (x *genericData) parse(data []byte) (err error) {
	return fmt.Errorf("parse not implemented")
}

func (x *genericData) metricsString(id *identificationData) string {
	return fmt.Sprintf("# No metrics for %T\n", x)
}

type identificationData struct {
	genericData

	// 30000 STR 15
	model string
	// 30015 STR 10
	sn string
	// 30025 STR 10
	pn string

	// 30035 STR 15
	firmwareVersion string
	// 30050 STR 15
	softwareVersion string

	// 30068 U32 2
	protocolVersion uint32 // modbus
	// 30070 U16 1
	modelID uint16
	// 30071 U16 1
	numberOfStrings uint16
	// 30072 U16 1
	numberOfMPPTs uint16
	// 30073 U32 2 gain 1000 kW
	ratedPower float32
	// 30075 U32 2 gain 1000 kW
	maxActivePowerPmax float32
	// 30077 U32 2 gain 1000 kVA Pmax
	maxApparentPowerSmax float32
	// 30079 I32 2 gain 1000 kVar Qmax
	realtimeMaxReactivePowerQmaxFeedToGrid float32
	// 30081 I32 2 gain 1000 kVar -Qmax
	realtimeMaxReactivePowerQmaxAbsorbedFromGrid float32
	// 30083 U32 2 gain 1000 kW Pmax_real - 0<Pmax≤Smax≤Pmax_real≤Smax_real or 0<Pmax≤Pmax_real≤Smax≤Smax_real
	maxActiveCapabilityPmaxReal float32
	// 30085 U32 2 gain 1000 kVA Smax_real - 0<Pmax≤Smax≤Pmax_real≤Smax_real or 0<Pmax≤Pmax_real≤Smax≤Smax_real
	maxApparentCapabilitySmaxReal float32
}

func (x *identificationData) parse(data []byte) (err error) {
	if len(data) < 70 {
		return fmt.Errorf("data length %d < 70", len(data))
	}
	x.Lock()
	defer x.Unlock()

	var idx uint
	x.model, idx, _ = getSTR(data, idx, 15)
	x.sn, idx, _ = getSTR(data, idx, 10)
	x.pn, idx, _ = getSTR(data, idx, 10)
	x.firmwareVersion, idx, _ = getSTR(data, idx, 15)
	x.softwareVersion, idx, _ = getSTR(data, idx, 15)
	idx, _ = skipRecords(data, idx, 3)
	x.protocolVersion, idx, _ = getU32(data, idx)
	x.modelID, idx, _ = getU16(data, idx)
	x.numberOfStrings, idx, _ = getU16(data, idx)
	x.numberOfMPPTs, idx, _ = getU16(data, idx)
	u32, idx, _ := getU32(data, idx)
	x.ratedPower = float32(u32) / 1000
	u32, idx, _ = getU32(data, idx)
	x.maxActivePowerPmax = float32(u32) / 1000
	u32, idx, _ = getU32(data, idx)
	x.maxApparentPowerSmax = float32(u32) / 1000
	i32, idx, _ := getI32(data, idx)
	x.realtimeMaxReactivePowerQmaxFeedToGrid = float32(i32) / 1000
	i32, idx, _ = getI32(data, idx)
	x.realtimeMaxReactivePowerQmaxAbsorbedFromGrid = float32(i32) / 1000
	u32, idx, _ = getU32(data, idx)
	x.maxActiveCapabilityPmaxReal = float32(u32) / 1000
	u32, _, _ = getU32(data, idx)
	x.maxApparentCapabilitySmaxReal = float32(u32) / 1000

	return nil
}

func (x *identificationData) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Identification Data\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Model            = %q\n", x.model))
	sb.WriteString(fmt.Sprintf("# SN               = %q\n", x.sn))
	sb.WriteString(fmt.Sprintf("# PN               = %q\n", x.pn))
	sb.WriteString(fmt.Sprintf("# Firmware Version = %q\n", x.firmwareVersion))
	sb.WriteString(fmt.Sprintf("# Software Version = %q\n", x.softwareVersion))
	sb.WriteString(fmt.Sprintf("# Protocol Version = %#08x (D%d.%d)\n", x.protocolVersion, x.protocolVersion>>16&0xffff, x.protocolVersion&0xffff))
	sb.WriteString(fmt.Sprintf("# Model ID         = %d\n", x.modelID))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Number of Strings = %d\n", x.numberOfStrings))
	sb.WriteString(fmt.Sprintf("# Number of MPPTs   = %d\n", x.numberOfMPPTs))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Rated Power                                         = %2.3f kW\n", x.ratedPower))
	sb.WriteString(fmt.Sprintf("# Max Active Power Pmax                               = %2.3f kW\n", x.maxActivePowerPmax))
	sb.WriteString(fmt.Sprintf("# Max Apparent Power Smax                             = %2.3f kVA\n", x.maxApparentPowerSmax))
	sb.WriteString(fmt.Sprintf("# Realtime Max Reactive Power Qmax Feed to Grid       = %2.3f kVar\n", x.realtimeMaxReactivePowerQmaxFeedToGrid))
	sb.WriteString(fmt.Sprintf("# Realtime Max Reactive Power Qmax Absorbed from Grid = %2.3f kVar\n", x.realtimeMaxReactivePowerQmaxAbsorbedFromGrid))
	sb.WriteString(fmt.Sprintf("# Max Active Capability Pmax Real                     = %2.3f kW    0<Pmax≤Smax≤Pmax_real≤Smax_real or 0<Pmax≤Pmax_real≤Smax≤Smax_real\n", x.maxActiveCapabilityPmaxReal))
	sb.WriteString(fmt.Sprintf("# Max Apparent Capability Smax Real                   = %2.3f kVA    0<Pmax≤Smax≤Pmax_real≤Smax_real or 0<Pmax≤Pmax_real≤Smax≤Smax_real\n", x.maxApparentCapabilitySmaxReal))

	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() {
		sb.WriteString("# No identification data read yet\n")
	} else {
		sb.WriteString(fmt.Sprintf("sun2000_inverter_number_of_MPPTs{model=%q,sn=%q} %d\n", x.model, x.sn, x.numberOfMPPTs))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_number_of_strings{model=%q,sn=%q} %d\n", x.model, x.sn, x.numberOfStrings))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_rated_power{model=%q,sn=%q,unit=\"kW\",description=\"Rated Power\"} %.3f\n", x.model, x.sn, x.ratedPower))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_Pmax{model=%q,sn=%q,unit=\"kW\",description=\"Maximum Active Power Pmax\"} %.3f\n", x.model, x.sn, x.maxActivePowerPmax))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_Smax{model=%q,sn=%q,unit=\"kVA\",description=\"Maximum Apparent Power Smax\"} %.3f\n", x.model, x.sn, x.maxApparentPowerSmax))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_Qmax_feed_to_grid{model=%q,sn=%q,unit=\"kVar\",description=\"Realtime Max Reactive Power Qmax Feed to Grid\"} %.3f\n", x.model, x.sn, x.realtimeMaxReactivePowerQmaxFeedToGrid))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_Qmax_absorbed_from_grid{model=%q,sn=%q,unit=\"kVar\",description=\"Realtime Max Reactive Power Qmax Absorbed from Grid\"} %.3f\n", x.model, x.sn, x.realtimeMaxReactivePowerQmaxAbsorbedFromGrid))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_Pmax_real{model=%q,sn=%q,unit=\"kW\",description=\"Maximum Active Capability Pmax Real\"} %.3f\n", x.model, x.sn, x.maxActiveCapabilityPmaxReal))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_Smax_real{model=%q,sn=%q,unit=\"kVA\",description=\"Maximum Apparent Capability Smax Real\"} %.3f\n", x.model, x.sn, x.maxApparentCapabilitySmaxReal))

	}
	sb.WriteString("\n")
	return sb.String()
}

type productData struct {
	genericData

	// 30105 STR 2
	productSalesArea string
	// 30107 U16 1
	productSoftwareNumber uint16
	// 30108 U16 1
	productSoftwareVersionNumber uint16
	// 30109 U16 1
	gridStandardCodeProtocolVersion uint16
	// 30110 U16 1
	uniqueIDOfTheSoftware uint16
	// 30111 U16 1
	numberOfPackagesToBeUpgraded uint16
	// 30112-30130 U32 2x10
	subpackageInformation [10]uint32
}

func (x *productData) parse(data []byte) (err error) {
	if len(data) < 2*17 {
		return fmt.Errorf("data length %d < 2*17", len(data))
	}

	x.Lock()
	defer x.Unlock()

	var idx uint
	x.productSalesArea, idx, _ = getSTR(data, idx, 2)
	x.productSoftwareNumber, idx, _ = getU16(data, idx)
	x.productSoftwareVersionNumber, idx, _ = getU16(data, idx)
	x.gridStandardCodeProtocolVersion, idx, _ = getU16(data, idx)
	x.uniqueIDOfTheSoftware, idx, _ = getU16(data, idx)
	x.numberOfPackagesToBeUpgraded, idx, _ = getU16(data, idx)
	for i := 0; i < 10; i++ {
		x.subpackageInformation[i], idx, _ = getU32(data, idx)
	}

	return nil
}

func (x *productData) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Product Data\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Product Sales Area = %q\n", x.productSalesArea))
	sb.WriteString(fmt.Sprintf("# Product Software Number             = %d\t%#04x\n", x.productSoftwareNumber, x.productSoftwareNumber))
	sb.WriteString(fmt.Sprintf("# Product Software Version Number     = %d\t%#04x\n", x.productSoftwareVersionNumber, x.productSoftwareVersionNumber))
	sb.WriteString(fmt.Sprintf("# Grid Standard Code Protocol Version = %d\t%#04x\n", x.gridStandardCodeProtocolVersion, x.gridStandardCodeProtocolVersion))
	sb.WriteString(fmt.Sprintf("# Unique ID Of The Software           = %d\t%#04x\n", x.uniqueIDOfTheSoftware, x.uniqueIDOfTheSoftware))
	sb.WriteString(fmt.Sprintf("# Number Of Packages To Be Upgraded   = %d\n", x.numberOfPackagesToBeUpgraded))
	sb.WriteString("\n")
	for i, v := range x.subpackageInformation {
		sb.WriteString(fmt.Sprintf("# Subpackage %2d Information  = %#08x\tfileTypeID=%3d\tdeviceTypeId=%d\n", i+1, v, v>>16&0xff, v&0xffff))
	}
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No product or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()
		sb.WriteString(fmt.Sprintf("sun2000_inverter_unique_id_of_the_software{model=%q,sn=%q} %d\n", id.model, id.sn, x.uniqueIDOfTheSoftware))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_number_of_packages_to_be_upgraded{model=%q,sn=%q} %d\n", id.model, id.sn, x.numberOfPackagesToBeUpgraded))
	}
	sb.WriteString("\n")

	return sb.String()
}

type hardwareData1 struct {
	genericData

	// 30206 Bitfield16 1
	hardwareFunctionalUnitConfigurationIdentifier uint16
	// 30207 Bitfield32 2
	subdeviceSupportFlag uint32
	// 30209 Bitfield32 2
	subdeviceInPositionFlag uint32
	// 30211-30217 Bitfield32 2x4
	featureMask [4]uint32
	// 30219-30250 Bitfield16 1x32
	gridStandardCodeMask [32]uint16
}

func (x *hardwareData1) parse(data []byte) (err error) {
	size := 1 + 2 + 2 + 2*4 + 1*32
	if len(data) < size*2 {
		return fmt.Errorf("data length %d < %d*2", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var idx uint
	x.hardwareFunctionalUnitConfigurationIdentifier, idx, _ = getU16(data, idx)
	x.subdeviceSupportFlag, idx, _ = getU32(data, idx)
	x.subdeviceInPositionFlag, idx, _ = getU32(data, idx)
	for i := 0; i < 4; i++ {
		x.featureMask[i], idx, _ = getU32(data, idx)
	}
	for i := 0; i < 32; i++ {
		x.gridStandardCodeMask[i], idx, _ = getU16(data, idx)
	}

	return nil
}

func (x *hardwareData1) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Hardware Data Part 1\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	var text string
	switch x.hardwareFunctionalUnitConfigurationIdentifier {
	case 0:
		text = "no functional unit hardware configuration"
	case 1:
		text = "The hardware configuration of the functional unit is available"
	case 2:
		text = ""
	}
	sb.WriteString(fmt.Sprintf("# Hardware Functional Unit Configuration Identifier = %#04x\t%#016b\t%s\n", x.hardwareFunctionalUnitConfigurationIdentifier, x.hardwareFunctionalUnitConfigurationIdentifier, text))
	sb.WriteString(fmt.Sprintf("# Subdevice Support Flag     = %#08x\t%#032b\n", x.subdeviceSupportFlag, x.subdeviceSupportFlag))
	sb.WriteString(fmt.Sprintf("# Subdevice In Position Flag = %#08x\t%#032b\n", x.subdeviceInPositionFlag, x.subdeviceInPositionFlag))
	sb.WriteString("\n")
	for i, v := range x.featureMask {
		sb.WriteString(fmt.Sprintf("# Feature Mask %d             = %#08x\t%#032b\n", i+1, v, v))
	}
	sb.WriteString("\n")
	for i, v := range x.gridStandardCodeMask {
		sb.WriteString(fmt.Sprintf("# Grid Standard Code Mask %2d = %#04x\t%#016b\n", i+1, v, v))
	}
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No hardware or identification data read yet\n")
	}
	// else {
	// 	// no useful metrics here
	// 	id.RLock()
	// 	defer id.RUnlock()
	// }
	sb.WriteString("\n")

	return sb.String()
}

type hardwareData2 struct {
	genericData

	// 30300-30307 Bitfield16 1x8
	monitoringParameterMask [8]uint16
	// 30308-30326 Bitfield16 1x19
	powerParameterMask [19]uint16
}

func (x *hardwareData2) parse(data []byte) (err error) {
	size := 8 + 19
	if len(data) < size*2 {
		return fmt.Errorf("data length %d < %d*2", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var idx uint
	for i := 0; i < 8; i++ {
		x.monitoringParameterMask[i], idx, _ = getU16(data, idx)
	}
	for i := 0; i < 19; i++ {
		x.powerParameterMask[i], idx, _ = getU16(data, idx)
	}

	return nil
}

func (x *hardwareData2) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Hardware Data Part 2\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	for i, v := range x.monitoringParameterMask {
		sb.WriteString(fmt.Sprintf("# Monitoring Parameter Mask %2d = %#04x\t%#016b\n", i+1, v, v))
	}
	sb.WriteString("\n")
	for i, v := range x.powerParameterMask {
		sb.WriteString(fmt.Sprintf("# Power Parameter Mask %2d      = %#04x\t%#016b\n", i+1, v, v))
	}
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No hardware or identification data read yet\n")
	}
	// else {
	// 	// no useful metrics here
	// 	id.RLock()
	// 	defer id.RUnlock()
	// }
	sb.WriteString("\n")

	return sb.String()
}

type hardwareData3 struct {
	genericData

	// 30350 U16 1
	builtinPIDParameterMask uint16
}

func (x *hardwareData3) parse(data []byte) (err error) {
	if len(data) < 2 {
		return fmt.Errorf("data length %d < 2", len(data))
	}

	x.Lock()
	defer x.Unlock()

	var idx uint
	x.builtinPIDParameterMask, _, _ = getU16(data, idx)

	return nil
}

func (x *hardwareData3) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Hardware Data Part 3\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Builtin PID Parameter Mask = %#04x\t%#016b\n", x.builtinPIDParameterMask, x.builtinPIDParameterMask))
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No hardware or identification data read yet\n")
	}
	// else {
	// 	// no useful metrics here
	// 	id.RLock()
	// 	defer id.RUnlock()
	// }
	sb.WriteString("\n")

	return sb.String()
}

type hardwareData4 struct {
	genericData

	// 30364 U32 2
	realtimeMaxActiveCapability uint32
	// 30366 I32 2
	realtimeMaxCapacitiveReactiveCapacityPlus int32
	// 30368 I32 2
	realtimeMaxInductiveReactiveCapacityMinus int32
}

func (x *hardwareData4) parse(data []byte) (err error) {
	if len(data) < 6 {
		return fmt.Errorf("data length %d < 6", len(data))
	}

	x.Lock()
	defer x.Unlock()

	var idx uint
	x.realtimeMaxActiveCapability, idx, _ = getU32(data, idx)
	x.realtimeMaxCapacitiveReactiveCapacityPlus, idx, _ = getI32(data, idx)
	x.realtimeMaxInductiveReactiveCapacityMinus, _, _ = getI32(data, idx)

	return nil
}

func (x *hardwareData4) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Hardware Data Part 4\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Realtime Max Active Capability                 = %d\n", x.realtimeMaxActiveCapability))
	sb.WriteString(fmt.Sprintf("# Realtime Max Capacitive Reactive Capacity (+)  = %d\n", x.realtimeMaxCapacitiveReactiveCapacityPlus))
	sb.WriteString(fmt.Sprintf("# Realtime Max Inductive Reactive Capacity (-)   = %d\n", x.realtimeMaxInductiveReactiveCapacityMinus))
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No hardware or identification data read yet\n")
	}
	// else {
	// 	// no useful metrics here
	// 	id.RLock()
	// 	defer id.RUnlock()
	// }
	sb.WriteString("\n")

	return sb.String()
}

type hardwareData5 struct {
	genericData

	// 31000 STR 15
	hardwareVersion string
	// 31015 STR 10
	monitoringBoardSN string
	// 31025 STR 15
	monitoringSoftwareVersion string
	// 31040 STR 15
	primaryDSPVersion string
	// 31055 STR 15
	slaveDSPVersion string
	// 31070 STR 15
	cplDRevNo string
	// 31085 STR 15
	afciVersion string
	// 31100 STR 15
	builtinPID string
}

func (x *hardwareData5) parse(data []byte) (err error) {
	size := 15*7 + 10
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var idx uint
	x.hardwareVersion, idx, _ = getSTR(data, idx, 15)
	x.monitoringBoardSN, idx, _ = getSTR(data, idx, 10)
	x.monitoringSoftwareVersion, idx, _ = getSTR(data, idx, 15)
	x.primaryDSPVersion, idx, _ = getSTR(data, idx, 15)
	x.slaveDSPVersion, idx, _ = getSTR(data, idx, 15)
	x.cplDRevNo, idx, _ = getSTR(data, idx, 15)
	x.afciVersion, idx, _ = getSTR(data, idx, 15)
	x.builtinPID, _, _ = getSTR(data, idx, 15)

	return nil
}

func (x *hardwareData5) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Hardware Data Part 5\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Hardware Version            = %q\n", x.hardwareVersion))
	sb.WriteString(fmt.Sprintf("# Monitoring Board SN         = %q\n", x.monitoringBoardSN))
	sb.WriteString(fmt.Sprintf("# Monitoring Software Version = %q\n", x.monitoringSoftwareVersion))
	sb.WriteString(fmt.Sprintf("# Primary DSP Version         = %q\n", x.primaryDSPVersion))
	sb.WriteString(fmt.Sprintf("# Slave DSP Version           = %q\n", x.slaveDSPVersion))
	sb.WriteString(fmt.Sprintf("# CPL D Rev. No.              = %q\n", x.cplDRevNo))
	sb.WriteString(fmt.Sprintf("# AFCI Version                = %q\n", x.afciVersion))
	sb.WriteString(fmt.Sprintf("# Builtin PID                 = %q\n", x.builtinPID))
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No hardware or identification data read yet\n")
	}
	// else {
	// 	// no useful metrics here
	// 	id.RLock()
	// 	defer id.RUnlock()
	// }
	sb.WriteString("\n")

	return sb.String()
}

type hardwareData6 struct {
	genericData

	// 31130 STR 15
	elModuleSoftwareVersion string
	// 31145 STR 15
	afci2SoftwareVersion string
}

func (x *hardwareData6) parse(data []byte) (err error) {
	size := 15 * 2
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var idx uint
	x.elModuleSoftwareVersion, idx, _ = getSTR(data, idx, 15)
	x.afci2SoftwareVersion, _, _ = getSTR(data, idx, 15)

	return nil
}

func (x *hardwareData6) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Hardware Data Part 6\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# EL Module Software Version = %q\n", x.elModuleSoftwareVersion))
	sb.WriteString(fmt.Sprintf("# AFCI2 Software Version      = %q\n", x.afci2SoftwareVersion))
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No hardware or identification data read yet\n")
	}
	// else {
	// 	// no useful metrics here
	// 	id.RLock()
	// 	defer id.RUnlock()
	// }
	sb.WriteString("\n")

	return sb.String()
}

type remoteSignallingData struct {
	genericData

	// 32000 Bitfield16 1
	singleMachineTelesignalling uint16
	// 32001 Bitfield16 1
	runningStatusMonitoringProcessing uint16
	// 32002 Bitfield16 1
	runningStatusPowerProcessing uint16
}

func (x *remoteSignallingData) parse(data []byte) (err error) {
	if len(data) < 4 {
		return fmt.Errorf("data length %d < 4", len(data))
	}

	x.Lock()
	defer x.Unlock()

	var idx uint
	x.singleMachineTelesignalling, idx, _ = getU16(data, idx)
	x.runningStatusMonitoringProcessing, idx, _ = getU16(data, idx)
	x.runningStatusPowerProcessing, _, _ = getU16(data, idx)

	return nil
}

func (x *remoteSignallingData) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Remote Signalling Data\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Single Machine Telesignalling         = %#04x\t%#016b\n", x.singleMachineTelesignalling, x.singleMachineTelesignalling))
	sb.WriteString(fmt.Sprintf("# Running Status Monitoring Processing  = %#04x\t%#016b\n", x.runningStatusMonitoringProcessing, x.runningStatusMonitoringProcessing))
	sb.WriteString(fmt.Sprintf("# Running Status Power Processing       = %#04x\t%#016b\n", x.runningStatusPowerProcessing, x.runningStatusPowerProcessing))
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No remote signalling or identification data read yet\n")
	}
	// else {
	// 	// no useful metrics here
	// 	id.RLock()
	// 	defer id.RUnlock()
	// }
	sb.WriteString("\n")

	return sb.String()
}

const alarmCount = 5

type alarmData struct {
	genericData

	// 32008-32010 Bitfield16 1x3 - but Alarms chapters defines also Alarm4 and Alarm 5
	alarm [alarmCount]uint16
}

func (x *alarmData) parse(data []byte) (err error) {
	if len(data) < 6 {
		return fmt.Errorf("data length %d < 6", len(data))
	}

	x.Lock()
	defer x.Unlock()

	var idx uint
	for i := 0; i < 3; i++ {
		x.alarm[i], idx, _ = getU16(data, idx)
	}

	return nil
}

func (x *alarmData) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Alarm Data\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	for i, v := range x.alarm {
		sb.WriteString(fmt.Sprintf("# Alarm %d = %#04x\t%#016b\n", i+1, v, v))
	}

	alarms := getAlarms(x.alarm)
	for _, a := range alarms {
		sb.WriteString(fmt.Sprintf("# Alarm Triggered: %s\n", a))
	}

	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No alarm or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()

		for i, a := range x.alarm {
			sb.WriteString(fmt.Sprintf("sun2000_inverter_alarm{model=%q,sn=%q,name=\"Alarm%d\"} %d\n", id.model, id.sn, i+1, a))
		}

		// might be a bit much, but nice for some historical visibility
		for _, a := range sun2000Alarms {
			value := "0"
			if a.isTriggered(x.alarm) {
				value = "1"
			}
			sb.WriteString(fmt.Sprintf("sun2000_inverter_alarm_triggered{model=%q,sn=%q,name=%q,id=%d,level=%q} %s\n", id.model, id.sn, a.name, a.id, a.level, value))
		}

	}
	sb.WriteString("\n")

	return sb.String()
}

type sun2000AlarmLevel uint8

const (
	alarmLevelWarning sun2000AlarmLevel = iota
	alarmLevelMinor
	alarmLevelMajor
)

func (x sun2000AlarmLevel) String() string {
	switch x {
	case alarmLevelWarning:
		return "Warning"
	case alarmLevelMinor:
		return "Minor"
	case alarmLevelMajor:
		return "Major"
	default:
		return "Unknown"
	}
}

type sun2000Alarm struct {
	mask  [alarmCount]uint16
	name  string
	id    uint16
	level sun2000AlarmLevel
}

// TODO: figure out if bit0 is LSB, or MSB
var sun2000Alarms = []sun2000Alarm{
	{[alarmCount]uint16{0b1000000000000000}, "High String Input Voltage", 2001, alarmLevelMajor},
	{[alarmCount]uint16{0b0100000000000000}, "DC Arc Fault", 2002, alarmLevelMajor},
	{[alarmCount]uint16{0b0010000000000000}, "String Reverse Connection", 2011, alarmLevelMajor},
	{[alarmCount]uint16{0b0001000000000000}, "String Current Backfeed ", 2012, alarmLevelWarning},
	{[alarmCount]uint16{0b0000100000000000}, "Abnormal String Power", 2013, alarmLevelWarning},
	{[alarmCount]uint16{0b0000010000000000}, "AFCI Self-Check Fail", 2021, alarmLevelMajor},
	{[alarmCount]uint16{0b0000001000000000}, "Phase Wire Short-Circuited to PE", 2031, alarmLevelMajor},
	{[alarmCount]uint16{0b0000000100000000}, "Grid Loss", 2032, alarmLevelMajor},
	{[alarmCount]uint16{0b0000000010000000}, "Grid Undervoltage", 2033, alarmLevelMajor},
	{[alarmCount]uint16{0b0000000001000000}, "Grid Overvoltage", 2034, alarmLevelMajor},
	{[alarmCount]uint16{0b0000000000100000}, "Grid Volt. Imbalance", 2035, alarmLevelMajor},
	{[alarmCount]uint16{0b0000000000010000}, "Grid Overfrequency", 2036, alarmLevelMajor},
	{[alarmCount]uint16{0b0000000000001000}, "Grid Underfrequency", 2037, alarmLevelMajor},
	{[alarmCount]uint16{0b0000000000000100}, "Unstable Grid Frequency", 2038, alarmLevelMajor},
	{[alarmCount]uint16{0b0000000000000010}, "Output Overcurrent", 2039, alarmLevelMajor},
	{[alarmCount]uint16{0b0000000000000001}, "Output DC Component Overhigh", 2040, alarmLevelMajor},

	{[alarmCount]uint16{0, 0b1000000000000000}, "Abnormal Residual Current", 2051, alarmLevelMajor},
	{[alarmCount]uint16{0, 0b0100000000000000}, "Abnormal Grounding", 2061, alarmLevelMajor},
	{[alarmCount]uint16{0, 0b0010000000000000}, "Low Insulation Resistance", 2062, alarmLevelMajor},
	{[alarmCount]uint16{0, 0b0001000000000000}, "Overtemperature", 2063, alarmLevelMajor},
	{[alarmCount]uint16{0, 0b0000100000000000}, "Device Fault", 2064, alarmLevelMajor},
	{[alarmCount]uint16{0, 0b0000010000000000}, "Upgrade Failed or Version Mismatch", 2065, alarmLevelMinor},
	{[alarmCount]uint16{0, 0b0000001000000000}, "License Expired", 2066, alarmLevelWarning},
	{[alarmCount]uint16{0, 0b0000000100000000}, "Faulty Monitoring Unit", 61440, alarmLevelMinor},
	{[alarmCount]uint16{0, 0b0000000010000000}, "Faulty Power Collector", 2067, alarmLevelMajor},
	{[alarmCount]uint16{0, 0b0000000001000000}, "Battery Abnormal", 2068, alarmLevelMinor},
	{[alarmCount]uint16{0, 0b0000000000100000}, "Active Islanding", 2070, alarmLevelMajor},
	{[alarmCount]uint16{0, 0b0000000000010000}, "Passive Islanding", 2071, alarmLevelMajor},
	{[alarmCount]uint16{0, 0b0000000000001000}, "Transient AC Overvoltage", 2072, alarmLevelMajor},
	{[alarmCount]uint16{0, 0b0000000000000100}, "Peripheral Port Short Circuit", 2075, alarmLevelWarning},
	{[alarmCount]uint16{0, 0b0000000000000010}, "Churn Output Overload", 2077, alarmLevelMajor},
	{[alarmCount]uint16{0, 0b0000000000000001}, "Abnormal PV Module Configuration", 2080, alarmLevelMajor},

	{[alarmCount]uint16{0, 0, 0b1000000000000000}, "Optimizer Fault", 2081, alarmLevelWarning},
	{[alarmCount]uint16{0, 0, 0b0100000000000000}, "Built-in PID Operation Abnormal", 2085, alarmLevelMajor},
	{[alarmCount]uint16{0, 0, 0b0010000000000000}, "High Input String Voltage to Ground", 2014, alarmLevelMajor},
	{[alarmCount]uint16{0, 0, 0b0001000000000000}, "External Fan Abnormal", 2086, alarmLevelMajor},
	{[alarmCount]uint16{0, 0, 0b0000100000000000}, "Battery Reverse Connection", 2069, alarmLevelMajor},
	{[alarmCount]uint16{0, 0, 0b0000010000000000}, "On-grid/Off-grid Controller Abnormal", 2082, alarmLevelMajor},
	{[alarmCount]uint16{0, 0, 0b0000001000000000}, "PV String Loss", 2015, alarmLevelWarning},
	{[alarmCount]uint16{0, 0, 0b0000000100000000}, "Internal Fan Abnormal", 2087, alarmLevelMajor},
	{[alarmCount]uint16{0, 0, 0b0000000010000000}, "DC Protection Unit Abnormal", 2088, alarmLevelMajor},

	{[alarmCount]uint16{0, 0, 0, 0b0000000000100000}, "Management System Cert Valid Time Ineffective", 2095, alarmLevelMajor},
	{[alarmCount]uint16{0, 0, 0, 0b0000000000010000}, "Management System Cert Valid Time Being Overdue", 2096, alarmLevelMajor},
	{[alarmCount]uint16{0, 0, 0, 0b0000000000001000}, "Management System Cert Valid Time Overdue", 2097, alarmLevelMajor},

	{[alarmCount]uint16{0, 0, 0, 0, 0b0001000000000000}, "CT Disconnection", 2067, alarmLevelMajor},
	{[alarmCount]uint16{0, 0, 0, 0, 0b0000100000000000}, "PT Disconnection", 2067, alarmLevelMajor},
}

func (x sun2000Alarm) isTriggered(alarm [alarmCount]uint16) bool {
	for i := 0; i < alarmCount; i++ {
		if x.mask[i]&alarm[i] != 0 {
			return true
		}
	}
	return false
}

func getAlarms(id [alarmCount]uint16) (out []sun2000Alarm) {
	out = make([]sun2000Alarm, 0, 8)
	for _, a := range sun2000Alarms {
		for i := 0; i < alarmCount; i++ {
			if a.mask[i]&id[i] != 0 {
				out = append(out, a)
				break
			}
		}
	}
	return out
}

func (x sun2000Alarm) String() string {
	return fmt.Sprintf("%s : %s id=%d", x.level, x.name, x.id)
}

type pvData struct {
	genericData

	// 32015 U16 1
	deviceSNSignatureCode uint16
	// This could be either pv[x].voltage/current, or MPPT1[x].voltage/current
	pv [20]struct {
		// 32016 I16 1 gain 10 V
		voltage float32
		// 32017 I16 1 gain 100 A
		current float32
	}
}

func (x *pvData) parse(data []byte) (err error) {
	if len(data) < 2+2*20*2 {
		return fmt.Errorf("data length %d < 2+2*20*2", len(data))
	}

	x.Lock()
	defer x.Unlock()

	var idx uint
	x.deviceSNSignatureCode, idx, _ = getU16(data, idx)
	var i16 int16
	for i := 0; i < 20; i++ {
		i16, idx, _ = getI16(data, idx)
		x.pv[i].voltage = float32(i16) / 10
		i16, idx, _ = getI16(data, idx)
		x.pv[i].current = float32(i16) / 100
	}

	return nil
}

func (x *pvData) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# PV Data\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Device SN Signature Code = %#04x\t%#016b\n", x.deviceSNSignatureCode, x.deviceSNSignatureCode))
	sb.WriteString("\n")
	powerTotal := float32(0)
	for i, v := range x.pv {
		if v.voltage != 0 || v.current != 0 || i < int(id.numberOfStrings) {
			sb.WriteString(fmt.Sprintf("# PV%2d Voltage   = %7.3f V\n", i+1, v.voltage))
			sb.WriteString(fmt.Sprintf("# PV%2d Current   = %7.3f A\n", i+1, v.current))
			power := v.voltage * v.current
			powerTotal += power
			sb.WriteString(fmt.Sprintf("# PV%2d Power     = %7.3f kW\n", i+1, power/1000))
		}
	}
	sb.WriteString("# --------------------------\n")
	sb.WriteString(fmt.Sprintf("# PV Total Power = %7.3f kW\n", powerTotal/1000))
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No PV or identification data read yet\n")
	} else {
		// no useful metrics here
		id.RLock()
		defer id.RUnlock()

		for i, v := range x.pv {
			if v.voltage != 0 || v.current != 0 || i < int(id.numberOfStrings) {
				sb.WriteString(fmt.Sprintf("sun2000_inverter_pv_voltage{model=%q,sn=%q,pv=%d,unit=\"V\"} %.1f\n", id.model, id.sn, i+1, v.voltage))
				sb.WriteString(fmt.Sprintf("sun2000_inverter_pv_current{model=%q,sn=%q,pv=%d,unit=\"A\"} %.2f\n", id.model, id.sn, i+1, v.current))
				power := v.voltage * v.current
				sb.WriteString(fmt.Sprintf("sun2000_inverter_pv_power{model=%q,sn=%q,pv=%d,unit=\"kW\"} %.3f\n", id.model, id.sn, i+1, power/1000))
			}
		}
		sb.WriteString(fmt.Sprintf("sun2000_inverter_pv_total_power{model=%q,sn=%q,unit=\"kW\"} %.3f\n", id.model, id.sn, powerTotal/1000))
	}
	sb.WriteString("\n")

	return sb.String()
}

type gridData struct {
	genericData

	// 32064 I32 2 gain 1000 kW
	dcPower float32

	// 32066 U16 1 gain 10 V
	powergridABLineVoltage float32
	// 32067 U16 1 gain 10 V
	powergridBCLineVoltage float32
	// 32068 U16 1 gain 10 V
	powergridCALineVoltage float32

	// 32069 U16 1 gain 10 V
	powergridPhaseAVoltage float32
	// 32070 U16 1 gain 10 V
	powergridPhaseBVoltage float32
	// 32071 U16 1 gain 10 V
	powergridPhaseCVoltage float32

	// 32072 I32 2 gain 1000 A
	powergridPhaseACurrent float32
	// 32074 I32 2 gain 1000 A
	powergridPhaseBCurrent float32
	// 32076 I32 2 gain 1000 A
	powergridPhaseCCurrent float32

	// 32078 I32 2 gain 1000 kW
	peakActivePowerOfTheDay float32
	// 32080 I32 2 gain 1000 kW
	activePower float32
	// 32082 I32 2 gain 1000 kVar
	reactivePower float32
	// 32084 I16 1 gain 1000
	powerFactor float32

	// 32085 U16 1 gain 100 Hz
	powergridFrequency float32

	// 32086 U16 1 gain 100 %
	inverterEfficiency float32

	// 32087 I16 1 gain 10 ℃
	internalTemperature float32

	// 32088 U16 1 gain 1000 MΩ
	insulationImpedanceValue float32

	// 32089 E16 1
	deviceStatus deviceStatus

	// 32090 U16 1
	faultCode uint16

	// 32091 epoch 2
	startupTime time.Time
	// 32093 epoch 2
	shutdownTime time.Time

	// 32095 I32 2 gain 1000 kW
	activePowerFast float32
}

func (x *gridData) parse(data []byte) (err error) {
	size := 2 + 3 + 3 + 6 + 7 + 1*6 + 2*2 + 2
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var i32 int32
	var i16 int16
	var u16 uint16
	var epoch uint32
	var idx uint

	i32, idx, _ = getI32(data, idx)
	x.dcPower = float32(i32) / 1000

	u16, idx, _ = getU16(data, idx)
	x.powergridABLineVoltage = float32(u16) / 10
	u16, idx, _ = getU16(data, idx)
	x.powergridBCLineVoltage = float32(u16) / 10
	u16, idx, _ = getU16(data, idx)
	x.powergridCALineVoltage = float32(u16) / 10

	u16, idx, _ = getU16(data, idx)
	x.powergridPhaseAVoltage = float32(u16) / 10
	u16, idx, _ = getU16(data, idx)
	x.powergridPhaseBVoltage = float32(u16) / 10
	u16, idx, _ = getU16(data, idx)
	x.powergridPhaseCVoltage = float32(u16) / 10

	i32, idx, _ = getI32(data, idx)
	x.powergridPhaseACurrent = float32(i32) / 1000
	i32, idx, _ = getI32(data, idx)
	x.powergridPhaseBCurrent = float32(i32) / 1000
	i32, idx, _ = getI32(data, idx)
	x.powergridPhaseCCurrent = float32(i32) / 1000

	i32, idx, _ = getI32(data, idx)
	x.peakActivePowerOfTheDay = float32(i32) / 1000
	i32, idx, _ = getI32(data, idx)
	x.activePower = float32(i32) / 1000

	i32, idx, _ = getI32(data, idx)
	x.reactivePower = float32(i32) / 1000
	i16, idx, _ = getI16(data, idx)
	x.powerFactor = float32(i16) / 1000

	u16, idx, _ = getU16(data, idx)
	x.powergridFrequency = float32(u16) / 100

	u16, idx, _ = getU16(data, idx)
	x.inverterEfficiency = float32(u16) / 100

	i16, idx, _ = getI16(data, idx)
	x.internalTemperature = float32(i16) / 10

	u16, idx, _ = getU16(data, idx)
	x.insulationImpedanceValue = float32(u16) / 1000

	u16, idx, _ = getU16(data, idx)
	x.deviceStatus = deviceStatus(u16)

	x.faultCode, idx, _ = getU16(data, idx)

	epoch, idx, _ = getU32(data, idx)
	x.startupTime = time.Unix(int64(epoch), 0)
	epoch, idx, _ = getU32(data, idx)
	x.shutdownTime = time.Unix(int64(epoch), 0)

	i32, _, _ = getI32(data, idx)
	x.activePowerFast = float32(i32) / 1000

	return nil
}

func (x *gridData) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Grid Data\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# DC Power                       = %2.3f kW\n", x.dcPower))
	sb.WriteString("#\n")
	sb.WriteString(fmt.Sprintf("# Powergrid AB Line Voltage      = %3.1f V\n", x.powergridABLineVoltage))
	sb.WriteString(fmt.Sprintf("# Powergrid BC Line Voltage      = %3.1f V\n", x.powergridBCLineVoltage))
	sb.WriteString(fmt.Sprintf("# Powergrid CA Line Voltage      = %3.1f V\n", x.powergridCALineVoltage))
	sb.WriteString("#\n")
	sb.WriteString(fmt.Sprintf("# Powergrid Phase A Voltage      = %3.1f V\n", x.powergridPhaseAVoltage))
	sb.WriteString(fmt.Sprintf("# Powergrid Phase B Voltage      = %3.1f V\n", x.powergridPhaseBVoltage))
	sb.WriteString(fmt.Sprintf("# Powergrid Phase C Voltage      = %3.1f V\n", x.powergridPhaseCVoltage))
	sb.WriteString("#\n")
	sb.WriteString(fmt.Sprintf("# Powergrid Phase A Current      = %3.3f A\n", x.powergridPhaseACurrent))
	sb.WriteString(fmt.Sprintf("# Powergrid Phase B Current      = %3.3f A\n", x.powergridPhaseBCurrent))
	sb.WriteString(fmt.Sprintf("# Powergrid Phase C Current      = %3.3f A\n", x.powergridPhaseCCurrent))
	sb.WriteString("#\n")
	pA := x.powergridPhaseAVoltage * x.powergridPhaseACurrent
	pB := x.powergridPhaseBVoltage * x.powergridPhaseBCurrent
	pC := x.powergridPhaseCVoltage * x.powergridPhaseCCurrent
	sb.WriteString(fmt.Sprintf("# Powergrid Phase A Power        = %3.3f VA\n", pA))
	sb.WriteString(fmt.Sprintf("# Powergrid Phase B Power        = %3.3f VA\n", pB))
	sb.WriteString(fmt.Sprintf("# Powergrid Phase C Power        = %3.3f VA\n", pC))
	sb.WriteString(fmt.Sprintf("# Powergrid Total Power          = %3.3f VA\n", pA+pB+pC))
	sb.WriteString("#\n")
	sb.WriteString(fmt.Sprintf("# Peak Active Power of the Day   = %3.3f kW\n", x.peakActivePowerOfTheDay))
	sb.WriteString(fmt.Sprintf("# Active Power Fast              = %3.3f kW\n", x.activePowerFast))
	sb.WriteString(fmt.Sprintf("# Active Power                   = %3.3f kW\n", x.activePower))
	sb.WriteString(fmt.Sprintf("# Reactive Power                 = %3.3f kVar\n", x.reactivePower))
	sb.WriteString(fmt.Sprintf("# Power Factor                   = %3.3f\n", x.powerFactor))
	sb.WriteString(fmt.Sprintf("# Powergrid Frequency            = %2.2f Hz\n", x.powergridFrequency))
	sb.WriteString(fmt.Sprintf("# Inverter Efficiency            = %3.2f %%\n", x.inverterEfficiency))
	sb.WriteString("#\n")
	sb.WriteString(fmt.Sprintf("# Internal Temperature           = %3.1f ℃\n", x.internalTemperature))
	sb.WriteString(fmt.Sprintf("# Insulation Impedance Value     = %4.3f MΩ\n", x.insulationImpedanceValue))
	sb.WriteString(fmt.Sprintf("# Device Status                  = %d\t%s\n", x.deviceStatus, x.deviceStatus))
	sb.WriteString(fmt.Sprintf("# Fault Code                     = %#04x\t%#016b\n", x.faultCode, x.faultCode))
	sb.WriteString("#\n")
	sb.WriteString(fmt.Sprintf("# Startup Time                   = %s\n", x.startupTime.Format(time.RFC3339)))
	if x.shutdownTime.Unix() == 4294967295 {
		sb.WriteString("# Shutdown Time                  = N/A\n")
	} else {
		sb.WriteString(fmt.Sprintf("# Shutdown Time                  = %s\n", x.shutdownTime.Format(time.RFC3339)))
	}
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No grid or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()

		sb.WriteString(fmt.Sprintf("sun2000_inverter_dc_power{model=%q,sn=%q,unit=\"kW\"} %3.3f\n", id.model, id.sn, x.dcPower))

		sb.WriteString(fmt.Sprintf("sun2000_inverter_powergrid_ab_line_voltage{model=%q,sn=%q,unit=\"V\"} %3.1f\n", id.model, id.sn, x.powergridABLineVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_powergrid_bc_line_voltage{model=%q,sn=%q,unit=\"V\"} %3.1f\n", id.model, id.sn, x.powergridBCLineVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_powergrid_ca_line_voltage{model=%q,sn=%q,unit=\"V\"} %3.1f\n", id.model, id.sn, x.powergridCALineVoltage))

		sb.WriteString(fmt.Sprintf("sun2000_inverter_powergrid_phase_a_voltage{model=%q,sn=%q,unit=\"V\"} %3.1f\n", id.model, id.sn, x.powergridPhaseAVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_powergrid_phase_b_voltage{model=%q,sn=%q,unit=\"V\"} %3.1f\n", id.model, id.sn, x.powergridPhaseBVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_powergrid_phase_c_voltage{model=%q,sn=%q,unit=\"V\"} %3.1f\n", id.model, id.sn, x.powergridPhaseCVoltage))

		sb.WriteString(fmt.Sprintf("sun2000_inverter_powergrid_phase_a_current{model=%q,sn=%q,unit=\"A\"} %3.3f\n", id.model, id.sn, x.powergridPhaseACurrent))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_powergrid_phase_b_current{model=%q,sn=%q,unit=\"A\"} %3.3f\n", id.model, id.sn, x.powergridPhaseBCurrent))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_powergrid_phase_c_current{model=%q,sn=%q,unit=\"A\"} %3.3f\n", id.model, id.sn, x.powergridPhaseCCurrent))

		sb.WriteString(fmt.Sprintf("sun2000_inverter_powergrid_phase_a_power{model=%q,sn=%q,unit=\"VA\"} %3.3f\n", id.model, id.sn, pA))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_powergrid_phase_b_power{model=%q,sn=%q,unit=\"VA\"} %3.3f\n", id.model, id.sn, pB))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_powergrid_phase_c_power{model=%q,sn=%q,unit=\"VA\"} %3.3f\n", id.model, id.sn, pC))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_powergrid_total_power{model=%q,sn=%q,unit=\"VA\"} %3.3f\n", id.model, id.sn, pA+pB+pC))

		sb.WriteString(fmt.Sprintf("sun2000_inverter_peak_active_power_of_the_day{model=%q,sn=%q,unit=\"kW\"} %3.3f\n", id.model, id.sn, x.peakActivePowerOfTheDay))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_active_power_fast{model=%q,sn=%q,unit=\"kW\"} %3.3f\n", id.model, id.sn, x.activePowerFast))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_active_power{model=%q,sn=%q,unit=\"kW\"} %3.3f\n", id.model, id.sn, x.activePower))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_reactive_power{model=%q,sn=%q,unit=\"kVar\"} %3.3f\n", id.model, id.sn, x.reactivePower))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_power_factor{model=%q,sn=%q} %3.3f\n", id.model, id.sn, x.powerFactor))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_powergrid_frequency{model=%q,sn=%q,unit=\"Hz\"} %2.2f\n", id.model, id.sn, x.powergridFrequency))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_inverter_efficiency{model=%q,sn=%q,unit=\"%%\"} %3.2f\n", id.model, id.sn, x.inverterEfficiency))

		sb.WriteString(fmt.Sprintf("sun2000_inverter_internal_temperature{model=%q,sn=%q,unit=\"℃\"} %3.1f\n", id.model, id.sn, x.internalTemperature))

		sb.WriteString(fmt.Sprintf("sun2000_inverter_insulation_impedance_value{model=%q,sn=%q,unit=\"MΩ\"} %4.3f\n", id.model, id.sn, x.insulationImpedanceValue))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_device_status{model=%q,sn=%q} %d\n", id.model, id.sn, x.deviceStatus))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_fault_code{model=%q,sn=%q} %d\n", id.model, id.sn, x.faultCode))

		sb.WriteString(fmt.Sprintf("sun2000_inverter_startup_time{model=%q,sn=%q} %d\n", id.model, id.sn, x.startupTime.Unix()))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_shutdown_time{model=%q,sn=%q} %d\n", id.model, id.sn, x.shutdownTime.Unix()))

	}
	sb.WriteString("\n")

	return sb.String()
}

type deviceStatus uint16

func (x deviceStatus) String() string {
	switch x {
	case 0:
		return "Standby: initializing"
	case 1:
		return "Standby: detecting insulation resistance"
	case 2:
		return "Standby: detecting irradiation"
	case 3:
		return "Standby: grid detecting"
	case 256:
		return "Starting"
	case 512:
		return "On-grid: running"
	case 513:
		return "Grid connection: power limited"
	case 514:
		return "Grid connection: self-derating"
	case 515:
		return "Off-grid Running"
	case 768:
		return "Shutdown: fault"
	case 769:
		return "Shutdown: command"
	case 770:
		return "Shutdown: OVGR"
	case 771:
		return "Shutdown: communication disconnected"
	case 772:
		return "Shutdown: power limited"
	case 773:
		return "Shutdown: manual startup required"
	case 774:
		return "Shutdown: DC switches disconnected"
	case 775:
		return "Shutdown: rapid cutoff"
	case 776:
		return "Shutdown: input underpower"
	case 1025:
		return "Grid scheduling: cosΦ-P curve"
	case 1026:
		return "Grid scheduling: Q-U curve"
	case 1027:
		return "Grid scheduling: PF-U curve"
	case 1028:
		return "Grid scheduling: dry contact"
	case 1029:
		return "Grid scheduling: Q-P curve"
	case 1280:
		return "Spot-check ready"
	case 1281:
		return "Spot-checking"
	case 1536:
		return "Inspecting"
	case 1792:
		return "AFCI self check"
	case 2048:
		return "I-V scanning"
	case 2304:
		return "DC input detection"
	case 2560:
		return "Running: off-grid charging"
	case 40960:
		return "Standby: no irradiation"
	default:
		return fmt.Sprintf("Unknown status %d %#04x\t%#016b", x, x, x)
	}
}
