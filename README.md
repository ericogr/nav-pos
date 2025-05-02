# Navigation Position
A monitoring tool designed to help you detect nearby aircraft during your flight.

![navpos](docs/navpos-main.png?raw=true "NavPos")

## Run
Run in simulation mode for testing purposes:

```bash
./navpos start --host=localhost --port=8080 --tprovider=fake --aprovider=fake
```

Execute with the live API:

```bash
./navpos start --host=localhost --port=8080 --tparams=device=/dev/ttyUSB0,baudrate=115200
```