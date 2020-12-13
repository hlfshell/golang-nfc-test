package main

import (
	"fmt"

	nfc "github.com/clausecker/nfc/v2"
)

func init() {

}

func main() {
	dev, err := nfc.Open("pn532_uart:/dev/ttyS0")
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	if err := dev.InitiatorInit(); err != nil {
		panic(err)
	}

	fmt.Println("Opened Device", dev, dev.Connection())

	out, err := dev.Information()
	fmt.Println(out)
	if err != nil {
		panic(err)
	}

	var readTarget *nfc.Target

	for {
		targets, err := dev.InitiatorListPassiveTargets(nfc.Modulation{
			BaudRate: nfc.Nbr106,
			Type:     nfc.ISO14443a,
		})
		if err != nil {
			panic(err)
		} else {
			var scanTarget *nfc.Target
			for _, target := range targets {
				if card, ok := target.(*nfc.ISO14443aTarget); ok {
					scanTarget = &target
					if readTarget == nil {
						fmt.Println("Card read", card.UID)
						readTarget = &target
					}
					if err := dev.InitiatorInit(); err != nil {
						panic(err)
					}
					break
				}
			}
			if scanTarget == nil {
				if readTarget != nil {
					fmt.Println("Card removed")
					readTarget = nil
				}
			}
		}
	}
}
