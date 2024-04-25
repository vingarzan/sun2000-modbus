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
	identification      identificationData
	product             productData
	hardware1           hardwareData1
	hardware2           hardwareData2
	hardware3           hardwareData3
	hardware4           hardwareData4
	hardware5           hardwareData5
	hardware6           hardwareData6
	remoteSignalling    remoteSignallingData
	alarm1              alarmData1
	pv                  pvData
	inverter            inverterData
	cumulative1         cumulativeData1
	cumulative2         cumulativeData2
	cumulative3         cumulativeData3
	mppt1               mpptData1
	alarm2              alarmData2
	stringAccess        stringAccessData
	mppt2               mpptData2
	internalTemperature internalTemperatureData
	meter               meterData
	esu1                esu1Data
	esu2                esu2Data
	esuTemperatures     esuTemperaturesData
}

func init() {
	for i := 0; i < 3; i++ {
		parsedData.esu1.pack[i].esuId = 1
		parsedData.esu1.pack[i].id = i + 1
		parsedData.esu2.pack[i].esuId = 2
		parsedData.esu2.pack[i].id = i + 1
	}
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
	sb.WriteString(x.alarm1.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.pv.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.inverter.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.cumulative1.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.cumulative2.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.cumulative3.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.mppt1.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.alarm2.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.stringAccess.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.mppt2.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.internalTemperature.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.meter.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.esu1.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.esu2.metricsString(&x.identification))
	sb.WriteString("\n")
	sb.WriteString(x.esuTemperatures.metricsString(&x.identification))
	sb.WriteString("\n")

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
		sb.WriteString(fmt.Sprintf("sun2000_number_of_MPPTs{model=%q,sn=%q} %d\n", x.model, x.sn, x.numberOfMPPTs))
		sb.WriteString(fmt.Sprintf("sun2000_number_of_strings{model=%q,sn=%q} %d\n", x.model, x.sn, x.numberOfStrings))
		sb.WriteString(fmt.Sprintf("sun2000_rated_power{model=%q,sn=%q,unit=\"kW\",description=\"Rated Power\"} %.3f\n", x.model, x.sn, x.ratedPower))
		sb.WriteString(fmt.Sprintf("sun2000_Pmax{model=%q,sn=%q,unit=\"kW\",description=\"Maximum Active Power Pmax\"} %.3f\n", x.model, x.sn, x.maxActivePowerPmax))
		sb.WriteString(fmt.Sprintf("sun2000_Smax{model=%q,sn=%q,unit=\"kVA\",description=\"Maximum Apparent Power Smax\"} %.3f\n", x.model, x.sn, x.maxApparentPowerSmax))
		sb.WriteString(fmt.Sprintf("sun2000_Qmax_feed_to_grid{model=%q,sn=%q,unit=\"kVar\",description=\"Realtime Max Reactive Power Qmax Feed to Grid\"} %.3f\n", x.model, x.sn, x.realtimeMaxReactivePowerQmaxFeedToGrid))
		sb.WriteString(fmt.Sprintf("sun2000_Qmax_absorbed_from_grid{model=%q,sn=%q,unit=\"kVar\",description=\"Realtime Max Reactive Power Qmax Absorbed from Grid\"} %.3f\n", x.model, x.sn, x.realtimeMaxReactivePowerQmaxAbsorbedFromGrid))
		sb.WriteString(fmt.Sprintf("sun2000_Pmax_real{model=%q,sn=%q,unit=\"kW\",description=\"Maximum Active Capability Pmax Real\"} %.3f\n", x.model, x.sn, x.maxActiveCapabilityPmaxReal))
		sb.WriteString(fmt.Sprintf("sun2000_Smax_real{model=%q,sn=%q,unit=\"kVA\",description=\"Maximum Apparent Capability Smax Real\"} %.3f\n", x.model, x.sn, x.maxApparentCapabilitySmaxReal))

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
		sb.WriteString(fmt.Sprintf("sun2000_unique_id_of_the_software{model=%q,sn=%q} %d\n", id.model, id.sn, x.uniqueIDOfTheSoftware))
		sb.WriteString(fmt.Sprintf("sun2000_number_of_packages_to_be_upgraded{model=%q,sn=%q} %d\n", id.model, id.sn, x.numberOfPackagesToBeUpgraded))
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

type alarmData1 struct {
	genericData

	// 32008-32010 Bitfield16 1x3 - but Alarms chapters defines also Alarm4 and Alarm 5
	alarm [alarmCount]uint16
}

func (x *alarmData1) parse(data []byte) (err error) {
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

func (x *alarmData1) metricsString(id *identificationData) string {
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
			sb.WriteString(fmt.Sprintf("sun2000_alarm{model=%q,sn=%q,name=\"Alarm%d\"} %d\n", id.model, id.sn, i+1, a))
		}

		// might be a bit much, but nice for some historical visibility
		for _, a := range sun2000Alarms {
			value := "0"
			if a.isTriggered(x.alarm) {
				value = "1"
			}
			sb.WriteString(fmt.Sprintf("sun2000_alarm_triggered{model=%q,sn=%q,name=%q,id=\"%d\",level=%q} %s\n", id.model, id.sn, a.name, a.id, a.level, value))
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
				sb.WriteString(fmt.Sprintf("sun2000_pv_voltage{model=%q,sn=%q,pv=\"%d\",unit=\"V\"} %.1f\n", id.model, id.sn, i+1, v.voltage))
				sb.WriteString(fmt.Sprintf("sun2000_pv_current{model=%q,sn=%q,pv=\"%d\",unit=\"A\"} %.2f\n", id.model, id.sn, i+1, v.current))
				power := v.voltage * v.current
				sb.WriteString(fmt.Sprintf("sun2000_pv_power{model=%q,sn=%q,pv=\"%d\",unit=\"kW\"} %.3f\n", id.model, id.sn, i+1, power/1000))
			}
		}
		sb.WriteString(fmt.Sprintf("sun2000_pv_total_power{model=%q,sn=%q,unit=\"kW\"} %.3f\n", id.model, id.sn, powerTotal/1000))
	}
	sb.WriteString("\n")

	return sb.String()
}

type inverterData struct {
	genericData

	// 32064 I32 2 gain 1000 kW
	dcPower float32

	// 32066 U16 1 gain 10 V
	inverterABLineVoltage float32
	// 32067 U16 1 gain 10 V
	inverterBCLineVoltage float32
	// 32068 U16 1 gain 10 V
	inverterCALineVoltage float32

	// 32069 U16 1 gain 10 V
	inverterPhaseAVoltage float32
	// 32070 U16 1 gain 10 V
	inverterPhaseBVoltage float32
	// 32071 U16 1 gain 10 V
	inverterPhaseCVoltage float32

	// 32072 I32 2 gain 1000 A
	inverterPhaseACurrent float32
	// 32074 I32 2 gain 1000 A
	inverterPhaseBCurrent float32
	// 32076 I32 2 gain 1000 A
	inverterPhaseCCurrent float32

	// 32078 I32 2 gain 1000 kW
	peakActivePowerOfTheDay float32
	// 32080 I32 2 gain 1000 kW
	activePower float32
	// 32082 I32 2 gain 1000 kVar
	reactivePower float32
	// 32084 I16 1 gain 1000
	powerFactor float32

	// 32085 U16 1 gain 100 Hz
	inverterFrequency float32

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

func (x *inverterData) parse(data []byte) (err error) {
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
	x.inverterABLineVoltage = float32(u16) / 10
	u16, idx, _ = getU16(data, idx)
	x.inverterBCLineVoltage = float32(u16) / 10
	u16, idx, _ = getU16(data, idx)
	x.inverterCALineVoltage = float32(u16) / 10

	u16, idx, _ = getU16(data, idx)
	x.inverterPhaseAVoltage = float32(u16) / 10
	u16, idx, _ = getU16(data, idx)
	x.inverterPhaseBVoltage = float32(u16) / 10
	u16, idx, _ = getU16(data, idx)
	x.inverterPhaseCVoltage = float32(u16) / 10

	i32, idx, _ = getI32(data, idx)
	x.inverterPhaseACurrent = float32(i32) / 1000
	i32, idx, _ = getI32(data, idx)
	x.inverterPhaseBCurrent = float32(i32) / 1000
	i32, idx, _ = getI32(data, idx)
	x.inverterPhaseCCurrent = float32(i32) / 1000

	i32, idx, _ = getI32(data, idx)
	x.peakActivePowerOfTheDay = float32(i32) / 1000
	i32, idx, _ = getI32(data, idx)
	x.activePower = float32(i32) / 1000

	i32, idx, _ = getI32(data, idx)
	x.reactivePower = float32(i32) / 1000
	i16, idx, _ = getI16(data, idx)
	x.powerFactor = float32(i16) / 1000

	u16, idx, _ = getU16(data, idx)
	x.inverterFrequency = float32(u16) / 100

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

func (x *inverterData) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Grid Data\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# DC Power                       = %2.3f kW\n", x.dcPower))
	sb.WriteString("#\n")
	sb.WriteString(fmt.Sprintf("# Inverter AB Line Voltage       = %3.1f V\n", x.inverterABLineVoltage))
	sb.WriteString(fmt.Sprintf("# Inverter BC Line Voltage       = %3.1f V\n", x.inverterBCLineVoltage))
	sb.WriteString(fmt.Sprintf("# Inverter CA Line Voltage       = %3.1f V\n", x.inverterCALineVoltage))
	sb.WriteString("#\n")
	sb.WriteString(fmt.Sprintf("# Inverter Phase A Voltage       = %3.1f V\n", x.inverterPhaseAVoltage))
	sb.WriteString(fmt.Sprintf("# Inverter Phase B Voltage       = %3.1f V\n", x.inverterPhaseBVoltage))
	sb.WriteString(fmt.Sprintf("# Inverter Phase C Voltage       = %3.1f V\n", x.inverterPhaseCVoltage))
	sb.WriteString("#\n")
	sb.WriteString(fmt.Sprintf("# Inverter Phase A Current       = %3.3f A\n", x.inverterPhaseACurrent))
	sb.WriteString(fmt.Sprintf("# Inverter Phase B Current       = %3.3f A\n", x.inverterPhaseBCurrent))
	sb.WriteString(fmt.Sprintf("# Inverter Phase C Current       = %3.3f A\n", x.inverterPhaseCCurrent))
	sb.WriteString("#\n")
	pA := x.inverterPhaseAVoltage * x.inverterPhaseACurrent
	pB := x.inverterPhaseBVoltage * x.inverterPhaseBCurrent
	pC := x.inverterPhaseCVoltage * x.inverterPhaseCCurrent
	sb.WriteString(fmt.Sprintf("# Inverter Phase A Power         = %3.3f VA\n", pA))
	sb.WriteString(fmt.Sprintf("# Inverter Phase B Power         = %3.3f VA\n", pB))
	sb.WriteString(fmt.Sprintf("# Inverter Phase C Power         = %3.3f VA\n", pC))
	sb.WriteString(fmt.Sprintf("# Inverter Total Power           = %3.3f VA\n", pA+pB+pC))
	sb.WriteString("#\n")
	sb.WriteString(fmt.Sprintf("# Peak Active Power of the Day   = %3.3f kW\n", x.peakActivePowerOfTheDay))
	sb.WriteString(fmt.Sprintf("# Active Power Fast              = %3.3f kW\n", x.activePowerFast))
	sb.WriteString(fmt.Sprintf("# Active Power                   = %3.3f kW\n", x.activePower))
	sb.WriteString(fmt.Sprintf("# Reactive Power                 = %3.3f kVar\n", x.reactivePower))
	sb.WriteString(fmt.Sprintf("# Power Factor                   = %3.3f\n", x.powerFactor))
	sb.WriteString(fmt.Sprintf("# Inverter Frequency             = %2.2f Hz\n", x.inverterFrequency))
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

		sb.WriteString(fmt.Sprintf("sun2000_inverter_line_voltage{model=%q,sn=%q,line=\"AB\",unit=\"V\"} %3.1f\n", id.model, id.sn, x.inverterABLineVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_line_voltage{model=%q,sn=%q,line=\"BC\",unit=\"V\"} %3.1f\n", id.model, id.sn, x.inverterBCLineVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_line_voltage{model=%q,sn=%q,line=\"CA\",unit=\"V\"} %3.1f\n", id.model, id.sn, x.inverterCALineVoltage))

		sb.WriteString(fmt.Sprintf("sun2000_inverter_phase_voltage{model=%q,sn=%q,phase=\"A\",unit=\"V\"} %3.1f\n", id.model, id.sn, x.inverterPhaseAVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_phase_voltage{model=%q,sn=%q,phase=\"B\",unit=\"V\"} %3.1f\n", id.model, id.sn, x.inverterPhaseBVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_phase_voltage{model=%q,sn=%q,phase=\"C\",unit=\"V\"} %3.1f\n", id.model, id.sn, x.inverterPhaseCVoltage))

		sb.WriteString(fmt.Sprintf("sun2000_inverter_phase_current{model=%q,sn=%q,phase=\"A\",unit=\"A\"} %3.3f\n", id.model, id.sn, x.inverterPhaseACurrent))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_phase_current{model=%q,sn=%q,phase=\"B\",unit=\"A\"} %3.3f\n", id.model, id.sn, x.inverterPhaseBCurrent))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_phase_current{model=%q,sn=%q,phase=\"C\",unit=\"A\"} %3.3f\n", id.model, id.sn, x.inverterPhaseCCurrent))

		sb.WriteString(fmt.Sprintf("sun2000_inverter_phase_power{model=%q,sn=%q,phase=\"A\",unit=\"VA\"} %3.3f\n", id.model, id.sn, pA))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_phase_power{model=%q,sn=%q,phase=\"B\",unit=\"VA\"} %3.3f\n", id.model, id.sn, pB))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_phase_power{model=%q,sn=%q,phase=\"C\",unit=\"VA\"} %3.3f\n", id.model, id.sn, pC))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_total_power{model=%q,sn=%q,unit=\"VA\"} %3.3f\n", id.model, id.sn, pA+pB+pC))

		sb.WriteString(fmt.Sprintf("sun2000_inverter_peak_active_power_of_the_day{model=%q,sn=%q,unit=\"kW\"} %3.3f\n", id.model, id.sn, x.peakActivePowerOfTheDay))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_active_power_fast{model=%q,sn=%q,unit=\"kW\"} %3.3f\n", id.model, id.sn, x.activePowerFast))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_active_power{model=%q,sn=%q,unit=\"kW\"} %3.3f\n", id.model, id.sn, x.activePower))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_reactive_power{model=%q,sn=%q,unit=\"kVar\"} %3.3f\n", id.model, id.sn, x.reactivePower))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_power_factor{model=%q,sn=%q} %3.3f\n", id.model, id.sn, x.powerFactor))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_frequency{model=%q,sn=%q,unit=\"Hz\"} %2.2f\n", id.model, id.sn, x.inverterFrequency))
		sb.WriteString(fmt.Sprintf("sun2000_inverter_efficiency{model=%q,sn=%q,unit=\"%%\"} %3.2f\n", id.model, id.sn, x.inverterEfficiency))

		sb.WriteString(fmt.Sprintf("sun2000_internal_temperature{model=%q,sn=%q,sensor=\"main\",unit=\"℃\"} %3.1f\n", id.model, id.sn, x.internalTemperature))

		sb.WriteString(fmt.Sprintf("sun2000_insulation_impedance_value{model=%q,sn=%q,unit=\"MΩ\"} %4.3f\n", id.model, id.sn, x.insulationImpedanceValue))
		sb.WriteString(fmt.Sprintf("sun2000_device_status{model=%q,sn=%q,state=%q} %d\n", id.model, id.sn, x.deviceStatus, x.deviceStatus))
		sb.WriteString(fmt.Sprintf("sun2000_fault_code{model=%q,sn=%q} %d\n", id.model, id.sn, x.faultCode))

		sb.WriteString(fmt.Sprintf("sun2000_startup_time{model=%q,sn=%q} %d\n", id.model, id.sn, x.startupTime.Unix()))
		sb.WriteString(fmt.Sprintf("sun2000_shutdown_time{model=%q,sn=%q} %d\n", id.model, id.sn, x.shutdownTime.Unix()))

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

type cumulativeData1 struct {
	genericData

	// 32106 U32 2 gain 100 kWh
	cumulativeGeneratedElectricity float32
	// 32108 U32 2 gain 100 kWh
	totalDCInputPower float32
	// 32110 Epoch 2
	currentElectricityGenerationStatisticsTime time.Time
	// 32112 U32 2 gain 100 kWh
	electricityGeneratedInCurrentHour float32
	// 32114 U32 2 gain 100 kWh
	electricityGeneratedInCurrentDay float32
	// 32116 U32 2 gain 100 kWh
	electricityGeneratedInCurrentMonth float32
	// 32118 U32 2 gain 100 kWh
	electricityGeneratedInCurrentYear float32
}

func (x *cumulativeData1) parse(data []byte) (err error) {
	size := 7 * 2
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var u32 uint32
	var epoch uint32
	var idx uint

	u32, idx, _ = getU32(data, idx)
	x.cumulativeGeneratedElectricity = float32(u32) / 100

	u32, idx, _ = getU32(data, idx)
	x.totalDCInputPower = float32(u32) / 100

	epoch, idx, _ = getU32(data, idx)
	x.currentElectricityGenerationStatisticsTime = time.Unix(int64(epoch), 0)

	u32, idx, _ = getU32(data, idx)
	x.electricityGeneratedInCurrentHour = float32(u32) / 100

	u32, idx, _ = getU32(data, idx)
	x.electricityGeneratedInCurrentDay = float32(u32) / 100

	u32, idx, _ = getU32(data, idx)
	x.electricityGeneratedInCurrentMonth = float32(u32) / 100

	u32, _, _ = getU32(data, idx)
	x.electricityGeneratedInCurrentYear = float32(u32) / 100

	return nil
}

func (x *cumulativeData1) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Cumulative Data 1\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Cumulative Generated Electricity          = %3.3f kWh\n", x.cumulativeGeneratedElectricity))
	sb.WriteString(fmt.Sprintf("# Total DC Input Power                      = %3.3f kWh\n", x.totalDCInputPower))
	sb.WriteString(fmt.Sprintf("# Current Electricity Generation Statistics = %s\n", x.currentElectricityGenerationStatisticsTime.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Electricity Generated in Current Hour     = %3.3f kWh\n", x.electricityGeneratedInCurrentHour))
	sb.WriteString(fmt.Sprintf("# Electricity Generated in Current Day      = %3.3f kWh\n", x.electricityGeneratedInCurrentDay))
	sb.WriteString(fmt.Sprintf("# Electricity Generated in Current Month    = %3.3f kWh\n", x.electricityGeneratedInCurrentMonth))
	sb.WriteString(fmt.Sprintf("# Electricity Generated in Current Year     = %3.3f kWh\n", x.electricityGeneratedInCurrentYear))
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No cumulative or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()

		sb.WriteString(fmt.Sprintf("sun2000_cumulative_generate_electricity{model=%q,sn=%q,unit=\"kWh\"} %3.3f\n", id.model, id.sn, x.cumulativeGeneratedElectricity))
		sb.WriteString(fmt.Sprintf("sun2000_total_dc_input_power{model=%q,sn=%q,unit=\"kWh\"} %3.3f\n", id.model, id.sn, x.totalDCInputPower))

		sb.WriteString(fmt.Sprintf("sun2000_current_electricity_generation_statistics_time{model=%q,sn=%q} %d\n", id.model, id.sn, x.currentElectricityGenerationStatisticsTime.Unix()))
		sb.WriteString(fmt.Sprintf("sun2000_electricity_generated_in_current_hour{model=%q,sn=%q,unit=\"kWh\"} %3.3f\n", id.model, id.sn, x.electricityGeneratedInCurrentHour))
		sb.WriteString(fmt.Sprintf("sun2000_electricity_generated_in_current_day{model=%q,sn=%q,unit=\"kWh\"} %3.3f\n", id.model, id.sn, x.electricityGeneratedInCurrentDay))
		sb.WriteString(fmt.Sprintf("sun2000_electricity_generated_in_current_month{model=%q,sn=%q,unit=\"kWh\"} %3.3f\n", id.model, id.sn, x.electricityGeneratedInCurrentMonth))
		sb.WriteString(fmt.Sprintf("sun2000_electricity_generated_in_current_year{model=%q,sn=%q,unit=\"kWh\"} %3.3f\n", id.model, id.sn, x.electricityGeneratedInCurrentYear))
	}
	sb.WriteString("\n")

	return sb.String()
}

type cumulativeData2 struct {
	genericData

	// 32151 U16 1
	numberOfCriticalAlarms uint16
	// 32152 U16 1
	numberOfMajorAlarms uint16
	// 32153 U16 1
	numberOfMinorAlarms uint16
	// 32154 U16 1
	numberOfWarningAlarms uint16
	// 32155 U16 1
	alarmClearanceSerialNumber uint16
	// 32156 Epoch 2
	electricityStatisticsTimeInThePreviousHour time.Time
	// 32158 U32 2 gain 100 kWh
	electricityGeneratedInThePreviousHour float32
	// 32160 Epoch 2
	electricityStatisticsTimeOfThePreviousDay time.Time
	// 32162 U32 2 gain 100 kWh
	electricityGeneratedOnThePreviousDay float32
	// 32164 Epoch 2
	electricityStatisticsTimeOfThePreviousMonth time.Time
	// 32166 U32 2 gain 100 kWh
	electricityGeneratedInPreviousMonth float32
	// 32168 Epoch 2
	electricityStatisticsTimeOfThePreviousYear time.Time
	// 32170 U32 2 gain 100 kWh
	electricityGeneratedInPreviousYear float32
	// 32172 U32 2 gain 100 kWh
	latestActiveAlarmSerialNumber uint32
	// 32174 U32 2
	latestHistoricalAlarmSerialNumber uint32
	// 32176 I16 1 gain 10 V
	totalBusVoltage float32
	// 32177 I16 1 gain 10 V
	maximumPVVoltage float32
	// 32178 I16 1 gain 10 V
	minimumPVVoltage float32
	// 32179 I16 1 gain 10 V
	averagePVNegativeVoltageToGround float32
	// 32180 I16 1 gain 10 V
	maximumPVPositiveVoltageToGround float32
	// 32181 I16 1 gain 10 V
	minimumPVNegativeVoltageToGround float32
	// 32182 U16 1 gain 1 V
	inverterToPEVoltageTolerance inverterToPEVoltageTolerance
	// 32183 Bitfield16 1
	isoFeatureInformation uint16
}

func (x *cumulativeData2) parse(data []byte) (err error) {
	size := 13*2 + 10*4
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var u16 uint16
	var u32 uint32
	var i16 int16
	var epoch uint32
	var idx uint

	x.numberOfCriticalAlarms, idx, _ = getU16(data, idx)
	x.numberOfMajorAlarms, idx, _ = getU16(data, idx)
	x.numberOfMinorAlarms, idx, _ = getU16(data, idx)
	x.numberOfWarningAlarms, idx, _ = getU16(data, idx)
	x.alarmClearanceSerialNumber, idx, _ = getU16(data, idx)

	epoch, idx, _ = getU32(data, idx)
	x.electricityStatisticsTimeInThePreviousHour = time.Unix(int64(epoch), 0)
	u32, idx, _ = getU32(data, idx)
	x.electricityGeneratedInThePreviousHour = float32(u32) / 100

	epoch, idx, _ = getU32(data, idx)
	x.electricityStatisticsTimeOfThePreviousDay = time.Unix(int64(epoch), 0)
	u32, idx, _ = getU32(data, idx)
	x.electricityGeneratedOnThePreviousDay = float32(u32) / 100

	epoch, idx, _ = getU32(data, idx)
	x.electricityStatisticsTimeOfThePreviousMonth = time.Unix(int64(epoch), 0)
	u32, idx, _ = getU32(data, idx)
	x.electricityGeneratedInPreviousMonth = float32(u32) / 100

	epoch, idx, _ = getU32(data, idx)
	x.electricityStatisticsTimeOfThePreviousYear = time.Unix(int64(epoch), 0)
	u32, idx, _ = getU32(data, idx)
	x.electricityGeneratedInPreviousYear = float32(u32) / 100

	x.latestActiveAlarmSerialNumber, idx, _ = getU32(data, idx)
	x.latestHistoricalAlarmSerialNumber, idx, _ = getU32(data, idx)

	i16, idx, _ = getI16(data, idx)
	x.totalBusVoltage = float32(i16) / 10
	i16, idx, _ = getI16(data, idx)
	x.maximumPVVoltage = float32(i16) / 10
	i16, idx, _ = getI16(data, idx)
	x.minimumPVVoltage = float32(i16) / 10
	i16, idx, _ = getI16(data, idx)
	x.averagePVNegativeVoltageToGround = float32(i16) / 10
	i16, idx, _ = getI16(data, idx)
	x.maximumPVPositiveVoltageToGround = float32(i16) / 10
	i16, idx, _ = getI16(data, idx)
	x.minimumPVNegativeVoltageToGround = float32(i16) / 10

	u16, idx, _ = getU16(data, idx)
	x.inverterToPEVoltageTolerance = inverterToPEVoltageTolerance(u16)

	u16, _, _ = getU16(data, idx)
	x.isoFeatureInformation = u16

	return nil
}

func (x *cumulativeData2) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Cumulative Data 2\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Number of Critical Alarms             = %d\n", x.numberOfCriticalAlarms))
	sb.WriteString(fmt.Sprintf("# Number of Major Alarms                = %d\n", x.numberOfMajorAlarms))
	sb.WriteString(fmt.Sprintf("# Number of Minor Alarms                = %d\n", x.numberOfMinorAlarms))
	sb.WriteString(fmt.Sprintf("# Number of Warning Alarms              = %d\n", x.numberOfWarningAlarms))
	sb.WriteString(fmt.Sprintf("# Alarm Clearance Serial Number         = %d\n", x.alarmClearanceSerialNumber))
	sb.WriteString("#\n")
	sb.WriteString(fmt.Sprintf("# Electricity Statistics Time in the Previous Hour  = %s\n", x.electricityStatisticsTimeInThePreviousHour.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Electricity Generated in the Previous Hour        = %3.3f kWh\n", x.electricityGeneratedInThePreviousHour))
	sb.WriteString(fmt.Sprintf("# Electricity Statistics Time of the Previous Day   = %s\n", x.electricityStatisticsTimeOfThePreviousDay.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Electricity Generated on the Previous Day         = %3.3f kWh\n", x.electricityGeneratedOnThePreviousDay))
	sb.WriteString(fmt.Sprintf("# Electricity Statistics Time of the Previous Month = %s\n", x.electricityStatisticsTimeOfThePreviousMonth.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Electricity Generated in the Previous Month       = %3.3f kWh\n", x.electricityGeneratedInPreviousMonth))
	sb.WriteString(fmt.Sprintf("# Electricity Statistics Time of the Previous Year  = %s\n", x.electricityStatisticsTimeOfThePreviousYear.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Electricity Generated in the Previous Year        = %3.3f kWh\n", x.electricityGeneratedInPreviousYear))
	sb.WriteString("#\n")
	sb.WriteString(fmt.Sprintf("# Latest Active Alarm Serial Number     = %d\n", x.latestActiveAlarmSerialNumber))
	sb.WriteString(fmt.Sprintf("# Latest Historical Alarm Serial Number = %d\n", x.latestHistoricalAlarmSerialNumber))
	sb.WriteString("#\n")
	sb.WriteString(fmt.Sprintf("# Total Bus Voltage                     = %3.1f V\n", x.totalBusVoltage))
	sb.WriteString(fmt.Sprintf("# Maximum PV Voltage                    = %3.1f V\n", x.maximumPVVoltage))
	sb.WriteString(fmt.Sprintf("# Minimum PV Voltage                    = %3.1f V\n", x.minimumPVVoltage))
	sb.WriteString(fmt.Sprintf("# Average PV Negative Voltage to Ground = %3.1f V\n", x.averagePVNegativeVoltageToGround))
	sb.WriteString(fmt.Sprintf("# Maximum PV Positive Voltage to Ground = %3.1f V\n", x.maximumPVPositiveVoltageToGround))
	sb.WriteString(fmt.Sprintf("# Minimum PV Negative Voltage to Ground = %3.1f V\n", x.minimumPVNegativeVoltageToGround))
	sb.WriteString(fmt.Sprintf("# Inverter to PE Voltage Tolerance      = %d\t%s\n", x.inverterToPEVoltageTolerance, x.inverterToPEVoltageTolerance))
	sb.WriteString(fmt.Sprintf("# ISO Feature Information               = %d\n", x.isoFeatureInformation))
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No cumulative or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()

		sb.WriteString(fmt.Sprintf("sun2000_number_of_alarms{model=%q,sn=%q,level=\"Critical\"} %d\n", id.model, id.sn, x.numberOfCriticalAlarms))
		sb.WriteString(fmt.Sprintf("sun2000_number_of_alarms{model=%q,sn=%q,level=\"Major\"} %d\n", id.model, id.sn, x.numberOfMajorAlarms))
		sb.WriteString(fmt.Sprintf("sun2000_number_of_alarms{model=%q,sn=%q,level=\"Minor\"} %d\n", id.model, id.sn, x.numberOfMinorAlarms))
		sb.WriteString(fmt.Sprintf("sun2000_number_of_alarms{model=%q,sn=%q,level=\"Warning\"} %d\n", id.model, id.sn, x.numberOfWarningAlarms))
		sb.WriteString(fmt.Sprintf("sun2000_alarm_clearance_serial_number{model=%q,sn=%q} %d\n", id.model, id.sn, x.alarmClearanceSerialNumber))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("sun2000_electricity_statistics_time_in_the_previous_hour{model=%q,sn=%q} %d\n", id.model, id.sn, x.electricityStatisticsTimeInThePreviousHour.Unix()))
		sb.WriteString(fmt.Sprintf("sun2000_electricity_generated_in_the_previous_hour{model=%q,sn=%q,unit=\"kWh\"} %3.3f\n", id.model, id.sn, x.electricityGeneratedInThePreviousHour))
		sb.WriteString(fmt.Sprintf("sun2000_electricity_statistics_time_of_the_previous_day{model=%q,sn=%q} %d\n", id.model, id.sn, x.electricityStatisticsTimeOfThePreviousDay.Unix()))
		sb.WriteString(fmt.Sprintf("sun2000_electricity_generated_on_the_previous_day{model=%q,sn=%q,unit=\"kWh\"} %3.3f\n", id.model, id.sn, x.electricityGeneratedOnThePreviousDay))
		sb.WriteString(fmt.Sprintf("sun2000_electricity_statistics_time_of_the_previous_month{model=%q,sn=%q} %d\n", id.model, id.sn, x.electricityStatisticsTimeOfThePreviousMonth.Unix()))
		sb.WriteString(fmt.Sprintf("sun2000_electricity_generated_in_the_previous_month{model=%q,sn=%q,unit=\"kWh\"} %3.3f\n", id.model, id.sn, x.electricityGeneratedInPreviousMonth))
		sb.WriteString(fmt.Sprintf("sun2000_electricity_statistics_time_of_the_previous_year{model=%q,sn=%q} %d\n", id.model, id.sn, x.electricityStatisticsTimeOfThePreviousYear.Unix()))
		sb.WriteString(fmt.Sprintf("sun2000_electricity_generated_in_the_previous_year{model=%q,sn=%q,unit=\"kWh\"} %3.3f\n", id.model, id.sn, x.electricityGeneratedInPreviousYear))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("sun2000_latest_active_alarm_serial_number{model=%q,sn=%q} %d\n", id.model, id.sn, x.latestActiveAlarmSerialNumber))
		sb.WriteString(fmt.Sprintf("sun2000_latest_historical_alarm_serial_number{model=%q,sn=%q} %d\n", id.model, id.sn, x.latestHistoricalAlarmSerialNumber))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("sun2000_total_bus_voltage{model=%q,sn=%q,unit=\"V\"} %3.1f\n", id.model, id.sn, x.totalBusVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_maximum_pv_voltage{model=%q,sn=%q,unit=\"V\"} %3.1f\n", id.model, id.sn, x.maximumPVVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_minimum_pv_voltage{model=%q,sn=%q,unit=\"V\"} %3.1f\n", id.model, id.sn, x.minimumPVVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_average_pv_negative_voltage_to_ground{model=%q,sn=%q,unit=\"V\"} %3.1f\n", id.model, id.sn, x.averagePVNegativeVoltageToGround))
		sb.WriteString(fmt.Sprintf("sun2000_maximum_pv_positive_voltage_to_ground{model=%q,sn=%q,unit=\"V\"} %3.1f\n", id.model, id.sn, x.maximumPVPositiveVoltageToGround))
		sb.WriteString(fmt.Sprintf("sun2000_minimum_pv_negative_voltage_to_ground{model=%q,sn=%q,unit=\"V\"} %3.1f\n", id.model, id.sn, x.minimumPVNegativeVoltageToGround))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("sun2000_inverter_to_pe_voltage_tolerance{model=%q,sn=%q,unit=\"V\"} %d\n", id.model, id.sn, x.inverterToPEVoltageTolerance))
		sb.WriteString(fmt.Sprintf("sun2000_iso_feature_information{model=%q,sn=%q} %d\n", id.model, id.sn, x.isoFeatureInformation))
	}

	sb.WriteString("\n")

	return sb.String()
}

type inverterToPEVoltageTolerance uint16

func (x inverterToPEVoltageTolerance) String() string {
	switch x {
	case 0:
		return "1000V/1100V inverter"
	case 1500:
		return "HAV1"
	case 1502:
		return "HAV2"
	default:
		return fmt.Sprintf("Unknown tolerance %d", x)
	}
}

type cumulativeData3 struct {
	genericData

	// 32190 E16 1
	builtInPIDRunningStatus uint16
	// 32191 I16 1 gain 10 V
	pvNegativeVoltageToGround float32
}

func (x *cumulativeData3) parse(data []byte) (err error) {
	size := 2 * 2
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var i16 int16
	var idx uint

	x.builtInPIDRunningStatus, idx, _ = getU16(data, idx)

	i16, _, _ = getI16(data, idx)
	x.pvNegativeVoltageToGround = float32(i16) / 10

	return nil
}

func (x *cumulativeData3) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Cumulative Data 3\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Built-in PID Running Status           = %d\n", x.builtInPIDRunningStatus))
	sb.WriteString(fmt.Sprintf("# PV Negative Voltage to Ground         = %3.1f V\n", x.pvNegativeVoltageToGround))
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No cumulative or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()

		sb.WriteString(fmt.Sprintf("sun2000_built_in_pid_running_status{model=%q,sn=%q} %d\n", id.model, id.sn, x.builtInPIDRunningStatus))
		sb.WriteString(fmt.Sprintf("sun2000_pv_negative_voltage_to_ground{model=%q,sn=%q,unit=\"V\"} %3.1f\n", id.model, id.sn, x.pvNegativeVoltageToGround))
	}
	sb.WriteString("\n")

	return sb.String()
}

type mpptData1 struct {
	genericData

	// 32212-32230 U32 2x10 gain 100 kWh
	cumulativeDCEnergyYieldOfMPPT [10]float32
}

func (x *mpptData1) parse(data []byte) (err error) {
	size := 10 * 2
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var u32 uint32
	var idx uint

	for i := range x.cumulativeDCEnergyYieldOfMPPT {
		u32, idx, _ = getU32(data, idx)
		x.cumulativeDCEnergyYieldOfMPPT[i] = float32(u32) / 100
	}

	return nil
}

func (x *mpptData1) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# MPPT Data\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	for i, y := range x.cumulativeDCEnergyYieldOfMPPT {
		if y != 0 || i < int(id.numberOfMPPTs) {
			sb.WriteString(fmt.Sprintf("# Cumulative DC Energy Yield of MPPT %2d = %3.3f kWh\n", i+1, y))
		}
	}
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No MPPT or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()

		for i, y := range x.cumulativeDCEnergyYieldOfMPPT {
			if y != 0 || i < int(id.numberOfMPPTs) {
				sb.WriteString(fmt.Sprintf("sun2000_cumulative_dc_energy_yield_of_mppt{model=%q,sn=%q,mppt=\"%d\",unit=\"kWh\"} %3.3f\n", id.model, id.sn, i+1, y))
			}
		}
	}
	sb.WriteString("\n")

	return sb.String()
}

type alarmData2 struct {
	genericData

	// 32252-32254 Bitfield16 1x3
	// 32271-32272 Bitfield16 1x2
	monitoringAlarm [5]uint16
	// 32255-32270 Bitfield16 1x16
	// 32273-32274 Bitfield16 1x2
	externalPowerAlarm [18]uint16
}

func (x *alarmData2) parse(data []byte) (err error) {
	size := 5*2 + 18*2
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var idx uint

	for i := 0; i < 3; i++ {
		x.monitoringAlarm[i], idx, _ = getU16(data, idx)
	}

	for i := 0; i < 16; i++ {
		x.externalPowerAlarm[i], idx, _ = getU16(data, idx)
	}

	for i := 3; i < 5; i++ {
		x.monitoringAlarm[i], idx, _ = getU16(data, idx)
	}

	for i := 16; i < 18; i++ {
		x.externalPowerAlarm[i], idx, _ = getU16(data, idx)
	}

	return nil
}

func (x *alarmData2) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Alarm Data 2\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	for i, y := range x.monitoringAlarm {
		sb.WriteString(fmt.Sprintf("# Monitoring Alarm %2d = %#04x\t%#016b\n", i+1, y, y))
	}
	sb.WriteString("\n")
	for i, y := range x.externalPowerAlarm {
		sb.WriteString(fmt.Sprintf("# External Power Alarm %2d = %#04x\t%#016b\n", i+1, y, y))
	}
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No alarm or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()

		for i, y := range x.monitoringAlarm {
			sb.WriteString(fmt.Sprintf("sun2000_monitoring_alarm{model=%q,sn=%q,alarm=\"%d\"} %d\n", id.model, id.sn, i+1, y))
		}
		sb.WriteString("\n")
		for i, y := range x.externalPowerAlarm {
			sb.WriteString(fmt.Sprintf("sun2000_external_power_alarm{model=%q,sn=%q,alarm=\"%d\"} %d\n", id.model, id.sn, i+1, y))
		}
	}
	sb.WriteString("\n")

	return sb.String()
}

type stringAccessData struct {
	genericData

	// 32300-32317 E16 1x18
	stringAccessStatus [18]uint16
}

func (x *stringAccessData) parse(data []byte) (err error) {
	size := 18 * 2
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var idx uint

	for i := range x.stringAccessStatus {
		x.stringAccessStatus[i], idx, _ = getU16(data, idx)
	}

	return nil
}

func (x *stringAccessData) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# String Access Data\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	for i, y := range x.stringAccessStatus {
		sb.WriteString(fmt.Sprintf("# String Access Status %2d = %#04x\t%#016b\n", i+1, y, y))
	}
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No string access or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()

		for i, y := range x.stringAccessStatus {
			sb.WriteString(fmt.Sprintf("sun2000_string_access_status{model=%q,sn=%q,string=\"%d\"} %d\n", id.model, id.sn, i+1, y))
		}
	}
	sb.WriteString("\n")

	return sb.String()
}

type mpptData2 struct {
	genericData

	// 32324-32342 U32 2x10 gain 1000 kW
	mpptTotalInputPower [10]float32
}

func (x *mpptData2) parse(data []byte) (err error) {
	size := 10 * 2
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var u32 uint32
	var idx uint

	for i := range x.mpptTotalInputPower {
		u32, idx, _ = getU32(data, idx)
		x.mpptTotalInputPower[i] = float32(u32) / 1000
	}

	return nil
}

func (x *mpptData2) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# MPPT Data 2\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	for i, y := range x.mpptTotalInputPower {
		if y != 0 || i < int(id.numberOfMPPTs) {
			sb.WriteString(fmt.Sprintf("# MPPT %2d Total Input Power = %3.3f kW\n", i+1, y))
		}
	}
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No MPPT or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()

		for i, y := range x.mpptTotalInputPower {
			if y != 0 || i < int(id.numberOfMPPTs) {
				sb.WriteString(fmt.Sprintf("sun2000_mppt_total_input_power{model=%q,sn=%q,mppt=\"%d\",unit=\"kW\"} %3.3f\n", id.model, id.sn, i+1, y))
			}
		}
	}
	sb.WriteString("\n")

	return sb.String()
}

type internalTemperatureData struct {
	genericData

	// 35021-35032 I16 1x12 gain 10 ℃
	internalTemperature [12]float32
}

func (x *internalTemperatureData) parse(data []byte) (err error) {
	size := 12 * 2
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var i16 int16
	var idx uint

	for i := range x.internalTemperature {
		i16, idx, _ = getI16(data, idx)
		x.internalTemperature[i] = float32(i16) / 10
	}

	return nil
}

func (x *internalTemperatureData) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Internal Temperature Data\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	for i, y := range x.internalTemperature {
		if y != 0 {
			label := internalTemperatureLabel(i + 1)
			sb.WriteString(fmt.Sprintf("# Internal Temperature %2d = %3.1f ℃ %s\n", i+1, y, label))
		}
	}
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No internal temperature or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()

		for i, y := range x.internalTemperature {
			if y != 0 {
				label := internalTemperatureLabel(i + 1)
				sb.WriteString(fmt.Sprintf("sun2000_internal_temperature{model=%q,sn=%q,sensor=\"%d %s\",unit=\"℃\"} %3.1f\n", id.model, id.sn, i+1, label, y))
			}
		}
	}
	sb.WriteString("\n")

	return sb.String()
}

type internalTemperatureLabel int

func (x internalTemperatureLabel) String() string {
	switch x {
	case 1:
		return "Inverter Module A"
	case 2:
		return "Inverter Module B"
	case 3:
		return "Inverter Module C"
	case 4:
		return "Anti-reverse module 1"
	case 5:
		return "Output board relay / ambient temperature max"
	case 6:
		return "Output board / power board input / power board inverter max"
	case 7:
		return "Anti-reverse module 2"
	case 8:
		return "DC terminal 1/2 max"
	case 9:
		return "AC terminal 1/2/3 max"
	case 10:
		return ""
	case 11:
		return ""
	case 12:
		return ""
	default:
		return ""
	}
}

type meterData struct {
	genericData

	// 37100 U16 1
	meterStatus uint16
	// 37101 I32 2 gain 10 V
	gridPhaseAVoltage float32
	// 37103 I32 2 gain 10 V
	gridPhaseBVoltage float32
	// 37105 I32 2 gain 10 V
	gridPhaseCVoltage float32
	// 37107 I32 2 gain 100 A
	gridPhaseACurrent float32
	// 37109 I32 2 gain 100 A
	gridPhaseBCurrent float32
	// 37111 I32 2 gain 100 A
	gridPhaseCCurrent float32

	// 37113 I32 2 gain 1000 kW
	gridActivePower float32 // >0 feed-in to the grid, <0 supply from the grid
	// 37115 I32 2 gain 1 Var
	gridReactivePower float32
	// 37117 I16 1 gain 1000
	gridPowerFactor float32
	// 37118 I16 1 gain 100 Hz
	gridFrequency float32
	// 37119 I32 2 gain 100 kWh
	gridPositiveActiveElectricity float32 // Electricity fed by the inverter to the power grid
	// 37121 I32 2 gain 100 kWh
	gridReverseActivePower float32 // Power supplied to a distributed system from the power grid
	// 37123 I32 2 gain 100 kVar h
	gridAccumulatedReactivePower float32

	// 37125 U16 1
	meterType uint16 // 0: single-phase; 1: three-phase
	// 37126 I32 2 gain 10 V
	gridLineABVoltage float32
	// 37128 I32 2 gain 10 V
	gridLineBCVoltage float32
	// 37130 I32 2 gain 10 V
	gridLineCAVoltage float32
	// 37132 I32 2 gain 1000 kW
	gridPhaseAActivePower float32
	// 37134 I32 2 gain 1000 kW
	gridPhaseBActivePower float32
	// 37136 I32 2 gain 1000 kW
	gridPhaseCActivePower float32
	// 37138 U16 1
	meterModelDetectionResult uint16 // 0: being identified; 1: The selected model is the same as the actual model of the connected meter; 2: The selected model is different from the actual model of the connected meter
}

func (x *meterData) parse(data []byte) (err error) {
	size := 1 + 6*2
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var u16 uint16
	var i16 int16
	var i32 int32
	var idx uint

	u16, idx, _ = getU16(data, idx)
	x.meterStatus = u16

	i32, idx, _ = getI32(data, idx)
	x.gridPhaseAVoltage = float32(i32) / 10
	i32, idx, _ = getI32(data, idx)
	x.gridPhaseBVoltage = float32(i32) / 10
	i32, idx, _ = getI32(data, idx)
	x.gridPhaseCVoltage = float32(i32) / 10

	i32, idx, _ = getI32(data, idx)
	x.gridPhaseACurrent = float32(i32) / 100
	i32, idx, _ = getI32(data, idx)
	x.gridPhaseBCurrent = float32(i32) / 100
	i32, idx, _ = getI32(data, idx)
	x.gridPhaseCCurrent = float32(i32) / 100

	i32, idx, _ = getI32(data, idx)
	x.gridActivePower = float32(i32) / 1000
	i32, idx, _ = getI32(data, idx)
	x.gridReactivePower = float32(i32)
	i16, idx, _ = getI16(data, idx)
	x.gridPowerFactor = float32(i16) / 100

	i16, idx, _ = getI16(data, idx)
	x.gridFrequency = float32(i16) / 100
	i32, idx, _ = getI32(data, idx)
	x.gridPositiveActiveElectricity = float32(i32) / 100
	i32, idx, _ = getI32(data, idx)
	x.gridReverseActivePower = float32(i32) / 100
	i32, idx, _ = getI32(data, idx)
	x.gridAccumulatedReactivePower = float32(i32) / 100

	x.meterType, idx, _ = getU16(data, idx)

	i32, idx, _ = getI32(data, idx)
	x.gridLineABVoltage = float32(i32) / 10
	i32, idx, _ = getI32(data, idx)
	x.gridLineBCVoltage = float32(i32) / 10
	i32, idx, _ = getI32(data, idx)
	x.gridLineCAVoltage = float32(i32) / 10

	i32, idx, _ = getI32(data, idx)
	x.gridPhaseAActivePower = float32(i32) / 1000
	i32, idx, _ = getI32(data, idx)
	x.gridPhaseBActivePower = float32(i32) / 1000
	i32, idx, _ = getI32(data, idx)
	x.gridPhaseCActivePower = float32(i32) / 1000

	x.meterModelDetectionResult, _, _ = getU16(data, idx)

	return nil
}

func (x *meterData) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Meter Data\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")
	var txt, txt2, txt3 string
	switch x.meterStatus {
	case 0:
		txt = "offline"
	case 1:
		txt = "online"
	}
	sb.WriteString(fmt.Sprintf("# Meter Status       = %#04x\t%#016b\t%s\n", x.meterStatus, x.meterStatus, txt))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Grid Phase A Voltage = %3.1f V\n", x.gridPhaseAVoltage))
	sb.WriteString(fmt.Sprintf("# Grid Phase B Voltage = %3.1f V\n", x.gridPhaseBVoltage))
	sb.WriteString(fmt.Sprintf("# Grid Phase C Voltage = %3.1f V\n", x.gridPhaseCVoltage))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Grid Phase A Current = %3.2f A\n", x.gridPhaseACurrent))
	sb.WriteString(fmt.Sprintf("# Grid Phase B Current = %3.2f A\n", x.gridPhaseBCurrent))
	sb.WriteString(fmt.Sprintf("# Grid Phase C Current = %3.2f A\n", x.gridPhaseCCurrent))
	sb.WriteString("\n")
	pA := x.gridPhaseAVoltage * x.gridPhaseACurrent
	sb.WriteString(fmt.Sprintf("# Grid Phase A Power   = %5.2f VA\n", pA))
	pB := x.gridPhaseBVoltage * x.gridPhaseBCurrent
	sb.WriteString(fmt.Sprintf("# Grid Phase B Power   = %5.2f VA\n", pB))
	pC := x.gridPhaseCVoltage * x.gridPhaseCCurrent
	sb.WriteString(fmt.Sprintf("# Grid Phase C Power   = %5.2f VA\n", pC))
	pTotal := pA + pB + pC
	sb.WriteString("# -----------------------------\n")
	sb.WriteString(fmt.Sprintf("# Grid Total Power     = %5.2f VA\n", pTotal))
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("# Grid Active Power    = %6.3f kW\n", x.gridActivePower))
	sb.WriteString(fmt.Sprintf("# Grid Reactive Power  = %4.0f Var\n", x.gridReactivePower))
	sb.WriteString(fmt.Sprintf("# Grid Power Factor    = %3.2f\n", x.gridPowerFactor))
	sb.WriteString(fmt.Sprintf("# Grid Frequency       = %3.2f Hz\n", x.gridFrequency))
	sb.WriteString(fmt.Sprintf("# Grid Positive Active Electricity  = %6.2f kWh\n", x.gridPositiveActiveElectricity))
	sb.WriteString(fmt.Sprintf("# Grid Reverse Active Power         = %6.2f kWh\n", x.gridReverseActivePower))
	sb.WriteString(fmt.Sprintf("# Grid Accumulated Reactive Power   = %6.2f kVar h\n", x.gridAccumulatedReactivePower))
	sb.WriteString("\n")
	txt2 = ""
	switch x.meterType {
	case 0:
		txt2 = "single-phase"
	case 1:
		txt2 = "three-phase"
	}
	sb.WriteString(fmt.Sprintf("# Meter Type          = %d\t%s\n", x.meterType, txt2))

	sb.WriteString(fmt.Sprintf("# Grid Line AB Voltage = %3.1f V\n", x.gridLineABVoltage))
	sb.WriteString(fmt.Sprintf("# Grid Line BC Voltage = %3.1f V\n", x.gridLineBCVoltage))
	sb.WriteString(fmt.Sprintf("# Grid Line CA Voltage = %3.1f V\n", x.gridLineCAVoltage))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Grid Phase A Active Power = %3.3f kW\n", x.gridPhaseAActivePower))
	sb.WriteString(fmt.Sprintf("# Grid Phase B Active Power = %3.3f kW\n", x.gridPhaseBActivePower))
	sb.WriteString(fmt.Sprintf("# Grid Phase C Active Power = %3.3f kW\n", x.gridPhaseCActivePower))
	sb.WriteString("# -----------------------------\n")
	pTotalActive := x.gridPhaseAActivePower + x.gridPhaseBActivePower + x.gridPhaseCActivePower
	sb.WriteString(fmt.Sprintf("# Grid Total Active Power   = %3.3f kW\n", pTotalActive))
	sb.WriteString("\n")

	txt3 = ""
	switch x.meterModelDetectionResult {
	case 0:
		txt3 = "being identified"
	case 1:
		txt3 = "The selected model is the same as the actual model of the connected meter"
	case 2:
		txt3 = "The selected model is different from the actual model of the connected meter"
	}
	sb.WriteString(fmt.Sprintf("# Meter Model Detection Result = %d\t%s\n", x.meterModelDetectionResult, txt3))

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No meter or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()

		sb.WriteString(fmt.Sprintf("sun2000_meter_status{model=%q,sn=%q} %d\n", id.model, id.sn, x.meterStatus))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("sun2000_grid_phase_voltage{model=%q,sn=%q,phase=\"A\",unit=\"V\"} %3.1f\n", id.model, id.sn, x.gridPhaseAVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_grid_phase_voltage{model=%q,sn=%q,phase=\"B\",unit=\"V\"} %3.1f\n", id.model, id.sn, x.gridPhaseBVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_grid_phase_voltage{model=%q,sn=%q,phase=\"C\",unit=\"V\"} %3.1f\n", id.model, id.sn, x.gridPhaseCVoltage))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("sun2000_grid_phase_current{model=%q,sn=%q,phase=\"A\",unit=\"A\"} %3.2f\n", id.model, id.sn, x.gridPhaseACurrent))
		sb.WriteString(fmt.Sprintf("sun2000_grid_phase_current{model=%q,sn=%q,phase=\"B\",unit=\"A\"} %3.2f\n", id.model, id.sn, x.gridPhaseBCurrent))
		sb.WriteString(fmt.Sprintf("sun2000_grid_phase_current{model=%q,sn=%q,phase=\"C\",unit=\"A\"} %3.2f\n", id.model, id.sn, x.gridPhaseCCurrent))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("sun2000_grid_phase_power{model=%q,sn=%q,phase=\"A\",unit=\"VA\"} %3.2f\n", id.model, id.sn, pA))
		sb.WriteString(fmt.Sprintf("sun2000_grid_phase_power{model=%q,sn=%q,phase=\"B\",unit=\"VA\"} %3.2f\n", id.model, id.sn, pB))
		sb.WriteString(fmt.Sprintf("sun2000_grid_phase_power{model=%q,sn=%q,phase=\"C\",unit=\"VA\"} %3.2f\n", id.model, id.sn, pC))
		sb.WriteString(fmt.Sprintf("sun2000_grid_total_power{model=%q,sn=%q,unit=\"VA\"} %3.2f\n", id.model, id.sn, pTotal))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("sun2000_grid_active_power{model=%q,sn=%q,unit=\"kW\"} %3.3f\n", id.model, id.sn, x.gridActivePower))
		sb.WriteString(fmt.Sprintf("sun2000_grid_reactive_power{model=%q,sn=%q,unit=\"Var\"} %3.0f\n", id.model, id.sn, x.gridReactivePower))
		sb.WriteString(fmt.Sprintf("sun2000_grid_power_factor{model=%q,sn=%q} %3.2f\n", id.model, id.sn, x.gridPowerFactor))
		sb.WriteString(fmt.Sprintf("sun2000_grid_frequency{model=%q,sn=%q,unit=\"Hz\"} %3.2f\n", id.model, id.sn, x.gridFrequency))
		sb.WriteString(fmt.Sprintf("sun2000_grid_positive_active_electricity{model=%q,sn=%q,unit=\"kWh\"} %3.2f\n", id.model, id.sn, x.gridPositiveActiveElectricity))
		sb.WriteString(fmt.Sprintf("sun2000_grid_reverse_active_power{model=%q,sn=%q,unit=\"kWh\"} %3.2f\n", id.model, id.sn, x.gridReverseActivePower))
		sb.WriteString(fmt.Sprintf("sun2000_grid_accumulated_reactive_power{model=%q,sn=%q,unit=\"kVar h\"} %3.2f\n", id.model, id.sn, x.gridAccumulatedReactivePower))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("sun2000_meter_type{model=%q,sn=%q} %d\n", id.model, id.sn, x.meterType))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("sun2000_grid_line_voltage{model=%q,sn=%q,line=\"AB\",unit=\"V\"} %3.1f\n", id.model, id.sn, x.gridLineABVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_grid_line_voltage{model=%q,sn=%q,line=\"BC\",unit=\"V\"} %3.1f\n", id.model, id.sn, x.gridLineBCVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_grid_line_voltage{model=%q,sn=%q,line=\"CA\",unit=\"V\"} %3.1f\n", id.model, id.sn, x.gridLineCAVoltage))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("sun2000_grid_phase_active_power{model=%q,sn=%q,phase=\"A\",unit=\"kW\"} %3.3f\n", id.model, id.sn, x.gridPhaseAActivePower))
		sb.WriteString(fmt.Sprintf("sun2000_grid_phase_active_power{model=%q,sn=%q,phase=\"B\",unit=\"kW\"} %3.3f\n", id.model, id.sn, x.gridPhaseBActivePower))
		sb.WriteString(fmt.Sprintf("sun2000_grid_phase_active_power{model=%q,sn=%q,phase=\"C\",unit=\"kW\"} %3.3f\n", id.model, id.sn, x.gridPhaseCActivePower))
		sb.WriteString(fmt.Sprintf("sun2000_grid_total_active_power{model=%q,sn=%q,unit=\"kW\"} %3.3f\n", id.model, id.sn, pTotalActive))

		sb.WriteString(fmt.Sprintf("sun2000_meter_model_detection_result{model=%q,sn=%q} %d\n", id.model, id.sn, x.meterModelDetectionResult))

	}
	sb.WriteString("\n")

	return sb.String()
}

type esu1Data struct {
	genericData

	// 37000 U16 1
	runningStatus esuRunningStatus
	// 37001 I32 2 gain 1000 kW
	chargeAndDischargePower float32 // > 0: charging < 0: discharging
	// 37003 uint16 1 gain 10 V
	busVoltage float32
	// 37004 uint16 1 gain 10 %
	batterySOC float32
	// 37006 uint16 1
	workingMode esuWorkingMode
	// 37007 uint32 2 gain 1 W
	ratedChargePower uint32
	// 37009 uint32 2 gain 1 W
	ratedDischargePower uint32
	// 37014 uint16 1
	faultID uint16
	// 37015 uint32 2 gain 100 kWh
	currentDayChargeCapacity float32
	// 37017 uint32 2 gain 100 kWh
	currentDayDischargeCapacity float32
	// 37021 int16 1 gain 10 A
	busCurrent float32
	// 37022 int16 1
	batteryTemperature float32
	// 37025 uint16 1 mins
	remainingChargeDischargeTime uint16
	// 37026 string 10
	dcdcVersion string
	// 37036 string 10
	bmsVersion string
	// 37046 uint32 2 gain 1 W
	maximumChargePower uint32
	// 37048 uint32 2 gain 1 W
	maximumDischargePower uint32
	// 37052 string 10
	sn string
	// 37066 uint32 2 gain 100 kWh
	totalCharge float32
	// 37068 uint32 2 gain 100 kWh
	totalDischarge float32

	pack [3]batteryData
}

func (x *esu1Data) parse(data []byte) (err error) {
	size := 70
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var u16 uint16
	var i16 int16
	var i32 int32
	var u32 uint32
	var idx uint

	u16, idx, _ = getU16(data, idx)
	x.runningStatus = esuRunningStatus(u16)

	i32, idx, _ = getI32(data, idx)
	x.chargeAndDischargePower = float32(i32) / 1000

	u16, idx, _ = getU16(data, idx)
	x.busVoltage = float32(u16) / 10
	u16, idx, _ = getU16(data, idx)
	x.batterySOC = float32(u16) / 10

	idx, _ = skipRecords(data, idx, 1)

	u16, idx, _ = getU16(data, idx)
	x.workingMode = esuWorkingMode(u16)

	u32, idx, _ = getU32(data, idx)
	x.ratedChargePower = u32
	u32, idx, _ = getU32(data, idx)
	x.ratedDischargePower = u32

	idx, _ = skipRecords(data, idx, 3)

	x.faultID, idx, _ = getU16(data, idx)

	u32, idx, _ = getU32(data, idx)
	x.currentDayChargeCapacity = float32(u32) / 100
	u32, idx, _ = getU32(data, idx)
	x.currentDayDischargeCapacity = float32(u32) / 100

	idx, _ = skipRecords(data, idx, 2)

	i16, idx, _ = getI16(data, idx)
	x.busCurrent = float32(i16) / 10
	i16, idx, _ = getI16(data, idx)
	x.batteryTemperature = float32(i16) / 10

	idx, _ = skipRecords(data, idx, 2)

	u16, idx, _ = getU16(data, idx)
	x.remainingChargeDischargeTime = u16

	x.dcdcVersion, idx, _ = getSTR(data, idx, 10)
	x.bmsVersion, idx, _ = getSTR(data, idx, 10)

	u32, idx, _ = getU32(data, idx)
	x.maximumChargePower = u32
	u32, idx, _ = getU32(data, idx)
	x.maximumDischargePower = u32

	idx, _ = skipRecords(data, idx, 2)

	x.sn, idx, _ = getSTR(data, idx, 10)

	idx, _ = skipRecords(data, idx, 4)

	u32, idx, _ = getU32(data, idx)
	x.totalCharge = float32(u32) / 100
	u32, _, _ = getU32(data, idx)
	x.totalDischarge = float32(u32) / 100

	return nil
}

func (x *esu1Data) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# ESU1 Data\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("# SN                         = %q\n", x.sn))
	sb.WriteString(fmt.Sprintf("# Running Status             = %d\t%s\n", x.runningStatus, x.runningStatus))
	var txt string
	if x.chargeAndDischargePower > 0 {
		txt = "charging"
	} else if x.chargeAndDischargePower < 0 {
		txt = "discharging"
	}

	sb.WriteString(fmt.Sprintf("# Charge And Discharge Power = %6.3f kW %s\n", x.chargeAndDischargePower, txt))
	sb.WriteString(fmt.Sprintf("# Bus Voltage                = %3.1f V\n", x.busVoltage))
	sb.WriteString(fmt.Sprintf("# Bus Current                = %3.1f A\n", x.busCurrent))
	sb.WriteString(fmt.Sprintf("# Battery SOC                = %3.1f %%\n", x.batterySOC))
	sb.WriteString(fmt.Sprintf("# Remaining Charge Discharge Time = %d mins\n", x.remainingChargeDischargeTime))
	sb.WriteString(fmt.Sprintf("# Battery Temperature        = %3.1f ℃\n", x.batteryTemperature))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Working Mode              = %d\t%s\n", x.workingMode, x.workingMode))
	sb.WriteString(fmt.Sprintf("# Rated Charge Power        = %d W\n", x.ratedChargePower))
	sb.WriteString(fmt.Sprintf("# Rated Discharge Power     = %d W\n", x.ratedDischargePower))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Fault ID                  = %d\n", x.faultID))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# DCDC Version              = %q\n", x.dcdcVersion))
	sb.WriteString(fmt.Sprintf("# BMS Version               = %q\n", x.bmsVersion))
	sb.WriteString(fmt.Sprintf("# Maximum Charge Power      = %d W\n", x.maximumChargePower))
	sb.WriteString(fmt.Sprintf("# Maximum Discharge Power   = %d W\n", x.maximumDischargePower))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Current Day Charge Capacity     = %6.2f kWh\n", x.currentDayChargeCapacity))
	sb.WriteString(fmt.Sprintf("# Current Day Discharge Capacity  = %6.2f kWh\n", x.currentDayDischargeCapacity))
	sb.WriteString(fmt.Sprintf("# Total Charge              = %6.2f kWh\n", x.totalCharge))
	sb.WriteString(fmt.Sprintf("# Total Discharge           = %6.2f kWh\n", x.totalDischarge))
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() || len(x.sn) == 0 {
		sb.WriteString("# No battery ESU1 or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()

		tags := fmt.Sprintf("model=%q,sn=%q,esu=\"1\",esu_sn=%q", id.model, id.sn, x.sn)

		sb.WriteString(fmt.Sprintf("sun2000_ess_running_status{%s} %d\n", tags, x.runningStatus))
		sb.WriteString(fmt.Sprintf("sun2000_ess_charge_and_discharge_power{%s,unit=\"kW\"} %6.3f\n", tags, x.chargeAndDischargePower))
		sb.WriteString(fmt.Sprintf("sun2000_ess_bus_voltage{%s,unit=\"V\"} %3.1f\n", tags, x.busVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_ess_bus_current{%s,unit=\"A\"} %3.1f\n", tags, x.busCurrent))
		sb.WriteString(fmt.Sprintf("sun2000_ess_soc{%s,unit=\"%%\"} %3.1f\n", tags, x.batterySOC))
		sb.WriteString(fmt.Sprintf("sun2000_ess_remaining_charge_discharge_time{%s,unit=\"mins\"} %d\n", tags, x.remainingChargeDischargeTime))
		sb.WriteString(fmt.Sprintf("sun2000_ess_temperature{%s,unit=\"℃\"} %3.1f\n", tags, x.batteryTemperature))
		sb.WriteString(fmt.Sprintf("sun2000_ess_working_mode{%s} %d\n", tags, x.workingMode))
		sb.WriteString(fmt.Sprintf("sun2000_ess_rated_charge_power{%s,unit=\"W\"} %d\n", tags, x.ratedChargePower))
		sb.WriteString(fmt.Sprintf("sun2000_ess_rated_discharge_power{%s,unit=\"W\"} %d\n", tags, x.ratedDischargePower))
		sb.WriteString(fmt.Sprintf("sun2000_ess_fault_id{%s} %d\n", tags, x.faultID))
		sb.WriteString(fmt.Sprintf("sun2000_ess_maximum_charge_power{%s,unit=\"W\"} %d\n", tags, x.maximumChargePower))
		sb.WriteString(fmt.Sprintf("sun2000_ess_maximum_discharge_power{%s,unit=\"W\"} %d\n", tags, x.maximumDischargePower))
		sb.WriteString(fmt.Sprintf("sun2000_ess_current_day_charge_capacity{%s,unit=\"kWh\"} %6.2f\n", tags, x.currentDayChargeCapacity))
		sb.WriteString(fmt.Sprintf("sun2000_ess_current_day_discharge_capacity{%s,unit=\"kWh\"} %6.2f\n", tags, x.currentDayDischargeCapacity))
		sb.WriteString(fmt.Sprintf("sun2000_ess_total_charge{%s,unit=\"kWh\"} %6.2f\n", tags, x.totalCharge))
		sb.WriteString(fmt.Sprintf("sun2000_ess_total_discharge{%s,unit=\"kWh\"} %6.2f\n", tags, x.totalDischarge))
	}
	sb.WriteString("\n")

	for i := range x.pack {
		if len(x.pack[i].sn) != 0 {
			sb.WriteString(x.pack[i].metricsString(id))
		}
	}

	return sb.String()
}

type esuRunningStatus uint16

func (x esuRunningStatus) String() string {
	switch x {
	case 0:
		return "offline"
	case 1:
		return "stand-by"
	case 2:
		return "running"
	case 3:
		return "fault"
	case 4:
		return "sleep mode"
	default:
		return ""
	}
}

type esuWorkingMode uint16

func (x esuWorkingMode) String() string {
	switch x {
	case 0:
		return "none"
	case 1:
		return "Forcible charge/discharge"
	case 2:
		return "Time of Use(LG)"
	case 3:
		return "Fixed charge/discharge"
	case 4:
		return "Maximise selfconsumption"
	case 5:
		return "Fully fed to grid"
	case 6:
		return "Time of Use(LUNA2000)"
	case 7:
		return "remote scheduling maximum self-use"
	case 8:
		return "remote scheduling - full Internet access"
	case 9:
		return "remote scheduling - TOU"
	case 10:
		return "AI energy management and scheduling"
	default:
		return ""
	}
}

type esu2Data struct {
	genericData

	// 37700 STR 10
	sn string
	// 37738 U16 1 gain 10 %
	batterySOC float32
	// 37741 U16 1
	runningStatus esuRunningStatus
	// 37743 I32 2 gain 1 W
	chargeAndDischargePower float32 // > 0: charging < 0: discharging
	// 37746 U32 2 gain 100 kWh
	currentDayChargeCapacity float32
	// 37748 U32 2 gain 100 kWh
	currentDayDischargeCapacity float32
	// 37750 U16 1 gain 10 V
	busVoltage float32
	// 37751 I16 1 gain 10 A
	busCurrent float32
	// 37752 I16 1
	batteryTemperature float32
	// 37753 U32 2 gain 100 kWh
	totalCharge float32
	// 37755 U32 2 gain 100 kWh
	totalDischarge float32

	pack [3]batteryData
}

func (x *esu2Data) parse(data []byte) (err error) {
	size := 57
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var u16 uint16
	var i16 int16
	var i32 int32
	var u32 uint32
	var idx uint

	x.sn, idx, _ = getSTR(data, idx, 10)

	idx, _ = skipRecords(data, idx, 28)

	u16, idx, _ = getU16(data, idx)
	x.batterySOC = float32(u16) / 10

	idx, _ = skipRecords(data, idx, 2)

	u16, idx, _ = getU16(data, idx)
	x.runningStatus = esuRunningStatus(u16)

	idx, _ = skipRecords(data, idx, 1)

	i32, idx, _ = getI32(data, idx)
	x.chargeAndDischargePower = float32(i32)

	idx, _ = skipRecords(data, idx, 1)

	u32, idx, _ = getU32(data, idx)
	x.currentDayChargeCapacity = float32(u32) / 100
	u32, idx, _ = getU32(data, idx)
	x.currentDayDischargeCapacity = float32(u32) / 100

	u16, idx, _ = getU16(data, idx)
	x.busVoltage = float32(u16) / 10
	i16, idx, _ = getI16(data, idx)
	x.busCurrent = float32(i16) / 10
	i16, idx, _ = getI16(data, idx)
	x.batteryTemperature = float32(i16) / 10

	u32, idx, _ = getU32(data, idx)
	x.totalCharge = float32(u32) / 100
	u32, _, _ = getU32(data, idx)
	x.totalDischarge = float32(u32) / 100

	return nil
}

func (x *esu2Data) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# ESU2 Data\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("# SN                         = %q\n", x.sn))
	sb.WriteString(fmt.Sprintf("# Running Status             = %d\t%s\n", x.runningStatus, x.runningStatus))
	var txt string
	if x.chargeAndDischargePower > 0 {
		txt = "charging"
	} else if x.chargeAndDischargePower < 0 {
		txt = "discharging"
	}
	sb.WriteString(fmt.Sprintf("# Charge And Discharge Power = %6.3f kW %s\n", x.chargeAndDischargePower, txt))
	sb.WriteString(fmt.Sprintf("# Bus Voltage                = %3.1f V\n", x.busVoltage))
	sb.WriteString(fmt.Sprintf("# Bus Current                = %3.1f A\n", x.busCurrent))
	sb.WriteString(fmt.Sprintf("# Battery SOC                = %3.1f %%\n", x.batterySOC))
	sb.WriteString(fmt.Sprintf("# Battery Temperature        = %3.1f ℃\n", x.batteryTemperature))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Current Day Charge Capacity     = %6.2f kWh\n", x.currentDayChargeCapacity))
	sb.WriteString(fmt.Sprintf("# Current Day Discharge Capacity  = %6.2f kWh\n", x.currentDayDischargeCapacity))
	sb.WriteString(fmt.Sprintf("# Total Charge                    = %6.2f kWh\n", x.totalCharge))
	sb.WriteString(fmt.Sprintf("# Total Discharge                 = %6.2f kWh\n", x.totalDischarge))
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() || len(x.sn) == 0 {
		sb.WriteString("# No battery ESU2 or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()

		tags := fmt.Sprintf("model=%q,sn=%q,esu=\"2\",esu_sn=%q", id.model, id.sn, x.sn)

		sb.WriteString(fmt.Sprintf("sun2000_ess_running_status{%s} %d\n", tags, x.runningStatus))
		sb.WriteString(fmt.Sprintf("sun2000_ess_charge_and_discharge_power{%s,unit=\"kW\"} %6.3f\n", tags, x.chargeAndDischargePower))
		sb.WriteString(fmt.Sprintf("sun2000_ess_bus_voltage{%s,unit=\"V\"} %3.1f\n", tags, x.busVoltage))
		sb.WriteString(fmt.Sprintf("sun2000_ess_bus_current{%s,unit=\"A\"} %3.1f\n", tags, x.busCurrent))
		sb.WriteString(fmt.Sprintf("sun2000_ess_soc{%s,unit=\"%%\"} %3.1f\n", tags, x.batterySOC))
		sb.WriteString(fmt.Sprintf("sun2000_ess_temperature{%s,unit=\"℃\"} %3.1f\n", tags, x.batteryTemperature))
		sb.WriteString(fmt.Sprintf("sun2000_ess_current_day_charge_capacity{%s,unit=\"kWh\"} %6.2f\n", tags, x.currentDayChargeCapacity))
		sb.WriteString(fmt.Sprintf("sun2000_ess_current_day_discharge_capacity{%s,unit=\"kWh\"} %6.2f\n", tags, x.currentDayDischargeCapacity))
		sb.WriteString(fmt.Sprintf("sun2000_ess_total_charge{%s,unit=\"kWh\"} %6.2f\n", tags, x.totalCharge))
		sb.WriteString(fmt.Sprintf("sun2000_ess_total_discharge{%s,unit=\"kWh\"} %6.2f\n", tags, x.totalDischarge))
	}
	sb.WriteString("\n")

	for i := range x.pack {
		if len(x.pack[i].sn) != 0 {
			sb.WriteString(x.pack[i].metricsString(id))
		}
	}

	return sb.String()
}

type batteryData struct {
	genericData

	esuId int
	id    int

	// 38200 STR 10
	sn string
	// 38210 STR 15
	firmwareVersion string
	// 38228 U16 1
	workingStatus uint16
	// 38229 U16 1 gain 10 %
	soc float32
	// 38233 I32 2 gain 1000 kW
	chargeDischargePower float32 // > 0: charging < 0: discharging
	// 38235 U16 1 gain 10 V
	voltage float32
	// 38236 I16 1 gain 10 A
	current float32
	// 38238 U32 2 gain 100 kWh
	totalCharge float32
	// 38240 U32 2 gain 100 kWh
	totalDischarge float32
}

func (x *batteryData) parse(data []byte) (err error) {
	size := 45
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var u16 uint16
	var i16 int16
	var i32 int32
	var u32 uint32

	var idx uint

	x.sn, idx, _ = getSTR(data, idx, 10)
	x.firmwareVersion, idx, _ = getSTR(data, idx, 15)

	idx, _ = skipRecords(data, idx, 3)

	u16, idx, _ = getU16(data, idx)
	x.workingStatus = u16

	u16, idx, _ = getU16(data, idx)
	x.soc = float32(u16) / 10

	idx, _ = skipRecords(data, idx, 3)

	i32, idx, _ = getI32(data, idx)
	x.chargeDischargePower = float32(i32) / 1000

	u16, idx, _ = getU16(data, idx)
	x.voltage = float32(u16) / 10
	i16, idx, _ = getI16(data, idx)
	x.current = float32(i16) / 10

	idx, _ = skipRecords(data, idx, 1)

	u32, idx, _ = getU32(data, idx)
	x.totalCharge = float32(u32) / 100
	u32, _, _ = getU32(data, idx)
	x.totalDischarge = float32(u32) / 100

	return nil
}

func getESUSN(esuId int) (out string) {
	switch esuId {
	case 1:
		return parsedData.esu1.sn
	case 2:
		return parsedData.esu2.sn
	default:
		return ""
	}
}

func (x *batteryData) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()
	sb.WriteString(fmt.Sprintf("# ESU %d / Pack %d\n", x.esuId, x.id))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# SN                         = %q\n", x.sn))
	sb.WriteString(fmt.Sprintf("# Firmware Version           = %q\n", x.firmwareVersion))
	sb.WriteString(fmt.Sprintf("# Working Status             = %d\n", x.workingStatus))
	var txt string
	if x.chargeDischargePower > 0 {
		txt = "charging"
	} else if x.chargeDischargePower < 0 {
		txt = "discharging"
	}
	sb.WriteString(fmt.Sprintf("# Charge And Discharge Power = %6.3f kW %s\n", x.chargeDischargePower, txt))
	sb.WriteString(fmt.Sprintf("# Voltage                    = %3.1f V\n", x.voltage))
	sb.WriteString(fmt.Sprintf("# Current                    = %3.1f A\n", x.current))
	sb.WriteString(fmt.Sprintf("# SOC                        = %3.1f %%\n", x.soc))
	sb.WriteString(fmt.Sprintf("# Total Charge               = %6.2f kWh\n", x.totalCharge))
	sb.WriteString(fmt.Sprintf("# Total Discharge            = %6.2f kWh\n", x.totalDischarge))
	sb.WriteString("\n")

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() || len(x.sn) == 0 {
		sb.WriteString(fmt.Sprintf("# No battery ESU%d/Pack%d or identification data read yet\n", x.esuId, x.id))
	} else {
		id.RLock()
		defer id.RUnlock()

		esuSN := getESUSN(x.esuId)
		tags := fmt.Sprintf("model=%q,sn=%q,esu=\"%d\",esu_sn=%q,pack=\"%d\",pack_sn=%q", id.model, id.sn, x.esuId, esuSN, x.id, x.sn)

		sb.WriteString(fmt.Sprintf("sun2000_ess_pack_working_status{%s} %d\n", tags, x.workingStatus))
		sb.WriteString(fmt.Sprintf("sun2000_ess_pack_charge_and_discharge_power{%s,unit=\"kW\"} %6.3f\n", tags, x.chargeDischargePower))
		sb.WriteString(fmt.Sprintf("sun2000_ess_pack_voltage{%s,unit=\"V\"} %3.1f\n", tags, x.voltage))
		sb.WriteString(fmt.Sprintf("sun2000_ess_pack_current{%s,unit=\"A\"} %3.1f\n", tags, x.current))
		sb.WriteString(fmt.Sprintf("sun2000_ess_pack_soc{%s,unit=\"%%\"} %3.1f\n", tags, x.soc))
		sb.WriteString(fmt.Sprintf("sun2000_ess_pack_total_charge{%s,unit=\"kWh\"} %6.2f\n", tags, x.totalCharge))
		sb.WriteString(fmt.Sprintf("sun2000_ess_pack_total_discharge{%s,unit=\"kWh\"} %6.2f\n", tags, x.totalDischarge))
	}
	sb.WriteString("\n")

	return sb.String()
}

type esuTemperaturesData struct {
	genericData

	esu [2]struct {
		pack [3]struct {
			// 38452 INT16 1 gain 10 ℃
			maxTemperature float32
			// 38453 INT16 1 gain 10 ℃
			minTemperature float32
		}
	}
}

func (x *esuTemperaturesData) parse(data []byte) (err error) {
	size := 2 * 3 * 2
	if len(data) < size {
		return fmt.Errorf("data length %d < %d", len(data), size)
	}

	x.Lock()
	defer x.Unlock()

	var i16 int16
	var idx uint

	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			i16, idx, _ = getI16(data, idx)
			x.esu[i].pack[j].maxTemperature = float32(i16) / 10
			i16, idx, _ = getI16(data, idx)
			x.esu[i].pack[j].minTemperature = float32(i16) / 10
		}
	}

	return nil
}

func (x *esuTemperaturesData) metricsString(id *identificationData) string {
	sb := strings.Builder{}
	x.RLock()
	defer x.RUnlock()

	sb.WriteString("# Battery ESU Temperatures Data\n")
	sb.WriteString(fmt.Sprintf("# Last Read = %s\n", x.lastRead.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("# Next Read = %s\n", x.nextRead.Format(time.RFC3339)))
	sb.WriteString("\n")

	for i := 0; i < 2; i++ {
		esuSN := getESUSN(i + 1)
		if len(esuSN) == 0 {
			continue
		}
		for j := 0; j < 3; j++ {
			if (i == 0 && len(parsedData.esu1.pack[j].sn) > 0) ||
				(i == 1 && len(parsedData.esu2.pack[j].sn) > 0) {
				sb.WriteString(fmt.Sprintf("# ESU %d / Pack %d\n", i+1, j+1))
				sb.WriteString(fmt.Sprintf("# Max Temperature = %3.1f ℃\n", x.esu[i].pack[j].maxTemperature))
				sb.WriteString(fmt.Sprintf("# Min Temperature = %3.1f ℃\n", x.esu[i].pack[j].minTemperature))
				sb.WriteString("\n")
			}
		}
	}

	// skip metrics if the data is empty
	if x.lastRead.IsZero() || id.lastRead.IsZero() {
		sb.WriteString("# No battery ESU temperatures or identification data read yet\n")
	} else {
		id.RLock()
		defer id.RUnlock()

		for i := 0; i < 2; i++ {

			esuSN := getESUSN(i + 1)
			if len(esuSN) == 0 {
				continue
			}
			for j := 0; j < 3; j++ {
				var pack *batteryData
				switch i {
				case 0:
					pack = &parsedData.esu1.pack[j]
				case 1:
					pack = &parsedData.esu2.pack[j]
				}
				if len(pack.sn) == 0 {
					continue
				}
				tags := fmt.Sprintf("model=%q,sn=%q,esu=\"%d\",esu_sn=%q,pack=\"%d\",pack_sn=%q", id.model, id.sn, i+1, esuSN, j+1, pack.sn)

				sb.WriteString(fmt.Sprintf("sun2000_ess_pack_max_temperature{%s,unit=\"℃\"} %3.1f\n", tags, x.esu[i].pack[j].maxTemperature))
				sb.WriteString(fmt.Sprintf("sun2000_ess_pack_min_temperature{%s,unit=\"℃\"} %3.1f\n", tags, x.esu[i].pack[j].minTemperature))

			}
		}

	}

	sb.WriteString("\n")

	return sb.String()
}
