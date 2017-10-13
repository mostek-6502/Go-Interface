Rick Faszold 

September 26th, 2016 

XEn, LLC (c) Missouri, USA 

This program taks data from a controller board, formats that data into JSON and 
sends the JSON data to a web browser, if they are connected. 

In this case, the web browser is running a Python script in SSE mode for the browser. 

This program has four parts running concurrently. 

1. Communication with Controller \n Handles a handshake with the IoT controller board. 
	Takes the incoming traffic and puts it into a Channel 
	
2. Raw Data to JSON 
	A. Takes the Board data off of the channel, IE. 
	B.1 Breaks each element up into discreet JSON messages 
	B.2 While each message is being broken up into discreet elements, 
		parts of the data are changed from Status Codes to Readbale Text 
		Ex. 
			"WINDSTAT",2,33,44 is broken up into: 
				JSON { "Element: Header", "Data: WINDSTAT"	} 
				JSON { "Element: Speed", "Data: 2"	} 
				JSON { "Element: Status", "Data: OK"	} 
				JSON { "Element: Direction", "Data: NNW"	} 
		The down stream Python script simply takes the preformaated JSON message and 
		sends it to the browser 
	C. Each JSON Message is loaded onto a JSON Channel 
3. Web Browser Connections 
	This program limits browser connections to 10, utilizing up to 20 ports from a predefined list. 
	Once a connection is made on a PREDEFINED port, a simple hand shake ensues. 
	Python Sends -> HELLO 
	This Program -> 
		Searches for an Open Port from a predefined list 
		Responds with 
		HIYA,55058 (for instance) 
		Essentially, I see you (HIAY), let's talk on 55058 from here on out 
	Python Sends -> CONFIRM ,55058 
	
	From here, all UDP communications move to 55058 
	This program keeps a running list of connected Browsers 
	The reason for the port change is to only allow hand shakes on the predefined port. 
	One port can not handle ALL of the two way traffice from multiple browsers. 
	This allowed an easier way to identify what was connected 
4. Python to Browser via JSON message 
	The final function removes data from the JSON Channel and iterates 
	through the list of connected browsers. Once an active browser is 
	discovered, the JSON message is sent to that browser. Thus, if 8 
	browsers ar on the active list, the message is sent 8 times.