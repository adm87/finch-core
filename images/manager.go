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
