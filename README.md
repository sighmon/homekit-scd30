# HomeKit SCD-30 CO2 sensor

An Apple HomeKit accessory for the [Adafruit SCD-30 CO2 sensor](https://www.adafruit.com/product/4867).

![The accessory added to iOS](homekit-scd30.jpg)

## Hardware

* Adafruit SCD-30 CO2 sensor ([Core Electronics](https://core-electronics.com.au/adafruit-scd-30-ndir-co2-temperature-and-humidity-sensor-stemma-qt-qwiic.html))
* Raspberry Pi (3, 4, zero)

### Wiring

Adafruit have a great tutorial: [Learn SCD-30](https://learn.adafruit.com/adafruit-scd30/python-circuitpython)

| Raspberry Pi pin | SCD-30 sensor pin | Wire colour |
| - | - | - |
| `1` 3.3V | `VIN` Voltage in | Red |
| `3` GPIO 02 I2C SDA | `SDA` I2C SDA | Blue |
| `5` GPIO 03 I2C SCL | `SCL` I2C SCL | Yellow |
| `6` Ground | `GND` Ground | Black |

## Software

* Install [Go](http://golang.org/doc/install) >= 1.14 ([useful Gist](https://gist.github.com/pcgeek86/0206d688e6760fe4504ba405024e887c) for Raspberry Pi)
* Build: `go build homekit-scd30.go`
* Run: `go run homekit-scd30.go`
* In iOS Home app, click Add Accessory -> "More options..." and you should see "SCD-30"

### Prometheus exporter

To export the `co2`, `temperature`, and `humidity` for [Prometheus](https://prometheus.io) use the optional flag `-prometheusExporter`.

* Run: `go run homekit-scd30.go -prometheusExporter`

You'll then see the data on port `8000`: http://localhost:8000/metrics

```
# HELP co2 CO2 measured (ppm)
# TYPE co2 gauge
co2 513.689697265625

# HELP temperature Temperature measured (°C)
# TYPE temperature gauge
temperature 16.708629608154297

# HELP humidity Relative humidity measured (%)
# TYPE humidity gauge
humidity 66.6168212890625
```

## TODO

- [x] Read the sensor
- [x] Add HomeKit CO2, temperature, humidity accessory
- [x] Add Prometheus exporter
- [ ] Add Pull Request for sending a [forced calibration reference](https://learn.adafruit.com/adafruit-scd30/field-calibration), and maybe [altitude calibration](https://github.com/adafruit/Adafruit_CircuitPython_SCD30/blob/5566cb8133541c1c211d3c0c0430524d2890d71a/adafruit_scd30.py#L172)
