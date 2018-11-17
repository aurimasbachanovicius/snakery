package grafio

type TextOpts struct {
	Size       int32
	XCof, YCof float32
	R, G, B, A uint8
}

type Drawer interface {
	Text(txt string, opts TextOpts) error
	Background(r, g, b, a uint8) error
	Presentation(f func() error) error // somehow drop sdl.Renderer from interface
}

func sizeCal(size int32, cof float32) int32 {
	return int32(float32(size) * (float32(cof)))
}
