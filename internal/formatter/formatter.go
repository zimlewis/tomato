package formatter

import "github.com/zimlewis/tomato/internal/types"

// A formatter that is used by the output the format the response to different kind of string
type Formatter interface {
	Format(data types.CurrentResponse) (string, error)
}

// Take a string and return a formatter that matched the string, return DefaultFormatter by default
func NewFromString(s string) Formatter {
	switch s {
	case "waybar": return WaybarFormatter{}
	case "basic": return BasicFormatter{}
	case "default": return DefaultFormatter{}
	}
	return DefaultFormatter{}
}

// Basic structure of a formatter
// type f struct {}
//
// func (f) Format(value types.CurrentResponse) (string, error) {
// 	var clocks = []string{"pomodoro", "short", "long"}
//
// 	minute := value.TimeLeft / 60
// 	second := value.TimeLeft % 60
// 	clock := clocks[value.Clock]
//
// 	return fmt.Sprintf("%s %02d %02d", clock, minute, second), nil
// }
