package finch

import (
	"bytes"
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func RegisterImageAssetManager() {
	RegisterAssetManager(&AssetManager{
		AssetTypes: []AssetType{"png", "jpg", "jpeg", "bmp"},
		ProcessAssetFile: func(file AssetFile, data []byte) (any, error) {
			img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
			if err != nil {
				return nil, err
			}
			return img, nil
		},
		CleanupAssetFile: func(file AssetFile, data any) error {
			img, ok := data.(*ebiten.Image)
			if !ok {
				return errors.New("asset is not an *ebiten.Image")
			}
			img.Deallocate()
			return nil
		},
	})
}

func Get(file AssetFile) (*ebiten.Image, error) {
	img, err := GetAsset[*ebiten.Image](file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func MustGet(file AssetFile) *ebiten.Image {
	return MustGetAsset[*ebiten.Image](file)
}
