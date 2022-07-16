package main

import (
	"flag"
	"log"
	"time"

	"github.com/pvainio/scd30"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

var developmentMode bool
var timeBetweenReadings int

func init() {
	flag.BoolVar(&developmentMode, "dev", false, "development mode returns a random reading")
	flag.IntVar(&timeBetweenReadings, "timeBetweenReadings", 5, "The time in seconds between CO2 readings")
	flag.Parse()
}

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	sensor, err := scd30.Open(bus)
	if err != nil {
		log.Fatal(err)
	}

	sensor.StartMeasurements(uint16(timeBetweenReadings))

	for {
		time.Sleep(time.Duration(timeBetweenReadings) * time.Second)
		hasMeasurement, err := sensor.HasMeasurement()
		if err != nil {
			log.Fatalf("error %v", err)
		}
		if hasMeasurement {
			m, err := sensor.GetMeasurement()
			if err != nil {
				log.Fatalf("error %v", err)
			}

			log.Printf("%f ppm, %f°C, %f%%", m.CO2, m.Temperature, m.Humidity)
		} else {
			log.Print("Failed to get a measurement...")
		}
	}
}
