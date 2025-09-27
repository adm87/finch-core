package finch

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/adm87/finch-core/linq"
	"github.com/adm87/finch-core/types"
)

var (
	ErrAssetManagerNotFound    = errors.New("asset manager not found")
	ErrAssetManagerConflict    = errors.New("asset manager conflict")
	ErrAssetManagerNil         = errors.New("asset manager is nil")
	ErrAssetFilesystemConflict = errors.New("asset filesystem conflict")
	ErrAssetFilesystemNil      = errors.New("asset filesystem is nil")
	ErrAssetInvalidType        = errors.New("asset has invalid type")
	ErrAssetTypeMismatch       = errors.New("asset type mismatch")
	ErrAssetNotLoaded          = errors.New("asset not loaded")
	ErrAssetIsLoaded           = errors.New("asset is already loaded")
	ErrAssetIsLoading          = errors.New("asset is currently loading")
	ErrAssetTypeEmpty          = errors.New("asset type is empty")
	ErrAssetRootEmpty          = errors.New("asset root is empty")
)

// ======================================================
// Asset Manager
// ======================================================

// AssetManager managers allocation and deallocation of a specific asset types.
type AssetManager struct {
	Allocator   AssetAllocator
	Deallocator AssetDeallocator
	Types       []AssetType
}

// ======================================================
// Asset Allocation/Deallocation
// ======================================================

// AssetAllocator is a function that takes raw asset data and converts it into a usable form.
type AssetAllocator func(file AssetFile, filedata []byte) (any, error)

// AssetDeallocator is a function that takes a loaded asset and frees its resources.
type AssetDeallocator func(file AssetFile, data any) error

// ======================================================
// Asset Type
// ======================================================

// AssetType represents the type of an asset, typically derived from its file extension.
type AssetType string

func (t AssetType) String() string {
	return string(t)
}

func (t AssetType) IsValid() error {
	if t == "" {
		return ErrAssetTypeEmpty
	}
	return nil
}

// ======================================================
// Asset Root
// ======================================================

// AssetRoot represents the root directory of an asset file.
type AssetRoot string

func (r AssetRoot) String() string {
	return string(r)
}

func (r AssetRoot) IsValid() error {
	if r == "" {
		return ErrAssetRootEmpty
	}
	return nil
}

// ======================================================
// Asset File
// ======================================================

// AssetFile is a handle to an asset file and its associated type and data.
type AssetFile string

// Path returns the file path of the asset.
func (f AssetFile) Path() string {
	return string(f)
}

// Root returns the root directory of the asset file.
func (f AssetFile) Root() AssetRoot {
	fpath := strings.ReplaceAll(filepath.Clean(f.Path()), "\\", "/")
	fpath = strings.TrimPrefix(fpath, "/")

	parts := strings.Split(fpath, "/")
	if len(parts) == 0 {
		return ""
	}

	return AssetRoot(parts[0])
}

// Type returns the asset type based on the file extension.
func (f AssetFile) Type() AssetType {
	return AssetType(filepath.Ext(f.Path())[1:])
}

func (f AssetFile) Load() error {
	return LoadAssets(f)
}

func (f AssetFile) MustLoad() {
	if err := f.Load(); err != nil {
		panic(err)
	}
}

func (f AssetFile) Unload() error {
	return UnloadAssets(f)
}

func (f AssetFile) MustUnload() {
	if err := f.Unload(); err != nil {
		panic(err)
	}
}

func (f AssetFile) Get() (any, error) {
	return GetAsset[any](f)
}

func (f AssetFile) MustGet() any {
	data, err := f.Get()

	if err != nil {
		panic(err)
	}

	return data
}

// ======================================================
// Asset Management
// ======================================================

var (
	assetCache       = make(map[AssetFile]any)
	assetManagers    = make(map[AssetType]*AssetManager)
	assetFilesystems = make(map[AssetRoot]fs.FS)
	assetsLoading    = make(types.HashSet[AssetFile])
	assetsMu         = sync.RWMutex{}
)

func HasAssetTypeSupport(t AssetType) bool {
	_, exists := assetManagers[t]
	return exists
}

func RegisterAssetManager(manager *AssetManager) error {
	if manager == nil {
		return ErrAssetManagerNil
	}

	for _, t := range manager.Types {
		if err := t.IsValid(); err != nil {
			return fmt.Errorf("%s: %w", ErrAssetInvalidType, err)
		}

		if _, exists := assetManagers[t]; exists {
			return fmt.Errorf("%s: %w", ErrAssetManagerConflict, t)
		}

		assetManagers[t] = manager
	}

	return nil
}

func RegisterAssetFilesystem(root AssetRoot, filesystem fs.FS) error {
	if err := root.IsValid(); err != nil {
		return err
	}

	if filesystem == nil {
		return errors.New("filesystem is nil")
	}

	if _, exists := assetFilesystems[root]; exists {
		return fmt.Errorf("%s: %s", ErrAssetFilesystemConflict, root)
	}

	assetFilesystems[root] = filesystem

	return nil
}

func GetAsset[T any](file AssetFile) (T, error) {
	data, ok := assetCache[file]

	if !ok {
		return *new(T), fmt.Errorf("%s: %s", ErrAssetNotLoaded, file)
	}

	if typed, ok := data.(T); ok {
		return typed, nil
	}

	return *new(T), fmt.Errorf("%s: %s", ErrAssetTypeMismatch, file)
}

func MustGetAsset[T any](file AssetFile) T {
	data, err := GetAsset[T](file)

	if err != nil {
		panic(err)
	}

	return data
}

func LoadAssets(files ...AssetFile) error {
	if len(files) == 0 {
		return nil
	}

	requests := make(types.HashSet[AssetFile])
	errs := make([]error, 0)

	if err := build_asset_requests(requests, files); err != nil {
		errs = append(errs, err)
	}

	if len(requests) == 0 {
		return errors.Join(errs...)
	}

	if err := load_asset_batches(linq.Batch(requests.ToSlice(), 100)); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func MustLoadAssets(files ...AssetFile) {
	if err := LoadAssets(files...); err != nil {
		panic(err)
	}
}

func UnloadAssets(files ...AssetFile) error {
	assetsMu.Lock()
	defer assetsMu.Unlock()

	for _, file := range files {
		asset, exists := assetCache[file]
		if !exists {
			return fmt.Errorf("%s: %s", ErrAssetNotLoaded, file)
		}

		manager, exists := assetManagers[file.Type()]
		if !exists {
			return fmt.Errorf("%s: %s", ErrAssetManagerNotFound, file.Type())
		}

		if manager.Deallocator != nil {
			if err := manager.Deallocator(file, asset); err != nil {
				return fmt.Errorf("failed to deallocate asset %s: %w", file, err)
			}
		}

		delete(assetCache, file)
	}

	return nil
}

func MustUnloadAssets(files ...AssetFile) {
	if err := UnloadAssets(files...); err != nil {
		panic(err)
	}
}

func build_asset_requests(requests types.HashSet[AssetFile], files []AssetFile) error {
	errs := make([]error, 0)

	for _, file := range files {
		if _, exists := requests[file]; exists {
			continue
		}

		fileType := file.Type()

		if err := fileType.IsValid(); err != nil {
			errs = append(errs, err)
			continue
		}

		if _, exists := assetManagers[fileType]; !exists {
			errs = append(errs, fmt.Errorf("%s: %s", ErrAssetManagerNotFound, fileType))
			continue
		}

		requests.Add(file)
	}

	return errors.Join(errs...)
}

func load_asset_batches(batches [][]AssetFile) error {
	if len(batches) == 1 {
		return load_asset_batch(batches[0])
	}

	panicCh := make(chan error, len(batches))
	wg := sync.WaitGroup{}

	wg.Add(len(batches))
	for _, batch := range batches {
		go func(files []AssetFile) {
			defer wg.Done()

			defer func() {
				if r := recover(); r != nil {
					panicCh <- fmt.Errorf("panic in asset load batch: %v", r)
				}
			}()

			if err := load_asset_batch(files); err != nil {
				panicCh <- err
			}
		}(batch)
	}
	wg.Wait()

	close(panicCh)

	errs := make([]error, 0, len(panicCh))
	for err := range panicCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func load_asset_batch(files []AssetFile) error {
	if len(files) == 0 {
		return nil
	}

	errs := make([]error, 0)

	// Note: Errors don't interrupt loading subsequent assets.
	// Instead all errors are returned to be handled upstream.
	for _, file := range files {
		if err := load_asset_file(file); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func load_asset_file(file AssetFile) error {
	if err := try_load(file); err != nil {
		return err
	}

	defer func() {
		assetsMu.Lock()
		assetsLoading.Remove(file)
		assetsMu.Unlock()
	}()

	var data []byte
	var err error

	froot := file.Root()
	fpath := file.Path()

	filesystem, exists := assetFilesystems[froot]
	if exists {
		switch _, ok := filesystem.(embed.FS); {
		case ok:
			fpath = filepath.Join(froot.String(), fpath)
		default:
			fpath = strings.TrimPrefix(fpath, froot.String())
			fpath = strings.TrimPrefix(fpath, string(filepath.Separator))
		}
		data, err = fs.ReadFile(filesystem, fpath)
		if err != nil {
			return err
		}
	} else {
		data, err = os.ReadFile(fpath)
		if err != nil {
			return err
		}
	}

	manager, exists := assetManagers[file.Type()]
	if !exists {
		return fmt.Errorf("%s: %s", ErrAssetManagerNotFound, file.Type())
	}

	if manager.Allocator == nil {
		return fmt.Errorf("%s: %s", ErrAssetManagerNil, file.Type())
	}

	asset, err := manager.Allocator(file, data)
	if err != nil {
		return fmt.Errorf("failed to import asset %s: %w", file, err)
	}

	assetsMu.Lock()
	assetCache[file] = asset
	assetsMu.Unlock()

	return nil
}

func try_load(file AssetFile) error {
	assetsMu.Lock()
	defer assetsMu.Unlock()

	if _, exists := assetCache[file]; exists {
		return fmt.Errorf("%w: %s", ErrAssetIsLoaded, file)
	}
	if assetsLoading.Contains(file) {
		return fmt.Errorf("%w: %s", ErrAssetIsLoading, file)
	}

	assetsLoading.Add(file)

	return nil
}
