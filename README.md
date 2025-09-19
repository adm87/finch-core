# Finch Core
The Finch Core module is the base for all other Finch modules. It provides an wrapper to and launch an Ebitengine application, along with other utilities used across the Finch ecosystem.

This core module is agnostic of all other Finch modules and should be considered as the minimum entry point into a Finch application.

## Usage
```go
package main

import (
	"github.com/adm87/finch-core/finch"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	f := finch.NewApp().
		WithWindow(&finch.Window{
			Title:        "Finch Test",
			ResizingMode: ebiten.WindowResizingModeDisabled,
			Width:        800,
			Height:       600,
			Fullscreen:   false,
			RenderScale:  1.0,
		}).
		WithStartup(startup).
		WithDraw(draw)

	if err := finch.Run(f); err != nil {
		panic(err)
	}
}

func startup(ctx finch.Context) {
	ctx.Logger().Info("Hello, Finch!")
}

func draw(ctx finch.Context, screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, Finch!")
}
```