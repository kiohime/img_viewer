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

const windowPositionX = int32(200)
const windowPositionY = int32(50)
const defaultMask = "*.png;*.jpg;*.jpeg;*.ico;*.bmp;*.cur;*.pnm;*.xpm;*.lbm;*.pcx;*.gof;*.tga;*.tiff;*.xv;*.ppm;*.pgm;*.pbm;*.iff;*.ilbmo"

var (
	window              *sdl.Window
	screenWidth         = int32(1025)
	screenHeight        = int32(769)
	fullscreen          = false
	renderer            *sdl.Renderer
	imageWidth          int32
	imageHeight         int32
	textureImg          *sdl.Texture
	textureCheckerboard *sdl.Texture
	imageError          error                           = fmt.Errorf("aaa")
	Rescale             func(W, H int32) (int32, int32) = RescaleNone

	// fileCounter int = 0
	// outCache    []string = os.Args[1:]
	filelist []string
	// scanlist  []string
	fileindex      int = -1
	patternMode        = 0
	zalivkaChB         = false
	bgChB              = false
	patternZalivka string
	patternBg      string
	defaultZalivka = "checkerDark"
	defaultBg      = "checkerLight"
)

// Setup - starts SDL, creates window, pre-loads images, sets render quality
func Setup() (successful bool) {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize sdl: %s\n", err)
		return false
	}

	window, err = sdl.CreateWindow("IMG Viewer", windowPositionX, windowPositionY,
		screenWidth, screenHeight, sdl.WINDOW_RESIZABLE)
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
	dupIndex := 0
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
			dupIndex++
			ret = append(ret, filename)
			dupmap[filename] = dupIndex

		}
		// fmt.Println("dupmap is", dupmap)

		// fmt.Println("ret is", ret)
		fmt.Println("############")
	}

	if len(errList) > 0 {
		errList = nil
	}

	return ret, errList
}

func getArgument(arg, prefix string, result *string) bool {
	// fmt.Println("im here#########")
	if strings.HasPrefix(arg, prefix) {
		*result = strings.TrimPrefix(arg, prefix)
		return true
	}
	return false
}

func SetFullscreen(x bool) {
	mode := uint32(0)
	if x {
		mode = sdl.WINDOW_FULLSCREEN_DESKTOP
	}
	window.SetFullscreen(mode)
}

func ParseArgs(args []string) []string {
	ret := []string{}
	// noFlags := true
	// CustomFlag := false
	doSort := false
	maskList := []string{}
	filePosName := ""
	// isDefaultFilterMode := true
	// es := 0
	errors := []string{}
	for i, arg := range args[1:] {
		fmt.Println(args[i+1])
		customFlag := strings.HasPrefix(arg, "-")
		if customFlag {
			fmt.Println(arg, "is custom flag")
			// filters = append(filters, arg)
			switch arg {
			default:
				ok := false
				ok = ok || getArgument(arg, "-x:", &filePosName)
				ok = ok || getArgument(arg, "-zalivka:", &defaultZalivka)
				ok = ok || getArgument(arg, "-bg:", &defaultBg)

				if !ok {
					errors = append(errors, fmt.Sprintf("Unsupported argument %v", arg))
				}
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
			case "-fullscreen":
				fullscreen = true
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
	if filePosName != "" {
		for i, name := range ret {
			if name == filePosName {
				fileindex = i
			}
		}
	}

	return ret
}

// CreateImage - creates surfaces with sorce image sizes and puts it in texture
func CreateImage(file string) (successful bool) {

	currentFilename := file
	surfaceImg, err := img.Load(currentFilename)
	imageError = err
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load image: %s\n", imageError)
		return false
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
	modifers := false

	for !quit {
		if doDraw {
			Draw()
			doDraw = false
		}
		for event := sdl.WaitEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.WindowEvent:
				switch t.Event {
				case sdl.WINDOWEVENT_RESIZED:
					screenWidth, screenHeight = window.GetSize()
					doDraw = true
				}
			case *sdl.QuitEvent:
				quit = true
			case *sdl.KeyboardEvent:
				switch t.Keysym.Sym {
				case sdl.K_LALT, sdl.K_RALT, sdl.K_LSHIFT, sdl.K_RSHIFT, sdl.K_LCTRL, sdl.K_RCTRL:
					modifers = false
					if t.Type == sdl.KEYDOWN {
						modifers = true
						continue
					}
				}

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
						DrawRescaleIndicator()
					}
					doDraw = true

				case sdl.K_RETURN:
					fullscreen = !fullscreen
					SetFullscreen(fullscreen)
				case sdl.K_1:
					setIf(modifers, "black", &patternBg, &patternZalivka)
					doDraw = true
				case sdl.K_2:
					setIf(modifers, "white", &patternBg, &patternZalivka)
					doDraw = true
				case sdl.K_3:
					setIf(modifers, "red", &patternBg, &patternZalivka)
					doDraw = true
				case sdl.K_q:
					setIf(modifers, "checkerDark", &patternBg, &patternZalivka)
					doDraw = true
				case sdl.K_w:
					setIf(modifers, "checkerLight", &patternBg, &patternZalivka)
					doDraw = true
					// case sdl.K_6:123232
					// 	curBgColor = "grey"
					// 	doDraw = true
					// case sdl.K_7:
					// 	curBgColor = "red"
					// 	doDraw = true
					// case sdl.K_8:
					// 	curBgColor = "green"
					// 	doDraw = true
					// case sdl.K_4:
					// 	if !modifers {

					// 	} else {
					// 		curBgColor = "green"
					// 	}
					// 	zalivkaChB = true
					// 	doDraw = true
					// case sdl.K_5:
					// 	bgChB = true
					// 	doDraw = true
				}
			}
		}
	}
}

func DrawRescaleIndicator() {
	fmt.Println("is refitted##############")

	renderer.SetDrawColor(100, 255, 255, 255)
	renderer.FillRect(&sdl.Rect{300, 300, 300, 300})
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

func DrawPattern(name string) {
	// whatColor := ""
	switch name {
	case "black":
		DrawBlank(0, 0, 0)
	case "white":
		DrawBlank(255, 255, 255)
	case "grey":
		DrawBlank(172, 172, 172)
	case "red":
		DrawBlank(255, 0, 0)
	case "green":
		DrawBlank(0, 255, 0)
	case "blue":
		DrawBlank(0, 0, 255)
	case "checkerDark":
		c1 := sdl.Color{45, 45, 45, 255}
		c2 := sdl.Color{85, 85, 85, 255}
		DrawCheckerboard(c1, c2)
	case "checkerLight":
		c1 := sdl.Color{180, 180, 180, 180}
		c2 := sdl.Color{220, 220, 220, 255}
		DrawCheckerboard(c1, c2)
	}
}

func DrawBlank(r, g, b uint8) {
	renderer.SetDrawColor(r, g, b, 255)
	renderer.FillRect(&sdl.Rect{0, 0, screenWidth, screenHeight})
}

func DrawCheckerboard(color1, color2 sdl.Color) {
	// renderer.SetClipRect(&sdl.Rect{offsetX, offsetY, newWidth, newHeight})

	renderer.SetDrawColor(color1.R, color1.G, color1.B, 255)
	renderer.FillRect(&sdl.Rect{0, 0, screenWidth, screenHeight})

	renderer.SetDrawColor(color2.R, color2.G, color2.B, 255)

	newPosX := int32(0)
	newPosY := int32(0)

	squareSize := int32(8)

	chetCounter := 0
	// начать строку
	for stepY := squareSize; newPosY <= screenHeight; newPosY = newPosY + stepY {
		// fmt.Println("newPosX", newPosX)
		// fmt.Println("newPosY", newPosY)
		// fmt.Println("chetCounter", chetCounter)

		if chetCounter == 0 {
			newPosX = 0
			chetCounter++
		} else {
			newPosX = 0 + squareSize
			chetCounter--
		}
		for stepX := squareSize * 2; newPosX <= screenWidth; newPosX = newPosX + stepX {
			// fmt.Println("newPosX", newPosX)
			// fmt.Println("newPosY", newPosY)
			renderer.FillRect(&sdl.Rect{newPosX, newPosY, squareSize, squareSize})
			// fmt.Println("##end raw##")
		}
	}
}

// Draw - renders background and puts created textures in window
func Draw() {
	defer renderer.Present()
	fmt.Printf("draw start\n")

	newWidth, newHeight := Rescale(imageWidth, imageHeight)
	fmt.Println("image width", imageWidth)
	fmt.Println("image height", imageHeight)
	fmt.Println("scaled width", newWidth)
	fmt.Println("scaled height", newWidth)

	offsetX := (screenWidth - newWidth) / 2
	offsetY := (screenHeight - newHeight) / 2

	renderer.SetClipRect(&sdl.Rect{0, 0, screenWidth, screenHeight})
	DrawPattern(patternZalivka)

	renderer.SetClipRect(&sdl.Rect{offsetX, offsetY, newWidth, newHeight})
	DrawPattern(patternBg)

	renderer.SetClipRect(nil)

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

func setIf(flag bool, val string, resultTrue, resultFalse *string) {
	if flag {
		if resultTrue == nil {
			return
		}
		*resultTrue = val
	} else {
		if resultFalse == nil {
			return
		}
		*resultFalse = val
	}
}

// ############################################

func main() {

	if !Setup() {
		os.Exit(1)
	}
	filelist = ParseArgs(os.Args)
	SetFullscreen(fullscreen)
	patternZalivka = defaultZalivka
	patternBg = defaultBg
	CreateImage(getCurFile())

	HandleEvents()
	Shutdown()
}
