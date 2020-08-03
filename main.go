package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const screenWidth = 1024
const screenHeight = 768

var (
	window      *sdl.Window
	quit        bool
	event       sdl.Event
	renderer    *sdl.Renderer
	imageWidth  int32
	imageHeight int32
	textureImg  *sdl.Texture
	imageName   string
)

// Setup - starts SDL, creates window, pre-loads images, sets render quality
func Setup() (successful bool) {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize sdl: %s\n", err)
		return false
	}

	window, err = sdl.CreateWindow("IMG Viewer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create renderer: %s\n", err)
		return false
	}
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create renderer: %s\n", err)
		return false
	}
	renderer.Clear()

	// Unnecessary preloading of jpg and png libraries. Can be commented out and program will automatically load
	// the correct library when you use "img.Load()"
	img.Init(img.INIT_JPG | img.INIT_PNG)

	// SUGGEST to sdl that it use a certain scaling quality for images. Default is "0" a.k.a. nearest pixel sampling
	// try out settings 0, 1, 2 to see the differences with the rotating stick figure. Change the
	// time.Sleep(time.Millisecond * 10) into time.Sleep(time.Millisecond * 100) to slow down the speed of the rotating
	// stick figure and get a good look at how blocky the stick figure is at RENDER_SCALE_QUALITY 0 versus 1 or 2
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	return true
}

func ChangeCurrentImage() {
	switch imageName {
	case "":
		imageName = "01.png"
	case "01.png":
		imageName = "02.png"
	case "02.png":
		imageName = "03.png"
	case "03.png":
		imageName = "01.png"
	}
}

// CreateImage - creates surfaces with sorce image sizes and puts it in texture
func CreateImage() (successful bool) {
	ChangeCurrentImage()
	// Load the glorious programmer art stick figure into memory
	surfaceImg, err := img.Load(imageName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load PNG: %s\n", err)
		os.Exit(4)
	}

	// This is for getting the Width and Height of surfaceImg. Once surfaceImg.Free() is called we lose the
	// ability to get information about the image we loaded into ram
	imageWidth = surfaceImg.W
	imageHeight = surfaceImg.H

	// Take the surfaceImg and use it to create a hardware accelerated textureImg. Or in other words take the image
	// sitting in ram and put it onto the graphics card.
	textureImg, err = renderer.CreateTextureFromSurface(surfaceImg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		os.Exit(5)
	}
	// We have the image now as a texture so we no longer have need for surface. Time to let it go
	surfaceImg.Free()
	return true
}

// HandleEvents - reacts on:
// 1. quit program on cross button
// 2. quit program on ESC press
func HandleEvents() {
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			quit = true
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_RIGHT {
				CreateImage()
			}
			if t.Keysym.Sym == sdl.K_ESCAPE {
				quit = true
			}
		}
	}
}

// Draw - renders background and puts created textures in window
func Draw() {

	renderer.SetDrawColor(255, 255, 55, 255)
	renderer.Clear()

	// Draw the first stick figure using the simpler Copy() function. First parameter is the image we want to draw on
	// screen. Second parameter is the source sdl.Rect of what we want to draw. In this case we instead pass nil, a shortcut
	// telling sdl to draw the entire image. You could use a sdl.Rect to specify drawing only a part of the image - especially
	// useful for animation.
	//
	// The third parameter speficies where on the screen the image will go (X & Y) and how large/small it will be. Alter the
	// 50's to grow or shrink the width and height as desired - or use imageWidth and imageHeight instead to use the normal
	// size of the image.
	// renderer.Copy(textureImg, nil, &sdl.Rect{0, 0, 50, 50})

	// A different way of drawing onto the screen with more options. The first 3 parameters are the same. The fourth
	// parameter is angle of degrees - use 0 if you don't want the image angled.
	//
	// The fifth parameter is to specify a point that the image rotates around. We use nil to use the default
	// Width / 2 and Height / 2 (vertical and horizontal center of image)
	//
	// The Last parameter is the RenderFlip setting. Do you want your image looking normal? Use sdl.FLIP_NONE
	// Do you want your image looking the other way? sdl.FLIP_HORIZONTAL
	// Do you want your image upside down? sdl.SDL_FLIP_VERTICAL
	// Do you want your image upside down AND looking the other way? sdl.FLIP_HORIZONTAL | sdl.SDL_FLIP_VERTICAL
	var scale int32
	scale = 2
	scaledImageWidth := imageWidth / scale
	scaledImageHeight := imageHeight / scale

	renderer.CopyEx(textureImg, nil, &sdl.Rect{0, 0, scaledImageWidth, scaledImageHeight}, 0, nil, sdl.FLIP_NONE)

	renderer.Present()
}

// Shutdown - closes all process to quit program correctly
func Shutdown() {
	// free the texture memory
	textureImg.Destroy()
	// we may or may not use img.Init(), but it's good form to properly shut down the sdl_image library
	img.Quit()
	renderer.Destroy()
	window.Destroy()

	sdl.Quit()
}

func main() {

	if !Setup() {
		os.Exit(1)
	}

	renderer.Clear()

	if !CreateImage() {
		os.Exit(2)
	}

	ticker := time.NewTicker(time.Second / 30)

	for !quit {
		HandleEvents()
		Draw()
		<-ticker.C

	}
	Shutdown()
}
