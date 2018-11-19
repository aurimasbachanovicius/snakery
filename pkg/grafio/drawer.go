package grafio

type RGBA struct {
	R, G, B, A uint8
}
type TextAlign int

const Right TextAlign = 2

type TextOpts struct {
	Size       int32
	XCof, YCof float32
	Color      RGBA
	Align      TextAlign
}

type RectOpts struct {
	Texture string
	Color   RGBA
}

type Drawer interface {
	Background(r, g, b, a uint8) error

	Text(txt string, opts TextOpts) error
	ColorRect(x, y, w, h int32, rgba RGBA) error
	TextureRect(x, y, w, h int32, texture string) error

	Present(f func() error) error
	LoadResources(fontsPath, texturesPath string) (func() error, error)
	SetMainFont(fontFileName string) error

	ScreenHeight() int32
	ScreenWidth() int32
}

func sizeCal(size int32, cof float32) int32 {
	return int32(float32(size) * (float32(cof)))
}
