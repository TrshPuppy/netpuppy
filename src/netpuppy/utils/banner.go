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
func UserSelectionBanner(choice string, host string, port int) string {
	var selectionBanner string

	if choice == "connect_back" {
		mode := "Client"
		string1 := `
	bork!
     __  /  	
`
		string2 := fmt.Sprintf("(___()'';      |Host: %v\n", host)
		string3 := fmt.Sprintf("/ )   /'       |Port: %v\n", port)
		string4 := fmt.Sprintf("/\\'--/\\        |Mode: %v\n", mode)

		selectionBanner = fmt.Sprintf("%v%v%v%v", string1, string2, string3, string4)
	} else {
		mode := "Offensive Server"
		s1 := `
   .-.  *sneef sneef*
  / (_
`
		s2 := fmt.Sprintf(" ( \"  6\\___o    |Host: %v\n", host)
		s3 := fmt.Sprintf(" /  (  ___/     |Port: %v\n", port)
		s4 := fmt.Sprintf("/     /  U      |Mode: %v\n", mode)

		selectionBanner = fmt.Sprintf("%v%v%v%v", s1, s2, s3, s4)
	}
	return selectionBanner
}

// .-.  *sneef sneef*
// / (_
// ( "  6\\___o   |Host: {host}
// /  (  ___/    |Port: {port}
// /     /  U     |Mode: {mode}
