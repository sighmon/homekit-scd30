package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"

	"github.com/sighmon/homekit-scd30/promexporter"

	"github.com/pvainio/scd30"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

var acc *accessory.Thermometer
var co2 *service.CarbonDioxideSensor
var co2Level *characteristic.CarbonDioxideLevel
var humidity *service.HumiditySensor
var prometheusExporter bool
var scd30PrometheusExporter *promexporter.Exporter
var timeBetweenReadings int

func init() {
	flag.BoolVar(&prometheusExporter, "prometheusExporter", false, "Start a Prometheus exporter on port 1006")
	flag.IntVar(&timeBetweenReadings, "timeBetweenReadings", 5, "The time in seconds between CO2 readings")
	flag.Parse()
}

func readSensor() {
	// Setup the SCD30 CO2 sensor
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

			acc.TempSensor.CurrentTemperature.SetValue(float64(m.Temperature))
			humidity.CurrentRelativeHumidity.SetValue(float64(m.Humidity))
			co2Level.SetValue(float64(m.CO2))
			if m.CO2 > 850 {
				co2.CarbonDioxideDetected.SetValue(1)
			} else {
				co2.CarbonDioxideDetected.SetValue(0)
			}
			scd30PrometheusExporter.UpdateReadings(m.Temperature, m.Humidity, m.CO2)
			log.Printf("%f ppm, %f°C, %f%%", m.CO2, m.Temperature, m.Humidity)
		} else {
			log.Print("Failed to get a measurement...")
		}
	}
}

func startPrometheus() {
	scd30PrometheusExporter = promexporter.New(1006)
	scd30PrometheusExporter.Start()
}

func main() {
	// Setup the HomeKit accessory
	acc = accessory.NewTemperatureSensor(accessory.Info{
		Name:             "SCD-30",
		SerialNumber:     "ADAFRUIT-4867-SCD-30",
		Manufacturer:     "Adafruit",
		Model:            "SCD-30",
		Firmware: 		  "1.0.0",
	})

	// Add the humidity service
	humidity = service.NewHumiditySensor()
	acc.AddS(humidity.S)

	// Add the CO2 service
	co2 = service.NewCarbonDioxideSensor()
	co2Level = characteristic.NewCarbonDioxideLevel()
	co2.AddC(co2Level.C)
	acc.AddS(co2.S)

	// Store the data in the "./db" directory.
	fs := hap.NewFsStore("./db")

	// Create the hap server.
	server, err := hap.NewServer(fs, acc.A)
	if err != nil {
		// stop if an error happens
		log.Panic(err)
	}

	// Setup a listener for interrupts and SIGTERM signals
	// to stop the server.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		// Stop delivering signals.
		signal.Stop(c)
		// Cancel the context to stop the server.
		cancel()
	}()

	// Read the CO2 sensor
	go readSensor()

	// Start the Prometheus exporter
	if prometheusExporter {
		go startPrometheus()
	}

	// Run the server.
	server.ListenAndServe(ctx)
}
