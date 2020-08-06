package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const screenWidth = int32(1024)
const screenHeight = int32(768)

var (
	window      *sdl.Window
	renderer    *sdl.Renderer
	imageWidth  int32
	imageHeight int32
	textureImg  *sdl.Texture
	imageError  error                           = fmt.Errorf("aaa")
	Rescale     func(W, H int32) (int32, int32) = RescaleNone

	// fileCounter int = 0
	// outCache    []string = os.Args[1:]
	filelist []string
	// scanlist  []string
	fileindex   int = -1
	defaultMask     = "*.png;*.jpg;*.jpeg;*.ico;*.bmp;*.cur;*.pnm;*.xpm;*.lbm;*.pcx;*.gof;*.tga;*.tiff;*.xv;*.ppm;*.pgm;*.pbm;*.iff;*.ilbmo"
)

// Setup - starts SDL, creates window, pre-loads images, sets render quality
func Setup() (successful bool) {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize sdl: %s\n", err)
		return false
	}

	window, err = sdl.CreateWindow("IMG Viewer", 0, 0,
		screenWidth, screenHeight, sdl.WINDOW_BORDERLESS)
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

func Sobaka() {
	fmt.Println("Woof-woof")
}

func CustomGlob(glob string) ([]string, []string) {
	argList := strings.Split(glob, ";")
	ret := []string{}
	errList := []string{}
	dupmap := map[string]int{}
	for _, s := range argList {
		searchList, err := filepath.Glob(s)
		if err != nil {
			errList = append(errList, err.Error())
			continue
		}
		for _, filename := range searchList {
			val := dupmap[filename]
			if val != 0 {
				continue
			}
			ret = append(ret, filename)
			dupmap[filename] = 1

		}
		// newName := retMap[searchList]

		// for _, a := range searchList {
		// 	newName := retMap[a]
		// 	fmt.Println(newName)

		// }

		fmt.Println("ret is", ret)
		fmt.Println("############")
	}

	if len(errList) > 0 {
		errList = nil
	}

	return ret, errList
}

func ParseArgs(args []string) []string {
	ret := []string{}
	// noFlags := true
	// CustomFlag := false
	doSort := false
	maskList := []string{}
	// isDefaultFilterMode := true
	// es := 0
	errors := []string{}
	for _, arg := range args[1:] {
		customFlag := strings.HasPrefix(arg, "-")
		if customFlag {
			fmt.Println(arg, "is custom flag")
			// filters = append(filters, arg)
			switch arg {
			default:
				errors = append(errors, fmt.Sprintf("Unsupported argument %v", arg))
				continue
			case "-all":
				arg = defaultMask
			case "-sort":
				doSort = true
				continue
			case "-sobaka":
				fmt.Println(arg, "is fine")
				Sobaka()
				continue
			}
		}
		maskList = append(maskList, arg)
	}

	if len(maskList) == 0 {
		maskList = append(maskList, defaultMask)
	}

	errList := []string{}
	ret, errList = CustomGlob(strings.Join(maskList, ";"))
	errors = append(errors, errList...)

	if len(errors) > 0 {
		for _, s := range errors {
			fmt.Println(s)
		}
		os.Exit(2)
	}

	if doSort {
		fmt.Println(doSort)
		sort.Strings(ret)
	}

	return ret
}

// CreateImage - creates surfaces with sorce image sizes and puts it in texture
func CreateImage(file string) (successful bool) {
	// ChangeCurrentImage()

	// currentFilename := NextFile(wDir)
	currentFilename := file
	surfaceImg, err := img.Load(currentFilename)
	imageError = err
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load image: %s\n", imageError)
		return false
	}

	// )
	// os.Exit(4)

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
	currentFilename = ""
	return true
}

// HandleEvents - reacts on:
// 1. quit program on cross button
// 2. quit program on ESC press
func HandleEvents() {
	doDraw := true
	quit := false
	scaleNone := true
	for !quit {
		if doDraw {
			Draw()
			doDraw = false
		}
		for event := sdl.WaitEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				quit = true
			case *sdl.KeyboardEvent:
				if t.Type != sdl.KEYDOWN {
					continue
				}

				switch t.Keysym.Sym {
				case sdl.K_RIGHT:
					fmt.Println("next file")

					CreateImage(getNextFile())
					doDraw = true
				case sdl.K_LEFT:
					fmt.Println("prev file")

					CreateImage(getPrevFile())

					doDraw = true
				case sdl.K_ESCAPE:
					quitEvent := sdl.QuitEvent{Type: sdl.QUIT}
					sdl.PushEvent(&quitEvent)
				case sdl.K_f:
					Rescale = RescaleNone
					scaleNone = !scaleNone
					if !scaleNone {
						Rescale = RescaleFit
					}
					doDraw = true
				}
			}
		}
	}
}

func RescaleFit(w, h int32) (int32, int32) {
	if screenWidth <= 0 || screenHeight <= 0 || w <= 0 || h <= 0 {
		return RescaleNone(w, h)
	}

	imgW := float64(w)
	imgH := float64(h)
	scrW := float64(screenWidth)
	scrH := float64(screenHeight)

	k := scrW / imgW
	if (imgW / imgH) < (scrW / scrH) {
		k = scrH / imgH
	}

	scaledImageWidth := math.Round(k * imgW)
	scaledImageHeight := math.Round(k * imgH)
	return int32(scaledImageWidth), int32(scaledImageHeight)

}

func RescaleNone(w, h int32) (int32, int32) {
	return w, h
}

func DrawCross() {
	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.DrawLine(0, 0, screenWidth, screenHeight)
	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.DrawLine(screenWidth, 0, 0, screenHeight)
}

// Draw - renders background and puts created textures in window
func Draw() {
	fmt.Printf("draw start\n")

	defer renderer.Present()
	renderer.SetDrawColor(255, 255, 55, 255)
	renderer.Clear()

	if imageError != nil {

		DrawCross()
		return
	}

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
	newWidth, newHeight := Rescale(imageWidth, imageHeight)
	offsetX := (screenWidth - newWidth) / 2
	offsetY := (screenHeight - newHeight) / 2

	renderer.CopyEx(textureImg, nil, &sdl.Rect{offsetX, offsetY, newWidth, newHeight}, 0, nil, sdl.FLIP_NONE)
	fmt.Printf("draw end\n")

}

// ###############################

func getCurFile() string {
	if fileindex == -1 {
		fileindex = 0
	}
	if fileindex > len(filelist)-1 {
		fileindex = -1
	}
	ret := ""
	if fileindex != -1 {
		ret = filelist[fileindex]
	}
	fmt.Printf("%v/%v | %v\n", fileindex+1, len(filelist), ret)
	return ret
}
func getNextFile() string {
	fileindex++
	if fileindex > len(filelist)-1 {
		fileindex = 0
	}
	return getCurFile()
}
func getPrevFile() string {
	fileindex--
	if fileindex < 0 {
		fileindex = len(filelist) - 1
	}
	return getCurFile()
}

// ############################################

func main() {

	if !Setup() {
		os.Exit(1)
	}
	filelist = ParseArgs(os.Args)
	// scanlist = ScanDir(filelist)

	CreateImage(getCurFile())

	HandleEvents()
	Shutdown()
}
