package utils

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ip2location/ip2location-go/v9"
)

func GetLocation(c *fiber.Ctx) (string, string, error) {

	ipAddress := ""

	// Check if the X-Forwarded-For header exists
	if xff := c.Get(fiber.HeaderXForwardedFor); xff != "" {
		// Split the header value to get the client IP address
		ip := strings.Split(xff, ",")[0]
		fmt.Printf("Client IP: %s\n", ip)
	} else {
		// X-Forwarded-For header doesn't exist, use RemoteIP
		ipAddress = c.IP()
		fmt.Printf("Client IP: %s\n", ipAddress)
	}

	// This site or product includes IP2Location LITE data available from "https://lite.ip2location.com"

	db, err := ip2location.OpenDB("./IP2LOCATION-LITE-DB1.IPV6.BIN/IP2LOCATION-LITE-DB1.IPV6.BIN")

	if err != nil {
		fmt.Print(err)
		return "", "", err
	}

	results, err := db.Get_all(ipAddress)

	if err != nil {
		fmt.Print(err)
		return "", "", err
	}

	fmt.Printf("country_short: %s\n", results.Country_short)
	fmt.Printf("country_long: %s\n", results.Country_long)

	return results.Country_long, ipAddress, nil

}
