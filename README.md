midi-macro
==========

Use your MIDI controller (pads, knobs, sliders, keys etc.) to trigger macros.

To build:

`go build -o midimacro midi-macro/*.go`

Then:
`./midimacro list`

And pick your device from the list. Edit its name to the configuration file, and point an environment variable to it: 

`export MIDI_MACRO_PATH=/path/to/midi_macros.yml`

And

`./midimacro`
