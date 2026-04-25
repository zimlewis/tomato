package formatter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zimlewis/tomato/internal/types"
)

type WaybarFormatter struct {}

func (WaybarFormatter) Format(value types.CurrentResponse) (string, error) {

	var clock = []string{"pomodoro", "short", "long"}

	type returnBody struct {
		Text    string `json:"text"`
		Tooltip string `json:"tooltip"`
		Alt     string `json:"alt"`
		Class   string `json:"class"`
	}

	minute := value.TimeLeft / 60
	second := value.TimeLeft % 60

	resultStruct := returnBody {
		Text: fmt.Sprintf(
			"%.4s %02d:%02d",
			strings.ToUpper(clock[value.Clock]),
			minute,
			second,
		),
		Class:   "tomato",
		Tooltip: "+ Left click to start\n+ Right click to stop\n+ Scroll up to switch mod up\n+ Scroll down to switch mod down",
	}

	jsonData, err := json.Marshal(resultStruct)
	if err != nil {
		return "", fmt.Errorf("Cannot parse the response: %w\n", err)
	}



	return fmt.Sprintf("%s", string(jsonData)), nil
}
