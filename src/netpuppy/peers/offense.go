package peers

type Offense struct {
	/*
		listenTo string = userInput

		deliverTo string = socket

		RedirectUserInputToSocket()

		RedirectSocketOutputTo User()

		userinput ==> socket
		socket ==> print to user

	*/
}

/*
	userChannel = make(channel)
	socketChannel = make(channel)

	go putuserinputIntochannle(userChannel channel)
	go putSocketOutputIntoChannle(socketChannel channel)

	func putUserInputIntochannel(userChannel chan<-)
		userinput = input()
		userChannel <- userInput

	func putSocketOutputIntoChannel(socketChannel chan<-)
		socketdata = readSocket()
		socketChannel <- socketData

	forever/select:
	case userChannelData := userChannel<-
		socket.send(userChannelData)
	case socketChannelData := socketChannel <-
		printToUser(socketChannelData)



*/
