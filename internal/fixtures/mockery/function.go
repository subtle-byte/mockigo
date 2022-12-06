package mockery

import (
	"context"
)

type SendFunc func(ctx context.Context, data string) (int, error)
