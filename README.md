# Finch Core
The Finch Core module is the base for all other Finch modules. It provides a wrapper to configure and launch an Ebitengine application, along with other utilities used across the Finch ecosystem.

This core module is agnostic of all other Finch modules and should be considered the minimum entry point into a Finch application.

## Usage
### Quick Startup
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
### Assets
#### Loading
The goal is to make loading and accessing game assets simple and straight forward. A `finch.AssetFile` provides a handle for loading, unload, and accessing asset data.
```go
package main

import (
	...

	"github.com/adm87/finch-core/finch"
	"github.com/adm87/finch-core/images"
)

var myAwesomePng = finch.AssetFile("absolute/path/to/assets/image.png")

func main() {
	images.RegisterAssetManager()

	...

	if err := myAwesomePng.Load(); err != nil {
		panic(err)
	}

	// OR

	myAwesomePng.MustLoad()

	// OR

	if err := finch.LoadAssets(myAwesomePng); err != nil {
		panic(err)
	}

	// OR

	finch.MustLoadAssets(myAwesomePng)
}
```
Finch's asset management also provides a mechanism for registering `fs.FS` filesystem implementations. This provides a custom solution for reading asset data from user defined locations.
```go
package main

import (
	"os"

	...

	"github.com/adm87/finch-core/finch"
)
var assets = AssetRoot("assets")

func main() {
	finch.RegisterAssetFilesystem(assets, os.DirFS("absolute/path/to/assets"))

	...
}
```
Here, `AssetRoot("assets")` defines the root directory for a custom filesystem. Since this example is creating a filesystem for a local disk location, `absolute/path/to/assets` is the absolute path where that filesystem starts. This allows the `AssetFile()` paths to be relative to the filesystem they belong to.
```go
var myAwesomePng = finch.AssetFile("assets/image.png")
```
When Finch attempts to load the `AssetFile`, is will check the top level directory to see if a filesystem has been registered for it. If so, it will use that filesystem with the relative path of the `AssetFile` to load the asset. This can be useful if assets are located on a remote server. Implement a custom `fs.FS` that knows how to connect to and read from the server, and register it to the root of that filesystem.
> Note: If Finch doesn't find a registered filesystem, it will attempt to use the path of the `AssetFile` to load from disk.

#### Accessing
After an `AssetFile` has been loaded, you can use it to get the untyped data associated with it. The data will need to be type casted to the expected type.
```go
func draw(ctx finch.Context, screen *ebiten.Image) {
	img, err := myAwesomePng.Get()
	if err != nil {
		panic(err)
	}
	screen.DrawImage(img.(*ebiten.Image), nil
	
	// OR

	screen.DrawImage(myAwesomePng.MustGet().(*ebiten.Image), nil)

	// OR

	img, err := images.Get(myAwesomePng)
	if err != nil {
		panic(err)
	}
	screen.DrawImage(img, nil)

	// OR

	screen.DrawImage(images.MustGet(myAwesomePng), nil)
}
```
#### Custom Asset Types
Finch comes with built-in asset managers for loading asset types common to the Ebitengine. To use them simply call their `RegisterAssetManager()` methods. Or, you can build your own asset manager and have Finch use that instead.

A custom asset manager provides users a way to manage how loaded data is allocated and deallocated within Finch's asset framework.
```go
package images

import (
	"bytes"
	"errors"

	"github.com/adm87/finch-core/finch"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func RegisterAssetManager() {
	finch.RegisterAssetManager(&finch.AssetManager{
		Types:       []finch.AssetType{"png", "jpg", "jpeg", "bmp"},
		Allocator:   allocator,
		Deallocator: deallocator,
	})
}

func Get(file finch.AssetFile) (*ebiten.Image, error) {
	img, err := finch.GetAsset[*ebiten.Image](file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func MustGet(file finch.AssetFile) *ebiten.Image {
	return finch.MustGetAsset[*ebiten.Image](file)
}

func allocator(data []byte) (any, error) {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func deallocator(asset any) error {
	img, ok := asset.(*ebiten.Image)
	if !ok {
		return errors.New("asset is not an *ebiten.Image")
	}
	img.Deallocate()
	return nil
}
```
This is Finch's built-in image asset manager. It tells Finch how to allocate and deallocate files with the extension `png, jpg, jpeg, and bmp`. This package also provides some convenient accessor methods for returning typed data.
> Note: Only one asset manager per file extension is supported. If you would rather manage images yourself, create a custom asset manager and don't bother calling Finch's images' RegisterAssetManager() method.