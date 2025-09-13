package finch

// Screen represents the game's screen with width and height.
//
// Game logic should reference this to get the current screen dimensions.
type Screen struct {
	targetWidth  int
	targetHeight int
	width        int
	height       int
	renderScale  float64
	fullscreen   bool
}

func NewScreen(width, height int, scale float64, fullscreen bool) *Screen {
	return &Screen{
		targetWidth:  width,
		targetHeight: height,
		width:        width,
		height:       height,
		renderScale:  scale,
		fullscreen:   fullscreen,
	}
}

func (s *Screen) Width() int {
	return s.width
}

func (s *Screen) Height() int {
	return s.height
}

func (s *Screen) TargetWidth() int {
	return s.targetWidth
}

func (s *Screen) TargetHeight() int {
	return s.targetHeight
}

func (s *Screen) RenderScale() float64 {
	return s.renderScale
}

func (s *Screen) set_render_scale(scale float64) {
	if scale <= 0 {
		panic("render scale must be greater than 0")
	}
	s.renderScale = scale
}

func (s *Screen) set_target_size(width, height int) {
	if width <= 0 {
		panic("target width must be greater than 0")
	}
	if height <= 0 {
		panic("target height must be greater than 0")
	}
	s.targetWidth = width
	s.targetHeight = height
}

func (s *Screen) SetSize(width, height int) {
	if width <= 0 {
		panic("width must be greater than 0")
	}
	if height <= 0 {
		panic("height must be greater than 0")
	}
	s.width = width
	s.height = height
}

func (s *Screen) IsFullscreen() bool {
	return s.fullscreen
}

func (s *Screen) set_fullscreen(fullscreen bool) {
	s.fullscreen = fullscreen
}
