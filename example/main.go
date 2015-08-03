package main

import (
	"fmt"

	"github.com/scarlson/Locate"
)

func main() {
	location, err := locate.WhereAmI()
	if err != nil {
		panic(err)
	}
	fmt.Println("Latitude:", location.Latitude, "Longitude:", location.Longitude)
}
