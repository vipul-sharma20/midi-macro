package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"gitlab.com/gomidi/midi"
	. "gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/reader"
	"gopkg.in/yaml.v2"

	driver "gitlab.com/gomidi/rtmididrv"
)

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type config struct {
	Keys []struct {
		Name     string `yaml:"name"`
		Type     string `yaml:"type"`
		Task     string `yaml:"task"`
		MaxValue uint8  `yaml:"max_value"`
	} `yaml:"keys"`
}

func main() {
	drv, err := driver.New()
	must(err)

	// make sure to close all open ports at the end
	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	outs, err := drv.Outs()
	must(err)

	in, out := ins[0], outs[0]

	must(in.Open())
	must(out.Open())

	rd := reader.New(
		reader.NoLogger(),

		// Fetch every message
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			// inspect
			switch midi_message := msg.(type) {
			case NoteOn:
				messageHandler(midi_message)
			case ControlChange:
				controlChangeHandler(midi_message)
			}
		}),
	)

	exit := make(chan string)

	// listen for MIDI
	go rd.ListenTo(in)

	for {
		select {
		case <-exit:
			os.Exit(0)
		}
	}

}

// Knobs/Slider events
func controlChangeHandler(midi_message ControlChange) {
	var conf config
	conf.getConf()

	for _, controllerConf := range conf.Keys {
		midi_knob := fmt.Sprint(midi_message.Controller())
		if controllerConf.Name == midi_knob {
			if controllerConf.Task == "volume" {
				updateVolume(controllerConf.MaxValue, midi_message.Value())
			} else if controllerConf.Task == "brightness" {
				updateBrightness(controllerConf.MaxValue, midi_message.Value())
			}
		}
	}
}

// Button events
func messageHandler(midi_message NoteOn) {
	var conf config
	conf.getConf()

	for _, v := range conf.Keys {
		midi_key := fmt.Sprint(midi_message.Key())
		if v.Name == midi_key {
			command := getCommand(v.Task)
			name, args := command[0], command[1:]

			cmd := exec.Command(name, args...)
			cmd.Run()
		}
	}

}

// YAML handler
func (conf *config) getConf() *config {
	yamlFile, err := ioutil.ReadFile(os.Getenv("MIDI_MACRO_PATH"))
	if err != nil {
		fmt.Printf("Cannot Read file   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}

	return conf
}

// Get task command
func getCommand(task string) []string {
	split := strings.Split(task, ",")
	return split
}
