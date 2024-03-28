package main

import (
	"bytes"
	"fmt"
	_ "image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	baseScreenWidth, baseScreenHeight = 1250, 750
)

var (
	// GLOBAL VARIABLES
	screenWidth, screenHeight = 1250, 750
	maxWidth, maxHeight       = ebiten.ScreenSizeInFullscreen()
	widthRatio                = float64(maxWidth / baseScreenWidth)
	heightRatio               = float64(maxHeight / baseScreenHeight)

	audioContext *audio.Context
	audioPlayer  *audio.Player

	err   error
	scene string = "Main Menu"

	otherToMenuButton Button
	info              string
	infoButton        InfoPage

	reset       bool
	defaultFont Font
	codonFont   Font

	seedSignal int
	template   = [5]string{}

	adenine   Nucleobase
	thymine   Nucleobase
	guanine   Nucleobase
	cytosine  Nucleobase
	uracil    Nucleobase
	aminoAcid Nucleobase
	stop      Nucleobase
)

type Game struct {
	switchedScene bool
	stateMachine  *StateMachine

	state_array []GUI

	menuSprites          []GUI
	aboutSprites         []GUI
	levSelSprites        []GUI
	receptionSprites     []GUI
	transductionSprites  []GUI
	transcriptionSprites []GUI
	translationSprites   []GUI
}

func executableDir() string {
    // Get the absolute path of the executable
    exePath, err := os.Executable()
    if err != nil {
        panic(err) // Handle error
    }
    // Get the directory of the executable
    dir := filepath.Dir(exePath)
    return dir
}

func loadFile(image string) string {
    // Get the absolute path to the Assets/Images directory
    imageDir := filepath.Join(executableDir())
    // Construct the absolute path to the image file
    imagepath := filepath.Join(imageDir, image)
    // Open file
    file, err := os.Open(imagepath)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return "Error"
    }
    defer file.Close()
    return file.Name()
}

func loadFont(font string) string {
    // Get the absolute path to the Assets/Images directory
    fontDir := filepath.Join(executableDir())
    // Construct the absolute path to the image file
    fontpath := filepath.Join(fontDir, font)
    // Open file
    file, err := os.Open(fontpath)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return "Error"
    }
    defer file.Close()
    return file.Name()
}

func loadMusic(music string) string {
    // Get the absolute path to the Assets/Images directory
   musicDir := filepath.Join(executableDir())
    // Construct the absolute path to the image file
    musicpath := filepath.Join(musicDir, music)
    // Open file
    file, err := os.Open(musicpath)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return "Error"
    }
    defer file.Close()
    return file.Name()
}

func (g *Game) init() {
	ebiten.SetWindowPosition(100, 0)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeOnlyFullscreenEnabled)

	defaultFont = newFont(loadFont("CourierPrime-Regular.ttf"), 32)
	codonFont = newFont(loadFont("BlackOpsOne-Regular.ttf"), 60)

	//	maxWidth, maxHeight       = g.Layout(maxWidth, maxHeight)

	// 	Initialize audio context
	audioContext = audio.NewContext(44100)

	g.switchedScene = true

	mp3Bytes, err := os.ReadFile(loadMusic("Signaling_of_the_Cell_MenuScreen.mp3"))
	if err != nil {
		log.Fatal(err)
	}
	mp3Stream, err := mp3.DecodeWithoutResampling(bytes.NewReader(mp3Bytes))
	if err != nil {
		log.Fatal(err)
	}
	audioPlayer, err = audioContext.NewPlayer(mp3Stream)
	if err != nil {
		log.Fatal(err)
	}

	infoButton = newInfoPage("infoButton.png", "infoPage.png", newRect(850, 0, 165, 165), "btn")
	otherToMenuButton = newButton("menuButton.png", newRect(1000, 0, 300, 200), ToMenu)
	adenine = newNucleobase("A", newRect(100, 500, 65, 150), 0, false)
	thymine = newNucleobase("T", newRect(100, 500, 65, 150), 0, false)
	guanine = newNucleobase("G", newRect(100, 500, 65, 150), 0, false)
	cytosine = newNucleobase("C", newRect(100, 500, 65, 150), 0, false)
	uracil = newNucleobase("U", newRect(100, 500, 65, 150), 0, false)
	aminoAcid = newNucleobase("X", newRect(0, 0, 60, 60), 1, false)
	stop = newNucleobase("STOP", newRect(0, 0, 60, 60), 1, false)

	var s_map = SceneConstructorMap{
		"Main Menu": newMainMenu, "About": newAbout, "Level Selection": newLevelSelection,
		"Signal Reception": newReceptionLevel, "Signal Transduction": newTransductionLevel,
		"Transcription": newTranscriptionLevel, "Translation": newTranslationLevel,
	}

	g.stateMachine = newStateMachine(s_map)
	ToMenu(g)
}

func (g *Game) Update() error {

	ebiten.SetWindowTitle("CSPS - " + scene)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ebiten.SetFullscreen(false)
	}

	g.stateMachine.update(g)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	if ebiten.IsFullscreen() {
		// Use this if statement to set sizes of graphics to fullscreen scale, else normal scale.
		if screenWidth == baseScreenWidth && screenHeight == baseScreenHeight {
			screenWidth, screenHeight = ebiten.ScreenSizeInFullscreen()
			g.stateMachine.Scale(g, screen)
			defaultFont = newFont(loadFont("CourierPrime-Regular.ttf"), 32*int(heightRatio))
			codonFont = newFont(loadFont("BlackOpsOne-Regular.ttf"), 60*int(heightRatio))
		}
	} else {
		if screenWidth != baseScreenWidth && screenHeight != baseScreenHeight {
			screenWidth, screenHeight = baseScreenWidth, baseScreenHeight
			g.stateMachine.Scale(g, screen)
			defaultFont = newFont(loadFont("CourierPrime-Regular.ttf"), 32)
			codonFont = newFont(loadFont("BlackOpsOne-Regular.ttf"), 60)
		}
	}

	if g.switchedScene {
		g.stateMachine.Scale(g, screen)
		g.switchedScene = false
	}

	g.stateMachine.draw(g, screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if ebiten.IsFullscreen() {
		return outsideWidth, outsideHeight
	} else {
		return baseScreenWidth, baseScreenHeight
	}
}

func main() {
	game := &Game{}
	game.init()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

}
