package finch

import (
	"context"
	"log/slog"
)

type ContextKey string

type Context interface {
	Context() context.Context
	Screen() *Screen
	Time() *Time
	Logger() *slog.Logger
	SetLogger(logger *slog.Logger) Context
	Get(key ContextKey) any
	Set(key ContextKey, value any) Context
}

func NewContext(ctx context.Context, logger *slog.Logger, screen *Screen, time *Time) Context {
	return &finchCtx{
		ctx:    ctx,
		logger: logger,
		screen: screen,
		time:   time,
	}
}

type finchCtx struct {
	ctx    context.Context
	logger *slog.Logger
	screen *Screen
	time   *Time
}

func (c *finchCtx) Context() context.Context {
	return c.ctx
}

func (c *finchCtx) Screen() *Screen {
	return c.screen
}

func (c *finchCtx) Time() *Time {
	return c.time
}

func (c *finchCtx) Logger() *slog.Logger {
	return c.logger
}

func (c *finchCtx) SetLogger(logger *slog.Logger) Context {
	c.logger = logger
	return c
}

func (c *finchCtx) Get(key ContextKey) any {
	return c.ctx.Value(key)
}

func (c *finchCtx) Set(key ContextKey, value any) Context {
	c.ctx = context.WithValue(c.ctx, key, value)
	return c
}
