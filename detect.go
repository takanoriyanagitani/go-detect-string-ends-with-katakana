package detect

import (
	"context"
)

type Last3 [9]byte

type EndsWithKatakana3 func(context.Context, Last3) (detected bool, err error)
