package formatter

import "github.com/zimlewis/tomato/internal/types"

type Formatter interface {
	Format(data types.CurrentResponse) (string, error)
}
