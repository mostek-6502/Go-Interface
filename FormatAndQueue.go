// Rick Faszold
//
// September 26th, 2016
//
// XEn, LLC (c)   Missouri, USA
//
// Raw Data to JSON
//    A. Takes the Board data off of the channel,  IE.
//    B.1 Breaks each element up into discreet JSON messages
//    B.2 While each message is being broken up into discreet elements,
//        parts of the data are changed from Status Codes to Readbale Text
//       Ex.
//         "WINDSTAT",2,33,44  is broken up into:
//             JSON { "Element: Header", "Data: WINDSTAT"	}
//             JSON { "Element: Speed", "Data: 2"			}
//             JSON { "Element: Status", "Data: OK"			}
//             JSON { "Element: Direction", "Data: NNW"		}
//          The down stream Python script simply takes the preformaated JSON message and
//          sends it to the browser
//    C. Each JSON Message is loaded onto a JSON Channel


package main

import (
	"fmt"
	"strings"
	"time"
	"strconv"
)


var MESSAGE_IDs = []string { "Header",
							"Error_Message", }

var EEPROM_IDs = []string { "Header",
							"EEPROM_Version",
							"EEPROM_RebootCount",
							"EEPROM_ProcessorSpeed",
							"EEPROM_TemperatureResolution",
							"EEPROM_HorizontalAdjustment",
							"EEPROM_VerticalAdjustment",
							"EEPROM_Max_Wind_Speed",
							"EEPROM_Max_Wind_Speed_Delay",
							"EEPROM_UART_Telemetry",
							"EEPROM_UDP_Telemetry",
							"EEPROM_Use_DNS",
							"EEPROM_UDP_Port",
							"EEPROM_Server_Name",
							"EEPROM_Dark_Threshold",
							"EEPROM_Delay_Sudden_Moveback",
							"EEPROM_Move_Threshold" }

var SENSOR_IDs = []string { "Header",
							"SEN_WSFP_WindSpeed",
							"SEN_WSFP_WindStatus",
							"SEN_WSFP_Failsafe",
							"SEN_WSFP_ProxRight",
							"SEN_WSFP_ProxLeft",
							"SEN_WSFP_ProxUp",
							"SEN_WSFP_ProxDown" }

var PUMP_IDs = []string { "Header",
							"PUMP_Thermostat_On",
							"PUMP_DishToResevoirPWMPercent",
							"PUMP_DishToResevoirStatus",
							"PUMP_Immd_ResevoirToHousePWMPercent",
							"PUMP_Immd_ResevoirToHouseStatus",
							"PUMP_Hold_ResevoirToHousePWMPercent",
							"PUMP_Hold_ResevoirToHouseStatus",
							"PUMP_AUXToResevoirPWMPercent",
							"PUMP_AUXToResevoirStatus" }


var RUNTIME_IDs = []string { "Header",
							"RUNTIME_Elapsed_Seconds",
							"RUNTIME_Temperature_Count",
							"RUNTIME_Dish_Count",
							"RUNTIME_Pump_Count",
							"RUNTIME_ADC_Count",
							"RUNTIME_UDP_Listen",
							"RUNTIME_Request_Heat",
							"RUNTIME_TT_Heat_Requested",
							"RUNTIME_EEPROM_Count",
							"RUNTIME_DataOuput_Count",
							"RUNTIME_TempPumps_Second",
							"RUNTIME_ADCDish_Second" }

var DISH_IDs = []string { "Header",
							"DISH_RightLow",
							"DISH_LeftLow",
							"DISH_RightHigh",
							"DISH_LeftHigh",
							"DISH_H_ResultCalc",
							"DISH_H_MoveFlag",
							"DISH_UpLow",
							"DISH_DownLow",
							"DISH_UpHigh",
							"DISH_DownHigh",
							"DISH_V_ResultCalc",
							"DISH_V_MoveFlag" }


var TEMP_IDs = []string { "Header",
							"TEMP_M_0",
							"TEMP_C_0",
							"TEMP_F_0",
							"TEMP_R_0",
							"TEMP_M_1",
							"TEMP_C_1",
							"TEMP_F_1",
							"TEMP_R_1",
							"TEMP_M_2",
							"TEMP_C_2",
							"TEMP_F_2",
							"TEMP_R_2",
							"TEMP_M_3",
							"TEMP_C_3",
							"TEMP_F_3",
							"TEMP_R_3",
							"TEMP_M_4",
							"TEMP_C_4",
							"TEMP_F_4",
							"TEMP_R_4",
							"TEMP_M_5",
							"TEMP_C_5",
							"TEMP_F_5",
							"TEMP_R_5",
							"TEMP_M_6",
							"TEMP_C_6",
							"TEMP_F_6",
							"TEMP_R_6",
							"TEMP_M_7",
							"TEMP_C_7",
							"TEMP_F_7",
							"TEMP_R_7",
							"TEMP_M_8",
							"TEMP_C_8",
							"TEMP_F_8",
							"TEMP_R_8",
							"TEMP_M_9",
							"TEMP_C_9",
							"TEMP_F_9",
							"TEMP_R_9",
							"TEMP_M_A",
							"TEMP_C_A",
							"TEMP_F_A",
							"TEMP_R_A",
							"TEMP_M_B",
							"TEMP_C_B",
							"TEMP_F_B",
							"TEMP_R_B",
							"TEMP_M_C",
							"TEMP_C_C",
							"TEMP_F_C",
							"TEMP_R_C",
							"TEMP_M_D",
							"TEMP_C_D",
							"TEMP_F_D",
							"TEMP_R_D",
							"TEMP_M_E",
							"TEMP_C_E",
							"TEMP_F_E",
							"TEMP_R_E",
							"TEMP_M_F",
							"TEMP_C_F",
							"TEMP_F_F",
							"TEMP_R_F" }


var mDishStatus = map[string]string {
							"0" : "MOVE_OK_TO_MOVE",
							"1" : "MOVE_H_RIGHT",
							"2" : "MOVE_H_LEFT",
							"3" : "MOVE_V_UP",
							"4" : "MOVE_V_DOWN",
							"5" : "MOVE_TRANSITION",
							"6" : "MOVE_ERROR_WIND_SPEED",
							"7" : "MOVE_H_ERROR_DUAL_PROXIMITY",
							"8" : "MOVE_V_ERROR_DUAL_PROXIMITY",
							"9" : "NO_H_MOVEMENT_NEEDED",
							"10" : "NO_V_MOVEMENT_NEEDED",
							"11" : "NO_MOVE_PROXIMITY_DETECT_RIGHT",
							"12" : "NO_MOVE_PROXIMITY_DETECT_LEFT",
							"13" : "NO_MOVE_PROXIMITY_DETECT_UP",
							"14" : "NO_MOVE_PROXIMITY_DETECT_DOWN",
							"15" : "NO_H_MOVEMENT_PERCENT_RANGE",
							"16" : "NO_V_MOVEMENT_PERCENT_RANGE",
							"17" : "NO_MOVEMENT_FAILSAFE",
							"18" : "NO_MOVEMENT_MOTORS_OFF",
							"19" : "FORCE_MOVE_H_RIGHT",
							"20" : "FORCE_MOVE_H_LEFT",
							"21" : "FORCE_MOVE_V_UP",
							"22" : "FORCE_MOVE_V_DOWN",
							"23" : "FORCE_MOVE_H_PAUSE",
							"24" : "FORCE_MOVE_V_PAUSE",
							"25" : "TEMP_TOO_HIGH_MOVE_AWAY_FROM_SUN_H",
							"26" : "TEMP_TOO_HIGH_MOVE_AWAY_FROM_SUN_V",
							"27" : "ERROR_HORIZONTAL_MOVEMENT_FLAG",
							"28" : "ERROR_VERTICAL_MOVEMENT_FLAG",
							"29" : "ERROR_PHOTO_RESISTOR_OUT_OF_RANGE",
							"30" : "Unknown", }

var mTemperatureStatus = map[string]string{
							"0"  : "OK",
							"20" : "Chip Reset",
							"30" : "Channel-1",
							"40" : "ROM Rtrv",
							"50" : "No Cfg",
							"60" : "Temp Strt",
							"85" : "Channel-2",
							"90" : "Temp Rtrv",
							"99" : "Unknown", }

var mWindAlert = map[string]string {
							"0" : "OK",
							"1" : "Caution",
							"2" : "Alert", }

var mPumpStatus = map[string]string {
							"0" : "PUMP_NORMAL_OPERATIONS",
							"1" : "PUMP_OFF",
							"2" : "PUMP_TEMPERATURE_ERROR",
							"3" : "PUMP_TEMPERATURE_TOO_LOW",
							"4" : "PUMP_TEMPERATURE_TOO_HIGH",
							"5" : "PUMP_TO_HOUSE_THERMOSTAT_OFF",
							"6" : "PUMP_FAILSAFE_ERROR",
							"7" : "PUMP_HOUSE_TEMPERATURE_ERROR",
							"8" : "PUMP_HOUSE_TEMP_TOO_HIGH", }

var mThermostatOn_Off = map[string]string {
							"0" : "Off",
							"1" : "On", }

var mFailsafe = map[string]string {
							"0" : "Failed!",
							"1" : "OK", }


var mProximity = map[string]string {
							"0" : "OK",
							"1" : "Limit!", }



//  print('event: message\n' + 'data: {"Element" : "EEPROM_RebootCount", "Data" : ' + '"' + str(iCount) + '"' + '}\n\n')

var tcks string = "\""
var JSON_Start =  "event: message\ndata: {" + tcks + "Element" + tcks + " : " + tcks
var JSON_Middle = tcks + ", " + tcks + "Data" + tcks + " : " + tcks
var JSON_End = tcks + "}\n\n"

func Format_JSON_String(sID string, sValue string) string {

	strJSON := JSON_Start + sID + JSON_Middle + sValue + JSON_End

	return strJSON

}


func Process_EEPROM_Data(sData string, aID []string) []string {

	var sElement, sTemp string
	asReturn := make([]string, 0, 100)

	sSplit := strings.Split(sData, ",")

	for iIndex := 0; iIndex < len(EEPROM_IDs); iIndex++ {
		sTemp = strings.TrimSpace(sSplit[iIndex])

		sElement = Format_JSON_String(aID[iIndex], sTemp)

		asReturn = append(asReturn, sElement)
	}

	return asReturn
}


func Process_Data_Message(sData string, aID []string) [] string {

	var sElement, sTemp string
	asReturn := make([]string, 0, 100)

	sSplit := strings.Split(sData, ",")

	for iIndex := 0; iIndex < len(MESSAGE_IDs); iIndex++ {
		sTemp = strings.TrimSpace(sSplit[iIndex])

		sElement = Format_JSON_String(aID[iIndex], sTemp)

		//println("2. Format->", sElement)

		asReturn = append(asReturn, sElement)
	}

	return asReturn


}


func Process_Pump_Data(sData string, aID []string) []string {

	var sElement, sTemp string

	asReturn := make([]string, 0, 100)

	sSplit := strings.Split(sData, ",")

	for iIndex := 0; iIndex < len(PUMP_IDs); iIndex++ {

		sTemp = strings.TrimSpace(sSplit[iIndex])

		if iIndex == 1 {
			sTemp = mThermostatOn_Off[sTemp]
		} else if iIndex == 2 || iIndex == 4 || iIndex == 6 || iIndex == 8 {
			sTemp = sTemp + "%"
		} else if iIndex == 3 || iIndex == 5 || iIndex == 7 || iIndex == 9 {
			sTemp = mPumpStatus[sTemp]
		}

		sElement = Format_JSON_String(aID[iIndex], sTemp)

		asReturn = append(asReturn, sElement)
	}

	return asReturn
}



func Process_Sensor_Data(sData string, aID []string) []string {

	var sWindSpeed string
	var sTemp string
	var sElement string
	asReturn := make([]string, 0, 100)

	sSplit := strings.Split(sData, ",")

	for iIndex := 0; iIndex < len(SENSOR_IDs); iIndex++ {

		sTemp = strings.TrimSpace(sSplit[iIndex])

		if iIndex == 0  {
			sElement = Format_JSON_String(aID[iIndex], sTemp) 		// Header
			asReturn = append(asReturn, sElement)
		} else if (iIndex == 1) || (iIndex == 2) {

			// wind speed and wind speed status have to go together because of color changes....
			if iIndex == 1 {
				sWindSpeed = sTemp
			}  else if iIndex == 2 {

				sElement = Format_JSON_String(aID[iIndex - 1], sWindSpeed)
				asReturn = append(asReturn, sElement)

				sElement = Format_JSON_String(aID[iIndex], mWindAlert[sTemp])
				asReturn = append(asReturn, sElement)
			}
		} else if iIndex == 3 {
			sTemp = mFailsafe[sTemp]

			sElement = Format_JSON_String(aID[iIndex], sTemp)		// Failsafe
			asReturn = append(asReturn, sElement)
		} else {													// The rest of these are proximities
			sTemp = mProximity[sTemp]

			sElement = Format_JSON_String(aID[iIndex], sTemp)
			asReturn = append(asReturn, sElement)
		}
	}

	//fmt.Println("Process_Sensor_Data Complete!")

	return asReturn
}



func Process_Dish_Data(sData string, aID []string) []string {

	var sElement, sTemp string
	asReturn := make([]string, 0, 100)

	sSplit := strings.Split(sData, ",")

	for iIndex := 0; iIndex < len(DISH_IDs); iIndex++ {
		sTemp = strings.TrimSpace(sSplit[iIndex])

		if iIndex == 6 { sTemp = mDishStatus[sTemp] }
		if iIndex == 12 { sTemp = mDishStatus[sTemp] }

		sElement = Format_JSON_String(aID[iIndex], sTemp)

		asReturn = append(asReturn, sElement)
	}

	return asReturn
}


func Process_Data_RunTime(sData string, aID []string) []string {

	var iUpTime int
	var iTempCycles int
	var iPumpCycles int

	var sDays string

	var iSeconds int
	var iMinutes int
	var iHours int
	var iDays int


	var sElement, sTemp string
	asReturn := make([]string, 0, 100)

	sSplit := strings.Split(sData, ",")

	//CheckLength("Process_RunTime_Data()", sSplit, RUNTIME_IDs)

	// we are subtracting 2 because the data does not support 2 more elements, only the computed output...
	for iIndex := 0; iIndex < len(RUNTIME_IDs) - 2; iIndex++ {

		//iDays := 0
		//if (iTime > 86400) {
		//	iDays = iTime / 86400;  // number of days.
		//	iTime = iTime % 86400;  // remainder after "subtracting" out days
		//}
		//String strLapsedTime = String.format("%d %02d:%02d:%02d", iDays, iTime / 3600, (iTime / 60) % 60, iTime % 60);

		sTemp = strings.TrimSpace(sSplit[iIndex])
		if iIndex == 1 {  // this is the Duration in Seconds that the Board has been up

			// this needs to be cleaned up
			iSeconds, _ = strconv.Atoi(sTemp)
			iUpTime = iSeconds

			iDays = iSeconds / 86400
			iHours = (iSeconds - (iDays * 86400)) / 3600
			iMinutes = (iSeconds - ((iDays * 86400) + (iHours * 3600))) / 60

			if len(sDays) == 0 {
				sTemp = fmt.Sprintf( "%02d:%02d:%02d", iHours, iMinutes, iSeconds % 60)
			} else {
				sTemp = fmt.Sprintf("%d.%02d:%02d:%02d", iDays, iHours, iMinutes, iSeconds % 60)
			}

			sElement = Format_JSON_String(aID[iIndex], sTemp)
			asReturn = append(asReturn, sElement)
		} else {

			sElement = Format_JSON_String(aID[iIndex], sTemp)
			asReturn = append(asReturn, sElement)

			if iIndex == 2 {
				iTempCycles, _ = strconv.Atoi(sTemp)
			} else if iIndex == 4 {
				iPumpCycles, _ = strconv.Atoi(sTemp)
			}
		}
	}

	var iIndex int
	var f float64
	var s64 string

	if iUpTime > 0 {
		f = (float64)(iTempCycles / iUpTime)
	} else {
		f = 0.0
	}
	s64 = strconv.FormatFloat(f, 'f', 2, 64)

	iIndex = 11
	sElement = Format_JSON_String(aID[iIndex], s64)
	asReturn = append(asReturn, sElement)


	if iUpTime > 0 {
		f = (float64)(iPumpCycles / iUpTime)
	} else {
		f = 0.0
	}
	s64 = strconv.FormatFloat(f, 'f', 2, 64)

	iIndex = 12
	sElement = Format_JSON_String(aID[iIndex], s64)
	asReturn = append(asReturn, sElement)


	//fmt.Println("Process_Data_RunTime Complete!")

	return asReturn
}



func Process_Temp_Data(sData string, aID []string) []string {

	var sElement string
	asReturn := make([]string, 0, 200)

	sSplit := strings.Split(sData, ",")

	var iOffset, iOutputOffset int
	var sTemp string


	//CheckLength("Process_Temp_Data()", sSplit, TEMP_IDs)

	// it's just a lot easier to take the header out of this... and hard code it here...
	sElement = Format_JSON_String(aID[iOutputOffset], "TEMPS")
	asReturn = append(asReturn, sElement)

	for iIndex := 0; iIndex < 16; iIndex++ {

		// there are 16 groups of data to be processed in each message
		// within each group, there are  F =
		iOffset = (iIndex * 8) + 1


		//fmt.Println("TEMP Index: ", iIndex, "Data Offset:", iOffset, "Output Offset: ", iOutputOffset)


		// STATUS
		iOutputOffset++
		sTemp = strings.TrimSpace(sSplit[iOffset + 0])
		sTemp = mTemperatureStatus[sTemp]
		sElement = Format_JSON_String(aID[iOutputOffset], sTemp)
		asReturn = append(asReturn, sElement)


		// DEGREES C
		iOutputOffset++
		sWhole := strings.TrimSpace(sSplit[iOffset + 1])
		sFraction := sSplit[iOffset + 2]
		sSign := sSplit[iOffset + 3]

		sTemp = ""
		if sSign == "1" {
			sTemp = "-"
		}
		sTemp = strings.TrimSpace(sTemp) + strings.TrimSpace(sWhole) + "." + strings.TrimSpace(sFraction)

		sElement = Format_JSON_String(aID[iOutputOffset], sTemp)
		asReturn = append(asReturn, sElement)


		// DEGREES F
		iOutputOffset++
		sWhole = sSplit[iOffset + 4]
		sFraction = sSplit[iOffset + 5]
		sSign = sSplit[iOffset + 6]

		sTemp = ""
		if sSign == "1" {
			sTemp = "-"
		}
		sTemp = strings.TrimSpace(sTemp) + strings.TrimSpace(sWhole) + "." + strings.TrimSpace(sFraction)


		sElement = Format_JSON_String(aID[iOutputOffset], sTemp)
		asReturn = append(asReturn, sElement)


		// ROM CODES
		iOutputOffset++
		sROM := sSplit[iOffset + 7]
		// 0011223344556677 - 00-11-22-33-44-55-66-77

		sROM = sROM[0:16]

		sTemp = sROM[0:2] + "-" + sROM[2:4] + "-" + sROM[4:6] + "-" + sROM[6:8] + "-" + sROM[8:10] + "-" + sROM[10:12] + "-" + sROM[12:14] + "-" + sROM[14:16]

		// fmt.Println("Process ROM String->", sROM, "<-  ->", sTemp, "<-")


		sElement = Format_JSON_String(aID[iOutputOffset], sTemp)
		asReturn = append(asReturn, sElement)
	}

	// fmt.Println("Process_Temp_Data Complete!")

	return asReturn
}




func FormatAndQueueData(ChannelBoardData <-chan string, ChannelBrowserData chan<- string) {

	var sPlace string = "FormatAndQueueData()"

	var bOK bool
	var sData string
	var sFormattedData []string

	var bFirstTime bool
	bFirstTime = true

	fmt.Println(sPlace, "Wait For Board Data...")

	for {

		if ActiveWebConnections() == false {
			if len(ChannelBoardData) > 0 {
				fmt.Println(sPlace, "There are no Browser Connections.  Draining The Board Channel!  Queue Depth: ", len(ChannelBoardData))

				for len(ChannelBoardData) > 0 {
					sData = <-ChannelBoardData
				}
			}
		}

		sData = <- ChannelBoardData

		if bFirstTime == true {
			fmt.Println(sPlace, "Processing Board Data...")
			bFirstTime = false
		}

		// call the magic
		bOK = true
		if strings.Contains(sData,  "EEPROM  ") {
			sFormattedData = Process_EEPROM_Data(sData, EEPROM_IDs)
		} else if strings.Contains(sData, "SEN-WSFP") {
			sFormattedData = Process_Sensor_Data(sData, SENSOR_IDs)
		} else if strings.Contains(sData, "PUMPS   ") {
			sFormattedData = Process_Pump_Data(sData, PUMP_IDs)
		} else if strings.Contains(sData, "DISH    ") {
			sFormattedData = Process_Dish_Data(sData, DISH_IDs)
		} else if strings.Contains(sData, "TEMPS   ") {
			sFormattedData = Process_Temp_Data(sData, TEMP_IDs)
		} else if strings.Contains(sData, "RUNTIME ") {
			sFormattedData = Process_Data_RunTime(sData, RUNTIME_IDs)
		} else if strings.Contains(sData, "MESSAGE ") {
			sFormattedData = Process_Data_Message(sData, MESSAGE_IDs)
		} else {
			bOK = false
			fmt.Println(sPlace, "Unrecognized Header from Board Data ->", sData, "<-")
		}


		if bOK == true {
			for iIndex := range sFormattedData {
				sTemp := sFormattedData[iIndex]

				//fmt.Println(sPlace, "2b. Data To Channel->", sTemp, "<-")

				ChannelBrowserData <- sTemp
			}

			// testing purposes only...
			//s := Format_JSON_String("Error_Message", "Really Really REALLY Bad Error! (just testing)")
			//ChannelBrowserData <- s
		}

		sFormattedData = nil

		time.Sleep(time.Millisecond)
	}

	// should never get here
	wg.Done()

}
