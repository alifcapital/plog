package log

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestCycle(t *testing.T) {
	l := New(os.Stdout, false)
	ctx := context.Background()
	l.Info(ctx, "test", P("param", "value"))
	l.Debug(ctx, "test", P("param", "value"))
	l.Err(ctx, "test", errors.New("some err"), P("param", "value"))
	l.Warn(ctx, "test", P("param", "value"))
}
