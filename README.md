# Navigation Position

**NavPos** is a simple and lightweight application that displays the real-time position of a MAVLink-enabled aircraft on an interactive map, alongside nearby aircraft detected using the [OpenSky Network](https://opensky-network.org/). It uses open standards and formats, making it easy to run both online and offline.

![navpos](docs/navpos-main.png?raw=true "NavPos")

## Features

- Real-time telemetry from MAVLink-compatible systems
- Aircraft visualization with heading and location
- Nearby aircraft displayed using data from the OpenSky Network
- Offline map support with vector tiles in `.mbtiles` format
- Clean and simple web-based UI
- Built using Go for the backend and Leaflet for the frontend

## Getting Started

### Run in simulation mode for testing purposes:

```bash
./navpos start --host=localhost --port=8080 --tsprovider=fake --tprovider=fake --aprovider=fake
```

### Run with a live telemetry source (e.g. via USB):

```bash
./navpos start --host=localhost --port=8080 --tparams=device=/dev/ttyUSB0,baudrate=115200
```

## Offline Maps

NavPos supports offline maps using `.mbtiles` vector tile files. You can download OpenStreetMap-compatible offline map data from [MapTiler](https://www.maptiler.com/data/).

## License

This project is open-source and uses only open standards and freely available map data.

---

NavPos is ideal for UAV enthusiasts, developers, or researchers who want to visualize telemetry data in real time with minimal setup.