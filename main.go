// Rick Faszold
//
// September 26th, 2016
//
// XEn, LLC (c)   Missouri, USA
//
// This program taks data from a controller board, formats that data into JSON and
// sends the JSON data to a web browser, if they are connected.
//
// In this case, the web browser is running a Python script in SSE mode for the browser.
//
// This program has four parts that run concurrently.
// 1. Communication with Controller
//    Handles a handshake with the IoT controller board.
//    Takes the incoming traffic and puts it into a Channel
// 2. Raw Data to JSON
//    A. Takes the Board data off of the channel,  IE.
//    B.1 Breaks each element up into discreet JSON messages
//    B.2 While each message is being broken up into discreet elements,
//        parts of the data are changed from Status Codes to Readable Text
//       Ex.
//         "WINDSTAT",2,33,44  is broken up into:
//             JSON { "Element: Header", "Data: WINDSTAT"	}
//             JSON { "Element: Speed", "Data: 2"			}
//             JSON { "Element: Status", "Data: OK"			}
//             JSON { "Element: Direction", "Data: NNW"		}
//          The down stream Python script simply takes the preformatted JSON message and
//          sends it to the browser
//    C. From here, each JSON Message is loaded onto a JSON Channel
// 3. Web Browser Connections Initiated by Python
//    This program limits browser connections to 10, utilizing up to 20 ports from a predefined list.
//    Once a connection is made on a PREDEFINED port, a simple hand shake ensues.
//    Python Sends -> HELLO
//    This Program ->
//        Searches for an Open Port from a predefined list
//        Responds with HIYA,55058 (for instance)
//           Essentially, I see you (HIAY), let's talk on 55058 from here on out
//    Python Sends -> CONFIRM ,55058
//
//    From here, all UDP communications move to 55058
//    This program keeps a running list of connected Browsers
//    The reason for the port change is to only allow hand shakes on the predefined port.
//    When multiple browsers are connected, it'll be too much data for one port to handle.
//        One port can not handle ALL of the two way traffice from multiple browsers.
//        This allowed an easier way to identify what was connected
// 4. Channel to Python to Browser via JSON message
//    The final function removes data from the JSON Channel and iterates
//    through the list of connected browsers.  Once an active browser is
//    discovered, the JSON message is sent to that browser.  Thus, if 8
//    browsers are on the active list, the message is sent 8 times.
//
// Current Issues
// 1. Color Needs to be added to the JSON message.. IE. Statuses of Green = OK, Yellow = Caution, Red = Alert
// 2. The current color scheme is pre-beta
// 3. Add Pass through functionality from the Browser to the Board (and back) for configuration
//    purposes.  For instance, allow the browser to update the Temperature resolution from 9 to 11 (for instance)
//    Report from the Board a 'pending' update to the browser.  The browser can either finalize the update
//    or allow the board to be rebooted and reset back to original settings
// $. If all of the ports are in use, do we
// 		a. Open a connection to the browser
//		b. send a JSON alert message and close the port out?
//		Since this is a specialy application, the need for more than 20 open browsers is small
//		This is a low priority item.
// 5. Validate / Verify End to End Data Loss
//      This is somewhat problematc because the queue depth is really really tight.
//      With this in mind, data loss is inevitable.
//      The other factor in this related to board queue depth is browser responsiveness
//      The browser can easily bog down on all of the data being sent to it.
//      The ONLY message to be added is one from the Python script adding to the current MESSAGE.
//      Although the message looks sloppy, everything is in one place.


package main

import (
	"os"
	"sync"
	"fmt"
)


const MAX_DISTINCT_MESSAGES = 6		// there is one particular message added every 30 seconds.  Technicall, this is 7.
const MAX_CHANNELS_BOARD int = 18
const MAX_CHANNELS_BROWSER int = 500  // this is essentially 25x's bigger than the BOARD queue due to message expansion

var wg sync.WaitGroup

func main() {

	fmt.Println("main() begin...")

	wg.Add(4)

	// the input channel creates a LOT of data for the output channel
	ChannelBoardData := make(chan string, MAX_CHANNELS_BOARD)
	// there are six distinct messages that come out of the board.
	// on average there are 22 JSON messages created for each board message
	ChannelBrowserData := make(chan string, MAX_CHANNELS_BROWSER)

	sPort := "55056"

	// this 'thread' maintains the list of connected browsers and their ports.
	go ManageWebServerConnections(sPort)

	// this 'thread' connects to the board and inserts the data into a queue
	go ProcessBoardData(ChannelBoardData, sPort)

	// this takes the channel data, formats it to JSON and maintains a data queue for browser requests
	go FormatAndQueueData(ChannelBoardData, ChannelBrowserData)

	// although we send the sPort for the Baord listen, by default, we'll add 1 to the Port number for the web data
	go SendDataToWebServer(ChannelBrowserData)

	wg.Wait()

	os.Exit(1)

}
