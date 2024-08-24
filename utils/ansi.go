package utils

func Trie() string {
	return "tiddies"
}

// func RunSTDIN(c conn.ConnectionGetter) {
// 	// set up the pipe thing with io.Copy
// 	logfile, _ := os.OpenFile("teststdin", os.O_RDWR, 0777)
// 	fmt.Printf("logfile: %s\n", logfile.Name())

// 	reader := bufio.NewReader(os.Stdin)

// 	//x, err := reader.ReadBytes('\n')

// 	x, err := reader.ReadByte()
// 	// x, err := io.Copy(logfile, reader)
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 	}
// 	fmt.Printf("Read: %d\n", x)

// 	if x == byte(27) {
// 		fmt.Printf("tiddies, its a control sequence")
// 	}

// 	b := make([]byte, 1)
// 	b = append(b, x)
// 	os.WriteFile(logfile.Name(), b, 0777)

// 	return
// }
