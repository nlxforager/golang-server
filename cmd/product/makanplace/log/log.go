package log

import (
	"context"
	"fmt"
)

func HttpRequestPrefix(ctx context.Context) string {
	return fmt.Sprintf("%s %s %s", ctx.Value("METHOD"), ctx.Value("URL"), ctx.Value("ORIGIN"))
}
