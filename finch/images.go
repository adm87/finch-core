package finch

import (
	"bytes"
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	PngAssetType  = "png"
	JpgAssetType  = "jpg"
	JpegAssetType = "jpeg"
	BmpAssetType  = "bmp"
)

func RegisterImageAssetTypes() {
	RegisterAssetManager(&AssetManager{
		AssetTypes: []AssetType{
			PngAssetType,
			JpgAssetType,
			JpegAssetType,
			BmpAssetType,
		},
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
