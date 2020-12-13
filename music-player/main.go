package main

import (
	"fmt"
	"os"
	"time"

	"github.com/clausecker/nfc/v2"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

//Audio vars
var controller *beep.Ctrl
var stream beep.StreamSeekCloser

//NFC vars
var device nfc.Device
var readTarget *nfc.Target

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing required argument:  input mp3 file name")
		return
	}

	//Prepare audio
	err := prepareAudio(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	//Prepare NFC
	err = prepareNFC()
	if err != nil {
		panic(err)
	}
	defer device.Close()

	// Now scan - if the toggle function is on, play. Otherwise stop
	err = scanNFC(func(toggle bool) {
		if toggle {
			//Play
			speaker.Lock()
			controller.Paused = false
			speaker.Unlock()
		} else {
			speaker.Lock()
			controller.Paused = true
			speaker.Unlock()
		}
	})
	if err != nil {
		panic(err)
	}
}

func prepareAudio(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	var format beep.Format
	stream, format, err = mp3.Decode(file)
	if err != nil {
		return err
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	controller = &beep.Ctrl{
		Streamer: beep.Loop(-1, stream),
		Paused:   true,
	}

	speaker.Play(controller)

	return nil
}

func prepareNFC() error {
	var err error
	device, err = nfc.Open("pn532_uart:/dev/ttyS0")
	if err != nil {
		return err
	}

	if err = device.InitiatorInit(); err != nil {
		return err
	}

	fmt.Println("Opened Device", device, device.Connection())

	return nil
}

func scanNFC(toggleFunc func(on bool)) error {
	for {
		targets, err := device.InitiatorListPassiveTargets(nfc.Modulation{
			BaudRate: nfc.Nbr106,
			Type:     nfc.ISO14443a,
		})
		if err != nil {
			return err
		} else {
			var scanTarget *nfc.Target
			for _, target := range targets {
				if card, ok := target.(*nfc.ISO14443aTarget); ok {
					scanTarget = &target
					if readTarget == nil {
						fmt.Println("Card read", card.UID)
						readTarget = &target
						toggleFunc(true)
					}
					if err := device.InitiatorInit(); err != nil {
						return err
					}
					break
				}
			}
			if scanTarget == nil {
				if readTarget != nil {
					fmt.Println("Card removed")
					toggleFunc(false)
					readTarget = nil
				}
			}
		}
	}
}
