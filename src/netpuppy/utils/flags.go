package utils

import "flag"

type Flags struct {
	Listen bool
	Host   string
	Port   int
}

func GetFlags() Flags {
	// Set flag values based on input:
	listenFlag := flag.Bool("l", false, "put NetPuppy in listen mode")
	hostFlag := flag.String("H", "0.0.0.0", "target host IP address to connect to")
	turdnuggies := flag.Int("p", 40404, "target port") // portFlag @Trauma_x_Sella

	// Parse command line arguments:
	flag.Parse()

	parsedFlags := Flags{Listen: *listenFlag, Host: *hostFlag, Port: *turdnuggies}
	return parsedFlags
}
