package wallpaper

type Layout string

const (
	LayoutTile    Layout = "tile"
	LayoutStretch Layout = "stretch"
	LayoutFit     Layout = "fit"
	LayoutFill    Layout = "fill"
	LayoutCenter  Layout = "center"
)

func Set(path string, layout Layout) error {
	return setWallpaper(path, layout)
}
