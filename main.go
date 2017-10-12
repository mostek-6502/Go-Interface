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
// This program has four parts running concurrently.
// 1. Communication with Controller
//    Handles a handshake with the IoT controller board.
//    Takes the incoming traffic and puts it into a Channel
// 2. Raw Data to JSON
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
// 3. Web Browser Connections
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
//        One port can not handle ALL of the two way traffice from multiple browsers.
//        This allowed an easier way to identify what was connected
// 4. Python to Browser via JSON message
//    The final function removes data from the JSON Channel and iterates
//    through the list of connected browsers.  Once an active browser is
//    discovered, the JSON message is sent to that browser.  Thus, if 8
//    browsers ar on the active list, the message is sent 8 times.
//
// Current Issues
// 1. Color Needs to be added to the JSON message.. IE. Statuses of Green = OK, Yellow = Caution, Red = Alert
// 2. Add Pass through functionality from the Browser to the Board (and back) for configuration
//    purposes.  For instance, allow the browser to update the Temperature resolution from 9 to 11 (for instance)
//    Report from the Board a 'pending' update to the browser.  The browser can either finalize the update
//    or allow the board to be rebooted and reset back to original settings
// 3. Conplete the minor status conversions for eacn message.  IE. Convert data elements such as "0" to "OK"
// 4. If all of the ports are in use, do we
// 		a. Open a connection to the browser
//		b. send a JSON alert message and close the port out?
//		Since this is a specialy application, the need for more than 20 open browsers is small
//		This is a low priority item.
// 5. Validate / Verify End to End Data Loss
//		a. board - indicates how much data is being sent
//			the high level board data should be normalized to individual JSON message counts
//		b. intbb - indicates how many messages are being sent to python
//			i. track channel insertion JSON messages need to be created along with
//			ii. track channel delete messages and send JSON accordingly
//		c. a & b need to be measured over time and load -> for instance how many browser sessions are active
//		d. Python script needs to forumate it's own JSON message indicating sent message
//		e. accumulate all in the browser.
//		Since retining every piece of data is not critical this is a lower priority item.


package main

import (
	"os"
	"sync"
	"fmt"
)


const MAX_CHANNELS int = 2000

var wg sync.WaitGroup

func main() {

	fmt.Println("main() begin...")


	wg.Add(4)
	//kaprekar()

	ChanChannelData := make(chan string, MAX_CHANNELS)
	// ChanBrowserConnects := make(chan string)
	ChanBrowserRequests := make(chan string, MAX_CHANNELS)

	sPort := "55056"

	// these small conditionals allow for and easy on/off for any given process

	// this 'thread' maintains the list of connected browsers and their ports.
	sJunk := "0"
	if sJunk == "0" { go ManageWebServerConnections(sPort) }

	// this 'thread' connects to the board and inserts the data into a queue
	sJunk = "1"
	if sJunk == "1" { go ProcessBoardData(ChanChannelData, sPort) }

	// this takes the channel data, formats it to JSON and maintains a data queue for browser requests
	sJunk = "2"
	if sJunk == "2" { go FormatAndQueueData(ChanChannelData, ChanBrowserRequests) }

	// although we send the sPort for the Baord listen, by default, we'll add 1 to the Port number for the web data
	sJunk = "3"
	if sJunk == "3" { go SendDataToWebServer(ChanBrowserRequests) }

	wg.Wait()
	// handle browser requests
	//ProcessBrowserRequests(ChanBrowserRequests)

	// ok, so the next go routine should convert the data to a JSON type

	os.Exit(1)

}
