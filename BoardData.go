// Rick Faszold
//
// September 26th, 2016
//
// XEn, LLC (c)   Missouri, USA
// 		Communication with Controller
//    	Handles a handshake with the IoT controller board.
//		Receive SYNCH_2
//		Send SYNCH_3
//    	Takes the incoming traffic and puts it into a Channel

package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)


func ProcessBoardData(ChanChannelData chan<- string, sPort string) (bool, string) {

	var sData string
	var sError string
	var sPlace string = "ProcessBoardData()"

	fmt.Println(sPlace, " Calling ResolveUDPAddr()  Using Port ", sPort)

	ServerAddr, err := net.ResolveUDPAddr("udp", ":" + sPort)
	if err != nil {
		fmt.Println(sPlace, "::ResolveUDPAddr()  Port: ", sPort, "Error: ", err)
		return false, sError
	}

	fmt.Println(sPlace, " Calling ListenUDP()")

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		fmt.Println(sPlace, "::ListenUDP()  Port: ", sPort, "Error: ", err)
		return false, sError
	}

	defer ServerConn.Close()

	fmt.Println(sPlace, " Wait For In Coming Board Data")

	var bReceivedData = false

	var buf = make([]byte, 1024)
	for {

		nByteCount, addr, err := ServerConn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println(sPlace, "::ReadFromUDP()  Port: ", sPort, "Error: ", err)
			return false, sError
		}

		sData = string(buf[0:nByteCount])
		//fmt.Println("Received From Board ->", sData, "<- from ", addr, "ByteCount: ", nByteCount)

		if bReceivedData == false {
			fmt.Println(sPlace, "Receiving Data From Board!")
			bReceivedData = true
		}

		if strings.Contains(sData, "SYNCH_2") {

			fmt.Println(sPlace, "Sending ->SYNCH_3 <- to ", addr)
			_, err := ServerConn.WriteToUDP([]byte ("SYNCH_3 "), addr)

			if err != nil {
				fmt.Println(sPlace, "::WriteToUDP()::SYNCH_3  Port: ", sPort, "Error: ", err)
				return false, sError
			}

		} else {
			// The problem with this is this...  once we go into a 'full' state,
			// we can can easily begin dropping incoming board data.  This is an issue.
			// However, if we do not do this, we will the fill channel and ignore data anyway
			// By doing it this way, a full message is accepted.  Various whole messages
			// may get dropped by the O/S.  Since the data is not critical, some data loss
			// is expected.

			for len(ChanChannelData) >= MAX_CHANNELS {
				time.Sleep(time.Millisecond)
			}

			ChanChannelData <- sData

			//fmt.Println(sPlace, " ***After Board Insert: ", len(ChanChannelData))
		}

		time.Sleep(time.Millisecond)
	}

	// This should never get hit.

	wg.Done()

	return true, "OK"
}
