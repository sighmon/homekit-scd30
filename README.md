# HomeKit SCD30 CO2 sensor

An Apple HomeKit accessory for the [Adafruit SCD30 CO2 sensor](https://www.adafruit.com/product/4867).

## Software

* Build: `go build homekit-scd30.go`
* Run: `go run homekit-scd30.go`
* In iOS Home app, click Add Accessory -> "More options..." and you should see "SCD-30"

## TODO

- [x] Read the sensor
- [x] Add HomeKit CO2, temperature, humidity accessory
- [ ] Add Prometheus exporter
