package main

import (
	"fmt"

	nfc "github.com/clausecker/nfc/v2"
)

func init() {

}

func main() {
	dev, err := nfc.Open("")
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	if err := dev.InitiatorInit(); err != nil {
		panic(err)
	}

	fmt.Println("Opened Device", dev, dev.Connection())

	for {
		targets, err := dev.InitiatorListPassiveTargets(nfc.Modulation{
			BaudRate: nfc.Nbr106,
			Type:     nfc.ISO14443a,
		})
		if err != nil {
			panic(err)
		} else {
			for _, target := range targets {
				if card, ok := target.(*nfc.ISO14443aTarget); ok {
					fmt.Println("Card read", card.UID)
				}
			}
		}
	}
}
