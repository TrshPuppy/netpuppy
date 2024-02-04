package utils

import "fmt"

func Banner() string {
	var bannerOpening string = `
	Trash Puppy brings you...
	
	`

	var openingBanner string = `
|8PPPPe                  ___      .++.
|8    |8 |eeee |eeeee __/_, '.  .'    '. .
|8e   |8 |8      |8   \_,  | \_'  /   )'-')
|88   |8 |8eee   |8e   U ) '-'    \  (('"'
|88   |8 |88     |88   ___Y  ,    .'7 /| 
|88   |8_|88ee___|88__(_,___/___.'_(_/_/_

|8PPPPe
|8    |8 |e   .e |eeeee  |eeeee  |e   .e
|8eeee8  |8   |8 |8   |8 |8   |8 |8   |8
|88      |8e  |8 |8eee8  |8eee8  |8eee8
|88      |88  |8 |88     |88      |88
|88______|88ee8__|88_____|88______|88____

	`

	var bannerClosing string = `
         Launch a puppy to
       ~ sneef  and  fetch ~
         data   for   you!
		   `

	return fmt.Sprintf("%v%v%v\n", bannerOpening, openingBanner, bannerClosing)
}

func UserSelectionBanner(choice string, host string, portString string) string {
	var selectionBanner string

	if choice == "connect_back" {
		mode := 
		selectionBanner = `
		bork!
      __  /  
 (___()'';      |Host: {host}
 / )   /'       |Port: {port}
 /\'--/\        |Mode: {mode}
    
		`
	}
}