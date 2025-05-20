package log

import (
	"context"
	"fmt"
)

func HttpRequestPrefix(ctx context.Context) string {
	return fmt.Sprintf("METHOD:%s URL:%s ORIGIN:%s UA: %s", ctx.Value("METHOD"), ctx.Value("URL"), ctx.Value("ORIGIN"), ctx.Value("USER-AGENT"))
}
