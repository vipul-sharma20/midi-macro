package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"gitlab.com/gomidi/midi"
	channel "gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/reader"
	"gopkg.in/yaml.v2"

	driver "gitlab.com/gomidi/rtmididrv"
)

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type key struct {
	Name     string   `yaml:"name"`
	Aliases  []string `yaml:"aliases,flow"`
	Type     string   `yaml:"type"`
	Task     string   `yaml:"task"`
	MaxValue uint8    `yaml:"max_value"`
	KeyUp    string   `yaml:"key_up"`
	KeyDown  string   `yaml:"key_down"`
}

type config struct {
	Port string `yaml:"port"`
	Keys []key  `yaml:"keys"`
}

var prev_values map[string]uint8

func (c config) getKey(name string) (key, error) {
	for _, key := range c.Keys {
		if key.Name == name {
			return key, nil
		}
		for _, alias := range key.Aliases {
			if alias == name {
				return key, nil
			}
		}
	}
	return key{}, fmt.Errorf("No configuration found.")
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

	// Discover which port to use.
	if len(os.Args) == 2 && os.Args[1] == "list" {
		printInPorts(ins)
		return
	}

	in, out := getIn(ins), outs[0]

	must(in.Open())
	must(out.Open())

    prev_values = make(map[string]uint8)

	rd := reader.New(
		reader.NoLogger(),

		// Fetch every message
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			switch midi_message := msg.(type) {
			case channel.NoteOn:
				messageHandler(midi_message)
			case channel.ControlChange:
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
func controlChangeHandler(midi_message channel.ControlChange) {
	conf := getConf()
	midi_key := fmt.Sprint(midi_message.Controller())
	key, err := conf.getKey(midi_key)
	if err != nil {
		fmt.Printf("Error getting key %s: %v\n", midi_key, err)
		return
	}

	prev_value, has_prev := prev_values[key.Name]
	if !has_prev { prev_value = midi_message.Value() }
	prev_values[key.Name] = midi_message.Value()

	if key.Task == "volume" {
		updateVolume(key.MaxValue, midi_message.Value())
	} else if key.Task == "brightness" {
		updateBrightness(key.MaxValue, midi_message.Value())
	} else if key.Task == "keyPress" {
		keyPress(key.MaxValue, midi_message.Value(), prev_value, key);
	} else {
		fmt.Printf("Unknown task: %v\n", key.Task)
	}
}

// Button events
func messageHandler(midi_message channel.NoteOn) {
	conf := getConf()
	midi_key := fmt.Sprint(midi_message.Key())
	key, err := conf.getKey(midi_key)
	if err != nil {
		fmt.Printf("Error getting key %s: %v\n", midi_key, err)
		return
	}

	command := getCommand(key.Task)
	name, args := command[0], command[1:]

	cmd := exec.Command(name, args...)
	cmd.Run()
}

// YAML handler
func getConf() *config {
	var conf config
	yamlFile, err := ioutil.ReadFile(os.Getenv("MIDI_MACRO_PATH"))
	if err != nil {
		fmt.Printf("Cannot Read file '%s': #%v ", yamlFile, err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}

	return &conf
}

func getIn(ports []midi.In) (def midi.In) {
	conf := getConf()
	portIndex, err := strconv.Atoi(conf.Port)

	if err == nil {
		return ports[portIndex]
	}

	for _, in := range ports {
		if strings.Contains(in.String(), conf.Port) {
			return in
		}
	}

	fmt.Printf("Failed to find port: %s\n", conf.Port)
	os.Exit(1)

	return
}

// Get task command
func getCommand(task string) []string {
	split := strings.Split(task, ",")
	return split
}

func printPort(port midi.Port) {
	fmt.Printf("[%v] %s\n", port.Number(), port.String())
}

func printInPorts(ports []midi.In) {
	fmt.Printf("MIDI IN Ports\n")
	for _, port := range ports {
		printPort(port)
	}
	fmt.Printf("\n\n")
}
