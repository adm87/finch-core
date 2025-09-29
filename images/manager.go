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
		AssetTypes: []finch.AssetType{"png", "jpg", "jpeg", "bmp"},
		ProcessAssetFile: func(file finch.AssetFile, data []byte) (any, error) {
			img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
			if err != nil {
				return nil, err
			}
			return img, nil
		},
		CleanupAssetFile: func(file finch.AssetFile, data any) error {
			img, ok := data.(*ebiten.Image)
			if !ok {
				return errors.New("asset is not an *ebiten.Image")
			}
			img.Deallocate()
			return nil
		},
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
