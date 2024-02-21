package peers

type ConnectBack struct {
	/*
		listenTo string = shell(stdout)

		deliverTo = socket

		RedirectShellOutToSocket()

		RedirectSocketOutToShellIn()

		shellstdout ==> socket
		socket ==> shellstdin
	*/
}

/*
	shellChannel = make(channel)
	socketChannel = make(channel)

	go putSocketOutputIntoChannel(socketChannel)
	go putShellOutputIntoChannel(shellChannel)

	func putSocketOutputIntoChannel(socketChannel chan<-)
		forever
			socketData = readSocket()
			socketchannel <- socketData

	func putShellOutputIntoChannel(shellchannel chan <-)
		forever
			shellData = shell.stdout() OR shell.stderr()
			shellChannel <- shellData

	foreveer/select:
	case dataInSocketChan := socketChannel <-
		shell.stdin.send(dataInSocketChan)
	case dataInShellChan := shellChannel <-
		socket.send(dataInShellChan)
*/
