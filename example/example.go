//Margery Sprite Sheet from https://opengameart.org/content/margery-limited (CC0 2015)
package main

import(
	"io/ioutil"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	
	"github.com/RileySun/Sprite"
)

var window fyne.Window
var exSprite *sprite.Sprite
var img *canvas.Image
var contentDisplay *fyne.Container
var prevButton, nextButton, mirrorButton, loopButton, reverseButton *widget.Button
var playing bool = false

func init() {
	spriteData := getFile("Margery.png")
	//spriteSheet []byte, frameTotal int, rows int, columns int, frameWidth int, frameHeight int
	exSprite = sprite.NewSprite(spriteData, 17, 3, 6, 104, 112)
	exSprite.OnUpdate = func() {refreshSingle()}
	//name string, sprite *Sprite, startFrame int, endFrame int (remember -1 for frames, zero indexed)
	exSprite.AddCycle(sprite.NewCycle("Idle", exSprite, 12, 16, 3))
	exSprite.AddCycle(sprite.NewCycle("Run", exSprite, 0, 5, 5))
	exSprite.AddCycle(sprite.NewCycle("Jump", exSprite, 6, 11, 5))
	exSprite.SetCycle("Idle")
}


////

func main() {
	app := app.New()	
	window = app.NewWindow("Sprite Example")
	
	content := render()
	
	window.SetContent(content)
	window.CenterOnScreen()
	window.Resize(fyne.NewSize(400, 600))
	window.SetFixedSize(true)
	window.SetMaster()
	window.ShowAndRun()
}

func switchSprite(name string) {
	exSprite.SetCycle(name)
}

//Render
func render() *fyne.Container {
	sprite1 := widget.NewButton("Idle", func() {switchSprite("Idle")})
	sprite2 := widget.NewButton("Run", func() {switchSprite("Run")})
	sprite3 := widget.NewButton("Jump", func() {switchSprite("Jump")})
	spriteButtons := container.New(layout.NewGridLayout(3), sprite1, sprite2, sprite3)
	spriteButtonsMax := container.New(layout.NewMaxLayout(), spriteButtons)
	
	prevButton = widget.NewButton("Prev", prev)
	nextButton = widget.NewButton("Next", next)
	mirrorButton = widget.NewButton("Mirror", mirror)
	reverseButton = widget.NewButton("Reverse", reverse)
	loopButton = widget.NewButton("Loop", loop)
	stop := widget.NewButton("Stop", stop)
	
	buttons := container.New(layout.NewGridLayout(3), prevButton, nextButton, mirrorButton, reverseButton, loopButton, stop)
	buttonMax := container.New(layout.NewMaxLayout(), buttons)
	
	contentDisplay = container.New(layout.NewMaxLayout(), renderSingle())
	return container.NewBorder(spriteButtonsMax, buttonMax, nil, nil, contentDisplay)
}

func renderSingle() *canvas.Image {
	res := fyne.NewStaticResource("Sprite", exSprite.Image)
	return canvas.NewImageFromResource(res)
}

func refreshSingle() {
	contentDisplay.RemoveAll()
	contentDisplay.Add(renderSingle())
	contentDisplay.Refresh()
}

func disableButtons() {
	prevButton.Disable()
	nextButton.Disable()
	loopButton.Disable()
	reverseButton.Disable()
}

func enableButtons() {
	prevButton.Enable()
	nextButton.Enable()
	loopButton.Enable()
	reverseButton.Enable()
}

//Replace
func getFile(path string) []byte {
	f, _ := ioutil.ReadFile(path)
	return f
}

func prev() {
	exSprite.Cycle.Prev()
	refreshSingle()
}

func next() {
	exSprite.Cycle.Next()
	refreshSingle()
}

func mirror() {
	if exSprite.Cycle.Mirror {
		exSprite.Cycle.Mirror = false
	} else {
		exSprite.Cycle.Mirror = true
	}
	exSprite.Refresh()
	refreshSingle()
}

func stop() {
	exSprite.Stop()
	enableButtons()
}

func reverse() {
	disableButtons()
	exSprite.Cycle.Loop = true
	exSprite.Cycle.Reverse = true
	exSprite.Play()
}

func loop() {
	disableButtons()
	exSprite.Cycle.Loop = true
	exSprite.Cycle.Reverse = false
	exSprite.Play()
}