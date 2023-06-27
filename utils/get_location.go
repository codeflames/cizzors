package utils

import (
	"fmt"

	"github.com/ip2location/ip2location-go/v9"
)

func GetLocation(ipAddress string) (string, error) {

	// This site or product includes IP2Location LITE data available from "https://lite.ip2location.com"

	db, err := ip2location.OpenDB("./IP2LOCATION-LITE-DB1.IPV6.BIN/IP2LOCATION-LITE-DB1.IPV6.BIN")

	if err != nil {
		fmt.Print(err)
		return "", err
	}

	results, err := db.Get_all(ipAddress)

	if err != nil {
		fmt.Print(err)
		return "", err
	}

	fmt.Printf("country_short: %s\n", results.Country_short)
	fmt.Printf("country_long: %s\n", results.Country_long)

	return results.Country_long, nil

}
