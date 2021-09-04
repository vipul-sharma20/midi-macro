`midi_macros.yml` contains a minimal sample of macro triggers and mappings. You
can customize based on your preference of keys and triggers.

YAML structure below:

```yaml
# Put your device here (use "midimacro list" command to get device)
port: "Impact LX61+ MIDI1"

keys:
  # MIDI controller key name
  - name: "12"

  # One of button, knob
  - type: "button"

  # Action you want to map to this MIDI component
  - task: "/usr/bin/osascript,/path/to/task.scpt" 
```
