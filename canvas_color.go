package go_wasmcanvas

type Color uint32

func MakeColor(r uint8, g uint8, b uint8) uint32 {
	return uint32((r << 16) + (g << 8) + b)
}

const (
	COLOR_BLACK      = 0x0
	COLOR_WHITE      = 0xffffff
	COLOR_RED        = 0x88ffff
	COLOR_CYAN       = 0xaaffee
	COLOR_PURPLE     = 0xcc44cc
	COLOR_GREEN      = 0x00cc55
	COLOR_BLUE       = 0x0000aa
	COLOR_YELLOW     = 0xeeee77
	COLOR_ORANGE     = 0xdd8855
	COLOR_BROWN      = 0x664400
	COLOR_LIGHTRED   = 0xff7777
	COLOR_DARKGRAY   = 0x333333
	COLOR_GRAY       = 0x777777
	COLOR_LIGHTGREEN = 0xaaff66
	COLOR_LIGHTBLUE  = 0x0088ff
	COLOR_LIGHTGRAY  = 0xbbbbbb
)
