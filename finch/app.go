package finch

import (
	"context"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	firstUpdate   = true
	shouldExit    = false
	currentWidth  = 0
	currentHeight = 0
)

func Run(a *App) error {
	if window := a.window; window != nil {
		ebiten.SetWindowTitle(window.Title)
		ebiten.SetWindowSize(window.Width, window.Height)
		ebiten.SetWindowResizingMode(window.ResizingMode)
		ebiten.SetFullscreen(window.Fullscreen)
	}
	return ebiten.RunGame(a)
}

func Exit() {
	shouldExit = true
}

type (
	DrawFunc     func(ctx Context, screen *ebiten.Image)
	LayoutFunc   func(ctx Context, outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
	UpdateFunc   func(ctx Context)
	StartupFunc  func(ctx Context)
	ShutdownFunc func(ctx Context)
)

type Window struct {
	Title        string
	ResizingMode ebiten.WindowResizingModeType
	Width        int
	Height       int
	Fullscreen   bool
	RenderScale  float64
}

type App struct {
	ctx Context

	DrawFn        DrawFunc
	LayoutFn      LayoutFunc
	UpdateFn      UpdateFunc
	FixedUpdateFn UpdateFunc
	LateUpdateFn  UpdateFunc
	StartupFn     StartupFunc
	ShutdownFn    ShutdownFunc

	window *Window
}

func NewApp(ctx context.Context, logger *slog.Logger) *App {
	s := NewScreen(800, 600, 1.0, false)
	t := NewTime(60.0)
	return &App{ctx: NewContext(ctx, logger, s, t)}
}

func (a *App) WithDraw(drawFunc DrawFunc) *App {
	a.DrawFn = drawFunc
	return a
}

func (a *App) WithLayout(layoutFunc LayoutFunc) *App {
	a.LayoutFn = layoutFunc
	return a
}

func (a *App) WithUpdate(updateFunc UpdateFunc) *App {
	a.UpdateFn = updateFunc
	return a
}

func (a *App) WithFixedUpdate(fixedUpdateFunc UpdateFunc) *App {
	a.FixedUpdateFn = fixedUpdateFunc
	return a
}

func (a *App) WithLateUpdate(lateUpdateFunc UpdateFunc) *App {
	a.LateUpdateFn = lateUpdateFunc
	return a
}

func (a *App) WithStartup(startupFunc StartupFunc) *App {
	a.StartupFn = startupFunc
	return a
}

func (a *App) WithShutdown(shutdownFunc ShutdownFunc) *App {
	a.ShutdownFn = shutdownFunc
	return a
}

func (a *App) WithWindow(window *Window) *App {
	a.window = window
	a.ctx.Screen().set_target_size(window.Width, window.Height)
	a.ctx.Screen().set_render_scale(window.RenderScale)
	a.ctx.Screen().set_fullscreen(window.Fullscreen)
	return a
}

func (a *App) Draw(screen *ebiten.Image) {
	if draw := a.DrawFn; draw != nil {
		draw(a.ctx, screen)
	}
}

func (a *App) Update() error {
	if shouldExit {
		a.ctx.Logger().Info("Shutting down application")
		if shutdown := a.ShutdownFn; shutdown != nil {
			shutdown(a.ctx)
		}
		return ebiten.Termination
	}

	if firstUpdate {
		a.ctx.Logger().Info("Starting up application")
		a.ctx.Time().start()
		if startup := a.StartupFn; startup != nil {
			startup(a.ctx)
		}
		firstUpdate = false
	}

	a.ctx.Time().tick()

	if update := a.UpdateFn; update != nil {
		update(a.ctx)
	}
	if fixedUpdate := a.FixedUpdateFn; fixedUpdate != nil {
		for i := 0; i < a.ctx.Time().FixedFrames(); i++ {
			fixedUpdate(a.ctx)
		}
	}
	if lateUpdate := a.LateUpdateFn; lateUpdate != nil {
		lateUpdate(a.ctx)
	}

	return nil
}

func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if layout := a.LayoutFn; layout != nil {
		screenWidth, screenHeight = layout(a.ctx, outsideWidth, outsideHeight)
	} else {
		screenWidth = int(float64(a.ctx.Screen().TargetWidth()) * a.ctx.Screen().RenderScale())
		screenHeight = int(float64(a.ctx.Screen().TargetHeight()) * a.ctx.Screen().RenderScale())
		a.ctx.Screen().SetSize(screenWidth, screenHeight)
	}

	if currentWidth != screenWidth || currentHeight != screenHeight {
		oldWidth := currentWidth
		oldHeight := currentHeight
		currentWidth = screenWidth
		currentHeight = screenHeight

		a.ctx.Logger().Info("Resized screen",
			slog.Int("oldWidth", oldWidth),
			slog.Int("oldHeight", oldHeight),
			slog.Int("newWidth", currentWidth),
			slog.Int("newHeight", currentHeight),
		)
	}

	return screenWidth, screenHeight
}

func (a *App) Context() Context {
	return a.ctx
}
