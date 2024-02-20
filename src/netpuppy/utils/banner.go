package utils

import "fmt"

// Build and return main banner:
func Banner() string {

	var openingBanner string = `
	Trash Puppy brings you...

|8PPPPe
|8    |8 |eeee |eeeee    ___      .++.
|8e   |8 |8      |8   __/_, '.  .'    '. .
|88   |8 |8eee   |8e  \_,  | \_'  /   )'-')
|88   |8 |88     |88   U ) '-'    \  (('"'
|88   |8 |88ee   |88   ___Y  ,    .'7 /|
______________________(_,___/___.' (_/_/_
|8PPPPe
|8    |8 |e   .e |eeeee  |eeeee  |e   .e
|8eeee8  |8   |8 |8   |8 |8   |8 |8   |8
|88      |8e  |8 |8eee8  |8eee8  |8eee8
|88      |88  |8 |88     |88      |88
|88      |88ee8  |88     |88      |88
________________________________________
	
	  Launch a puppy to
   	~ sneef  and  fetch ~
	  data   for   you!
	  
	`
	return openingBanner
}

// Build a banner and return based on the type of peer the user started:
func UserSelectionBanner(choice string, host string, remotePort int, localPort int) string {
	var selectionBanner string
	var s0 string
	var s1 string
	var s2 string
	var s3 string
	var s4 string

	if choice == "connect_back" {
		mode := "Connect Back"
		s0 = `
	bork!
`
		s1 = fmt.Sprintf("     __  /\n")
		s2 = fmt.Sprintf("(___()'';    |Mode:  %v\n", mode)
		s3 = fmt.Sprintf("/ )   /'     |Host:  %v\n", host)
		s4 = fmt.Sprintf("/\\'--/\\      |RPort: %v", remotePort)
	} else {
		mode := "Offensive Server"
		s0 = `
    *sneef sneef*
   .-.
`
		s1 = fmt.Sprintf("  / (_\n")
		s2 = fmt.Sprintf(" ( \"  6\\___o   |Mode:  %v\n", mode)
		s3 = fmt.Sprintf(" /  (  ___/    |Host:  %v\n", host)
		s4 = fmt.Sprintf("/     /  U     |LPort: %v", localPort)
	}

	selectionBanner = fmt.Sprintf("%v%v%v%v%v", s0, s1, s2, s3, s4)
	return selectionBanner
}

func PrintMissingPortToBanner(peerType string, LPort int, RPort int) string {
	var missingBannerPiece string

	// If the current peer is the client, we need to add the local port to banner:
	if peerType == "connect_back" {
		missingBannerPiece = fmt.Sprintf("             |LPort: %v\n", LPort)
	} else { // If the current peer is the server, we need to add the remote port to banner:
		missingBannerPiece = fmt.Sprintf("               |RPort: %v\n", RPort)
	}

	return missingBannerPiece
}

//func worker(done chan bool) {
//	fmt.Print("working...")
//	time.Sleep(time.Second * 5)
//	fmt.Println("done")
//
//	done <- true
//	done <- true
//	done <- true
//	fmt.Println("after done")
//}

//func main() {
//	done := make(chan bool, 1)
//	done <- true
//	go worker(done)
//	//	<-done
//	fmt.Printf("tiddies\n")
//	time.Sleep(time.Second * 3)
//	fmt.Printf("tiddies 2\n")
//}

/* TO DO:
- make channels non blocking of each other?
	- move channel creation out of for loop?
	- do the channels need to be destroyed?
- fix prompt


*/
