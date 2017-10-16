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


func ProcessBoardData(ChannelBoardData chan<- string, sPort string) (bool, string) {

	var sData string
	var sError string
	var sPlace string = "ProcessBoardData()"

	var iSkippedData int
	var iProcessedData int
	var iIgnoreTheNextMessages int

	tHold := time.Now()

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
			// this is the 4th iteration of throttling the queue... funness...
			// data is coming out much faster than can be processed
			// the browser is also having issues dealing with this.
			// ... screen updates are very slow because of the flood of JSON messages.
			// For this application, data loss is inevitable... not so with the fat clients.
			// so... I kept the inbound queue very very small... in fact, there is
			// only room for two distinct group of messages.
			// once we hit max, which is basically 12, the next set of messages are ignored.
			// this a group is ignored and you avoid the possibility of one set being continuously ignored.

			if ActiveWebConnections() == true {
				if len(ChannelBoardData) >= MAX_CHANNELS_BOARD {
					iIgnoreTheNextMessages = MAX_DISTINCT_MESSAGES
				} else {
					if len(ChannelBoardData) < MAX_CHANNELS_BOARD {
						if iIgnoreTheNextMessages == 0 {
							iProcessedData++
							ChannelBoardData <- sData
						} else {
							iSkippedData++
							iIgnoreTheNextMessages--
						}
					}
				}
			}
		}

		// send out a performance message every 15 seconds....
		tNow := time.Now()
		if tNow.Sub(tHold).Seconds() >= 15 {
			var fS float64 = float64(iSkippedData) / 15.00
			var fP float64 = float64(iProcessedData) / 15.00
			var fT float64 = float64(iProcessedData + iSkippedData) / 15.00


			sTemp := fmt.Sprintf("MESSAGE ,Skipped: %d  Skip/Sec %.2f *** Processed: %d  Processed/Sec %.2f *** Total: %d  Total/Sec: %.2f *** JSON Sent: ",
						iSkippedData, fS, iProcessedData, fP, iSkippedData + iProcessedData, fT)

			ChannelBoardData <- sTemp

			iSkippedData = 0
			iProcessedData = 0

			tHold = time.Now()
		}


		time.Sleep(time.Millisecond)
	}

	// This should never get hit.

	wg.Done()

	return true, "OK"
}
