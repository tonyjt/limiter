package limiter

import "context"

//需求

type Limiter interface {
	Limit(ctx context.Context, key string)
}
