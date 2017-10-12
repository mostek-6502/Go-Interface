// Rick Faszold
//
// September 26th, 2016
//
// XEn, LLC (c)   Missouri, USA
//
// This provides two concurrent functions for sending data to a browser.
// Currently, the Board communicates on 55056 to this program.
// This program utilizes (+1) or 55057 as the initial communications port
//     for the browsers.
//
// 1. Web Browser Connections
//    This function manages the browser connections.
//    Browser connections are limited to 10, while utilizing up to 20 ports from a predefined list.
//        The reason more ports are utilized is due to refreshes causing a new session
//        and the time it takes to recycle that port.
//    Once a connection is made on the PREDEFINED port defined in the cgi-bin <dir>,
// 			a simple hand shake ensues.
//    Python Sends -> HELLO
//    This Program ->
//        Searches for an Open Port from a predefined list and when it finds one,
//        Responds with HIYA,55058 (for instance)
//           Essentially, I see you (HIAY), let's talk on 55058 from here on out
//    Python Sends -> CONFIRM ,55058
//
//    From here, all UDP communications move to 55058 or one of the next 19 ports
//    This program keeps a running list of connected Browsers
//    The reason for the port change is to only allow 'hand shakes' on the predefined port.
//        One port can not handle ALL of the two way traffice from multiple browsers.
//        This also allowed an easier way to identify what was connected
// 2. JSON Data to Browser
//		The final function removes data from the JSON Channel and iterates
//		through the list of connected browsers.
// 		Once an active browser is identified, the JSON message is sent to that browser.
// 		Thus, if 8 browsers ar on the active list, the message is sent 8 times.
//
// These concurrent functions were shared in this file to more closely associate the
// list of browser connections and what data needed to be sent to them.
//


package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"strings"
)


// shared structure between concurrent functions...
type WebServer_Connections struct {
	iStatus int					// 0 = Open, 1 = Port Reserved, 2 = Running
	tLastCheckInTime time.Time	// time of last "RUNNING"
	iPort int					// 55058, 55059, etc...
	iSentMessage int			// used to validate when a messages start flowing to the browser
	sIP string					// address of Web Server
	webSvrConn net.UDPConn		// connection to Web Server
}

var webConnections = make([]WebServer_Connections, 10)



func SetStatusToReceiveData(sIP string, iPort int) {

	var sPlace string = "SetStatusToReceiveData"

	fmt.Println(sPlace, " CONFIRMED Received From Browser", iPort)

	for iWCIdx := range webConnections {

		fmt.Println(sPlace, " Checking WebCon Port:", webConnections[iWCIdx].iPort)

		if iPort == webConnections[iWCIdx].iPort {
			webConnections[iWCIdx].iSentMessage = 0
			webConnections[iWCIdx].sIP = sIP
			webConnections[iWCIdx].iStatus = 2                   // 0 = Open, 1 = Port Reserved, 2 = Running
			webConnections[iWCIdx].tLastCheckInTime = time.Now() // time of last "RUNNING"

			fmt.Println(sPlace, "WebCon Set and Ready to Process Data")
			return
		}
	}

	fmt.Println(sPlace, "WebCon NOT Set - Data Can Not Be Processed...")

	return

}

func SetNextWebCon(iWebBasePort int) int {


	var sPlace string = "SetNextWebCon()"

	var iBasePort int
	iBasePort = iWebBasePort + 1

	// we are allocating 20 ports for this; but, we plan on using only 10 connections
	// this way, if a port is not avaialbe, then we simply go to the next port and
	// 'assume' that 20 pre-defined ports will be usable for the 10 connections
	var iPortsToUse [20]int

	for iIdx := range iPortsToUse {
		iPortsToUse[iIdx] = iBasePort
		iBasePort++
	}
	fmt.Println(sPlace, "        Base Ports Set At:", iPortsToUse)

	// create a port list, range through them, see what's being used and pick one to use
	for iCurrentPortIdx := range iPortsToUse {
		// loop through the list of what is already being allocated or in use
		for iWebConIdx := range webConnections {

			// if it is already being used, do not try to re-use it
			if webConnections[iWebConIdx].iPort == iPortsToUse[iCurrentPortIdx] {
				iPortsToUse[iCurrentPortIdx] = 0
			}
		}
	}
	fmt.Println(sPlace, "Base Ports After Clean Up:", iPortsToUse)


	// do we have any open Web Connections left in the structure?
	var iOpenWebConnections int = -1
	for iWCIdx := range webConnections {
		if webConnections[iWCIdx].iPort == 0 {
			iOpenWebConnections = iWCIdx
			break
		}
	}

	// if there are no web connections...  bail
	if iOpenWebConnections == -1 {
		fmt.Println(sPlace, " The are no Open Web Connections!")
		return -1
	}
	fmt.Println(sPlace, " WebCon Structure Index:", iOpenWebConnections)


	// do we have any open ports available?
	var iOpenServerPorts int = -1
	for iCurrentPort := range iPortsToUse {

		fmt.Println(sPlace, " WebCon Status Of Ports - Index: ", iCurrentPort, " Port: ", iPortsToUse[iCurrentPort])

		if iPortsToUse[iCurrentPort] != 0 {
			iOpenServerPorts = iPortsToUse[iCurrentPort]
			break
		}
	}

	// if there are no ports...  bail
	if iOpenServerPorts == -1 {
		fmt.Println(sPlace, " The are no Open Server Ports!")
		return -1
	}


	// we have a usable Web connection,
	// we have a usable port,
	// now let's allocate from the O/S ans start using them
	fmt.Println(sPlace, " Calling ResolveUDPAddr()  Using Port ", iOpenServerPorts)

	ServerAddr, err := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(iOpenServerPorts))
	if err != nil {
		fmt.Println(sPlace, "::ResolveUDPAddr()  Port: ", iOpenServerPorts, "Error: ", err)
		return -1
	}

	fmt.Println(sPlace, " Calling ListenUDP()")

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		fmt.Println(sPlace, "::ListenUDP()  Port: ", iOpenServerPorts, "Error: ", err)
		return -1
	}


	fmt.Println(sPlace, "Changing Port From->", webConnections[iOpenWebConnections].iPort, "<- To: ", iOpenServerPorts)

	// the addrHelloIPPort is (as the name says) from the 'general' communications port.
	// in this, we have the IP of the requesting server, but, the port needs to change
	// to the data port...

	webConnections[iOpenWebConnections].iPort = iOpenServerPorts		// this is the port number
	webConnections[iOpenWebConnections].iStatus = 1                   	// 0 = Open, 1 = Port Reserved, 2 = Running
	webConnections[iOpenWebConnections].tLastCheckInTime = time.Now() 	// time of last "RUNNING"
	webConnections[iOpenWebConnections].iSentMessage = 0				// allows a console message to confirm data being sent
	webConnections[iOpenWebConnections].webSvrConn = *ServerConn      	// connection to Web Server

	fmt.Println(sPlace, "WebCon Structure Set.  Port->", webConnections[iOpenWebConnections].iPort)

	return iOpenWebConnections
}


func ManageWebServerConnections(sIncomingPort string) {

	// open a port and listen for the python script...
	// once a python script makes a request for data from this connection.
	// a small handshake takes place
	// the connection is handed off to other port, freeing this for other connections

	var sPlace string = "ManageWebServerConnections()"
	var iWebPort int
	var buf = make([]byte, 128)


	// initialization of the structure...
	for _, webCon := range webConnections {
		webCon.iPort = 0
		webCon.iStatus = 0    // 0 = Open, 1 = Port Reserved, 2 = Running
	}


	// this is the default board port...  (55056)
	iWebPort, err := strconv.Atoi(sIncomingPort)
	if err != nil {
		iWebPort = 55056
	}
	// (55057) is the default port for listening to browser traffic.
	iWebPort++

	// this is the listening WebServer Port....  (55057) ... this always stays here...
	// this is a starting point for the browser...  (a python script)
	// the IP and Port for this intbb (server) is stored in a config file in the python script directory
	// once the BROWSER connected (and the pythin script is running), the script sends a HELLO
	// this (intbb) routine, responds with a HIYA,xxxxx  - where xxxxx is a Port such as 55058 (for example)
	// The pyton script close the original HELLO (55057) connection and opens a new connection based on xxxxx or 55058 (for example)
	// from here on out, everything coming out of the board is send to this Browser connection.
	// Every 10,000 messages, the pythin script sends a RUNNING message... this is aheart beat message.
	// once these stop and a time out is triggered, this (intbb) routine, frees the connection and recycles it.

	fmt.Println(sPlace, "::ResolveUDPAddr()  Listening For WebRequests on Port ", iWebPort)

	ServerAddr, err := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(iWebPort))
	if err != nil {
		fmt.Println(sPlace, "::ResolveUDPAddr()  Port: ", iWebPort, "Error: ", err)
		return
	}

	fmt.Println(sPlace, " Calling ListenUDP()")

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		fmt.Println(sPlace, "::ListenUDP()  Port: ", iWebPort, "Error: ", err)
		return
	}

	defer ServerConn.Close()


	fmt.Println(sPlace, "Listening For Connections....")

	var sData string

	// only two messages are processed here: HELLO and CONFIRM
	for {
		nByteCount, addrServerConversationPort, err := ServerConn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println(sPlace, "::ReadFromUDP()  Port: ", iWebPort, "Error: ", err)
			return
		}

		sData = string(buf[0:nByteCount])

		fmt.Println(sPlace, "Received:", sData, "from",  addrServerConversationPort.String())

		// need a current list of IP's trying to say hi....
		// so that we do not generate more ports than necessary
		if sData == "HELLO   " {

			iWCIdx := SetNextWebCon(iWebPort)

			if iWCIdx == -1 {
				fmt.Println("There was an error trying to allocate a new connection!  The connection attempt was ignored!")
			} else {
				sHelloReply := "HIYA," + strconv.Itoa(webConnections[iWCIdx].iPort)

				fmt.Println("HELLO Reply: ", sHelloReply, " Sent To IP/Port: ", addrServerConversationPort.String())

				_, err := ServerConn.WriteToUDP([]byte (sHelloReply), addrServerConversationPort)

				if err != nil {
					fmt.Println(sPlace, "::WriteToUDP()", sHelloReply, "IP/Port: ", addrServerConversationPort.String(), "Error: ", err)
				}
			}

		} else if strings.Contains(sData, "CONFIRM ,") {

			sSplit := strings.Split(sData, ",")

			iPort, _ := strconv.Atoi(sSplit[1])

			fmt.Println("CONFIRM Received from Port: ", iPort)

			sIP := addrServerConversationPort.IP.String()

			SetStatusToReceiveData(sIP, iPort)
		} else {
			fmt.Println(sPlace, "::Invalid Response From WebServer. ->", sData, "<-")
		}


		time.Sleep(time.Millisecond)
	}

	wg.Done()

	return

}


func ReCycleConnection(iWebConIdx int, sReason string, err error) {

	var sPlace string = "ReCycleConnection()"

	sError := ""
	if err != nil {
		sError = "  Error: " + err.Error()
	}

	fmt.Println(sPlace, sReason, webConnections[iWebConIdx].sIP, ":", webConnections[iWebConIdx].iPort, " Index: ", iWebConIdx, sError)


	webConnections[iWebConIdx].webSvrConn.Close()				// this should ALWAYS be open, if not thrown an error!

	webConnections[iWebConIdx].iStatus = 0						// 0 = Open, 1 = Port Reserved, 2 = Running
	webConnections[iWebConIdx].tLastCheckInTime = time.Now()	// time of last "RUNNING"
	webConnections[iWebConIdx].iPort = 0						// 55058, 55059, etc...
	webConnections[iWebConIdx].iSentMessage = 0					// used to validate when a messages start flowing to the browser
	webConnections[iWebConIdx].sIP = ""							// address of Web Server
	webConnections[iWebConIdx].webSvrConn = net.UDPConn{}		// null out the web connection

}


func CheckTimeOut(iWebConIdx int) {

	// if we did receive something; but, it was not running
	dur := time.Now().Sub(webConnections[iWebConIdx].tLastCheckInTime)

	if dur.Seconds() > 60 {
		// at this point, the connection has expired, go ahead and exit out
		ReCycleConnection(iWebConIdx, "Time Out!", nil)
		return
	}

}

func CheckAndUpdateRunning(iWebConIdx int) {

	var sPlace string = "CheckAndUpdateRunning"

	var buf = make([]byte, 128)


	// 15 milli seconds should be plenty of time to check the buffer and see if something is there.
	webConnections[iWebConIdx].webSvrConn.SetReadDeadline(time.Now().Add(15 * time.Millisecond))

	nByteCount, addrServerConversationPort, err := webConnections[iWebConIdx].webSvrConn.ReadFromUDP(buf)


	if err != nil {
		bTimeout := err.(net.Error).Timeout()

		CheckTimeOut(iWebConIdx)

		if bTimeout == false {
			// we had an error on read and it was not a timeout, we need to get out
			// we leave the connection open 'thinking' that we are still sending data.
			// this will clean up on it's own
			fmt.Println(sPlace, "::ReadFromUDP()  Remote IP/Port: ", addrServerConversationPort, "WebCon IP", webConnections[iWebConIdx].sIP, "WebCon Port: ", webConnections[iWebConIdx].iPort, "Error: ", err)
			return
		} else {
			// we had an error and it was a timeout, we need to get out
			// there is nothing going on, so we just need to get out

			return
		}
	}

	// there was no error of any sort, so we need to process data...

	sData := string(buf[0:nByteCount])

	if sData == "RUNNING " {

		webConnections[iWebConIdx].tLastCheckInTime = time.Now()
		//fmt.Println(sPlace, "::Received RUNNING  Remote IP/Port: ", addrServerConversationPort, "Updating Time for Index: ", iWebConIdx)

		return
	}

	fmt.Println(sPlace, "Unknown Receive From Data Port: ", sData)

	// probably a waste to check; but, the incoming buffer 'may' not be RUNNING... so, if running is delayed
	// we still need to check it bail if the wrong stuff is coming in...
	CheckTimeOut(iWebConIdx)

}


func SendDataToWebServer(ChanBrowserRequests <-chan string) {

	var sPlace string = "SendDataToWebServer()"
	var iCount int
	var bFirstTime bool = true
	var sJSONData string


	for {

		sJSONData = <- ChanBrowserRequests

		if bFirstTime == true {
			fmt.Println(sPlace, "Processing JSON Messages...")
			bFirstTime = false
		}


		var bFoundConnection bool = false

		for iWebConIdx := range webConnections {
			iCount++

			// is the connection active, if not, go to the next item in the loop
			if webConnections[iWebConIdx].iStatus == 2 {

				bFoundConnection = true

				if webConnections[iWebConIdx].iSentMessage == 0 {
					fmt.Println(sPlace, "Sending Data To->", webConnections[iWebConIdx].sIP, ":", webConnections[iWebConIdx].iPort, "<-")
					webConnections[iWebConIdx].iSentMessage = 1
				}

				RemoteConnection := net.UDPAddr{IP: net.ParseIP(webConnections[iWebConIdx].sIP), Port: webConnections[iWebConIdx].iPort}

				_, err := webConnections[iWebConIdx].webSvrConn.WriteToUDP([]byte (sJSONData), &RemoteConnection)
				if err != nil {
					ReCycleConnection(iWebConIdx, "Write Error!", err)
				}

				CheckAndUpdateRunning(iWebConIdx)
			}
		}

		if bFoundConnection == false {
			if len(ChanBrowserRequests) > MAX_CHANNELS {
				fmt.Println(sPlace, "There are no Browser Connections.  Draining The Channel!")

				for len(ChanBrowserRequests) > 0 {
					sJSONData = <-ChanBrowserRequests
				}
			}
		}


		time.Sleep(time.Millisecond)

	}

	wg.Done()

}
