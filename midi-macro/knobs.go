package main

import (
	"fmt"
	"os/exec"
)

// Volume controller
func updateVolume(maxValue uint8, value uint8) {
	volume := float32(value) / (float32(maxValue) / 10) // 10 is the max volume

	cmd := exec.Command("/usr/bin/osascript", "-e", "set volume "+fmt.Sprint(volume))
	cmd.Run()
}

// Brightness controller
func updateBrightness(maxValue uint8, value uint8) {
	brightness := float32(value) / (float32(maxValue) / 1) // 1 is the max brightness
	cmd := exec.Command("/usr/local/bin/brightness", fmt.Sprint(brightness))
	cmd.Run()
}

func keyPress(maxValue uint8, value uint8, prev_value uint8, key key) {
    var action string
    if value > prev_value {
        action = key.KeyUp
    } else if value < prev_value {
        action = key.KeyDown
    }

    if action != "" {
        cmd := exec.Command("/usr/bin/xdotool", "key", action)
        cmd.Run()
    }
}