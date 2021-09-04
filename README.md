# midi-macro

Use your MIDI controller (pads, knobs, sliders, keys etc.) to trigger macros.

## Usage

Build binary
```
go build -o midimacro midi-macro/*.go
```

Pick your device from the list of connected device with the command below.
```
./midimacro list
```

Add the device to the configuration file, and point an environment variable to it
```
export MIDI_MACRO_PATH=/path/to/midi_macros.yml
```

Run the executable binary
```
./midimacro
```
