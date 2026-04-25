package formatter

import "github.com/zimlewis/tomato/internal/types"

type Formatter interface {
	Format(data types.CurrentResponse) (string, error)
}

func NewFromString(s string) Formatter {
	switch s {
	case "waybar": return WaybarFormatter{}
	case "default": return Default{}
	}
	return Default{}
}
