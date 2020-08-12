package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	textColor   = sdl.Color{255, 255, 255, 255}
	defaultFont *ttf.Font
	// LatoRegular20 ...
	//LatoRegular20 *ttf.Font
	// LatoRegular24 ...
	//LatoRegular24 *ttf.Font
	// LatoRegular14 ...
	//LatoRegular14 *ttf.Font
	// LatoRegular12 ...
	//LatoRegular12 *ttf.Font
)

// InitFonts ...
func InitFonts() {
	fmt.Println("init font")
	err := ttf.Init()
	if err != nil {
		panic(err)
	}

	// rwops := sdl.RWFromMem(unsafe.Pointer(&goregular.TTF[0]), len(goregular.TTF))
	rwops, err := sdl.RWFromMem(goregular.TTF)
	if err != nil {
		panic(err)
	}
	defaultFont, err = ttf.OpenFontRW(rwops, 1, 15)
	if err != nil {
		panic(err)
	}
	// LatoRegular20, _ = ttf.OpenFontRW(rwops, 1, 20)
	// LatoRegular24, _ = ttf.OpenFontRW(rwops, 1, 24)
	// LatoRegular14, _ = ttf.OpenFontRW(rwops, 1, 14)
	// LatoRegular12, _ = ttf.OpenFontRW(rwops, 1, 12)
}

func SetTextColor(color sdl.Color) {
	textColor = color
}
func SetTextColorRGBA(r, g, b, a uint8) {
	SetTextColor(sdl.Color{r, g, b, a})
}

func WriteTextCustom(mode int, x, y int32, str string) (int32, int32) {
	w, h := int32(0), int32(0)
	switch mode {
	default:
		panic("WriteTextCustom : unsupported text mode")
	case 0:
		w, h = WriteText(x, y, str)
	case 1:
		SetTextColorRGBA(0, 0, 0, 255)
		WriteText(x+1, y+1, str)
		SetTextColorRGBA(255, 255, 255, 255)
		w, h = WriteText(x, y, str)
	case 2:
		SetTextColorRGBA(0, 0, 0, 255)
		WriteText(x+1, y, str)
		WriteText(x-1, y, str)
		WriteText(x, y+1, str)
		WriteText(x, y-1, str)
		SetTextColorRGBA(255, 255, 255, 255)
		w, h = WriteText(x, y, str)
	case 3:
		textureTxt := (*sdl.Texture)(nil)
		textureTxt, w, h = TextTexture(x, y, str)
		renderer.SetDrawColor(0, 0, 0, 128)
		renderer.FillRect(&sdl.Rect{x, y, w, h})
		renderer.CopyEx(textureTxt, nil, &sdl.Rect{x, y, w, h}, 0, nil, sdl.FLIP_NONE)
		textureTxt.Destroy()
	case 4:
		textureTxt := (*sdl.Texture)(nil)
		textureTxt, w, h = TextTexture(x, y, str)
		renderer.SetDrawColor(0, 0, 0, 128)
		renderer.FillRect(&sdl.Rect{x - 3, y, w + 6, h})
		renderer.CopyEx(textureTxt, nil, &sdl.Rect{x, y, w, h}, 0, nil, sdl.FLIP_NONE)
		textureTxt.Destroy()
	}
	return w + x, h + y
}

func TextTexture(x, y int32, str string) (*sdl.Texture, int32, int32) {
	fmt.Println("write text", defaultFont)
	// if str == "" {
	// 	return nil, x, y
	// }
	solid, err := defaultFont.RenderUTF8Blended(str, textColor)
	if err != nil {
		panic("defaultFont.RenderUTF8Blended : " + err.Error())
	}
	defer solid.Free()

	textureTxt, err := renderer.CreateTextureFromSurface(solid)
	if err != nil {
		panic("CreateTextureFromSurface : " + err.Error())
	}
	return textureTxt, solid.W, solid.H
}

func WriteText(x, y int32, str string) (int32, int32) {

	// fmt.Println("write text", defaultFont)
	// if str == "" {
	// 	return x, y
	// }
	// solid, err := defaultFont.RenderUTF8Blended(str, textColor)
	// if err != nil {
	// 	fmt.Fprint(os.Stderr, "Failed to render text: %s\n", err)
	// 	return 0, 0
	// }
	// defer solid.Free()

	// textureTxt, err := renderer.CreateTextureFromSurface(solid)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to create text texture: %s\n", err)
	// 	os.Exit(5)
	// }
	textureTxt, w, h := TextTexture(x, y, str)
	renderer.CopyEx(textureTxt, nil, &sdl.Rect{x, y, w, h}, 0, nil, sdl.FLIP_NONE)
	textureTxt.Destroy()
	return w, h
}
