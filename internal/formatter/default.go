package formatter

import (
	"fmt"

	"github.com/zimlewis/tomato/internal/types"
)

// DefaultFormatter will return this form: Your clock session has min:sec remaining
// eg: Your short session has 10:20 remaining
type DefaultFormatter struct {}

func (DefaultFormatter) Format(value types.CurrentResponse) (string, error) {
	var clocks = []string{"pomodoro", "short", "long"}

	minute := value.TimeLeft / 60
	second := value.TimeLeft % 60
	clock := clocks[value.Clock]

	return fmt.Sprintf("Your %s session has %02d:%02d remaining", clock, minute, second), nil
}
