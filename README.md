# midi-macro

Use your MIDI controller (pads, knobs, sliders, keys etc.) to trigger macros.

## Usage

Build binary
```
go build -o $GOPATH/bin/midimacro midi-macro/*.go
```

Usage doc for `midimacro` command
```
Tool to map macros to your MIDI controller

Usage:
  midimacro [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  list        List of connected devices
  run         Run MIDI event listener

Flags:
  -h, --help   help for midimacro

Use "midimacro [command] --help" for more information about a command.
```

Pick your device from the list of connected device with the command below.
```
midimacro list
```

Add the device to the configuration file, and point an environment variable to it
```
export MIDI_MACRO_PATH=/path/to/midi_macros.yml
```

Run the executable binary
```
midimacro run
```
