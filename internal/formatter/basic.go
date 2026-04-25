package formatter

import (
	"fmt"

	"github.com/zimlewis/tomato/internal/types"
)

type BasicFormatter struct {}

func (BasicFormatter) Format(value types.CurrentResponse) (string, error) {
	var clocks = []string{"pomodoro", "short", "long"}

	minute := value.TimeLeft / 60
	second := value.TimeLeft % 60
	clock := clocks[value.Clock]

	return fmt.Sprintf("%s %02d %02d", clock, minute, second), nil
}

