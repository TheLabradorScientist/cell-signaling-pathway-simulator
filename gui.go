package main

import (
	"fmt"
	"image/color"

	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ButtonFunc func(*Game)

type GUI interface {
	draw(screen *ebiten.Image)
	update(params ...interface{})
	scaleToScreen()
	//getStructType(GUI) reflect.Type
}

//func getStructType(g GUI) reflect.Type {
//	return reflect.TypeOf(g)
//}

// Create Sprite struct with fields for image, second image (optional),
// rectangle, scale factors, matrix draw option, + original images
type Sprite struct {
	image       *ebiten.Image
	image_2     *ebiten.Image
	rect        Rectangle
	scaleW      float64
	scaleH      float64
	op          ebiten.GeoM
	origImage   *ebiten.Image
	origImage_2 *ebiten.Image
}

type Button struct {
	//CommonDraw
	Sprite
	cmd ButtonFunc
}

type VolButton struct {
	Button
	player audio.Player
	status string
}

type Signal struct {
	Sprite
	is_dragged bool
	signalType string
}

type Receptor struct {
	Sprite
	is_touching_signal bool // MAY move outside to variable of plasma struct to let signal access it
	receptorType       string
}

type Kinase struct {
	Sprite
	is_moving     bool
	is_clicked_on bool
	delta         int
	kinaseType    string
}

type TFA struct {
	Sprite
	is_active bool
	tfaType   string
}

type Transcript struct {
	Sprite
	codon string
	isRNA bool
}

type Template struct {
	Sprite
	codon       string
	fragment    int
	is_complete bool
}

type RNAPolymerase struct {
	Sprite
	next bool
}

type Nucleobase struct {
	Sprite
	baseType   string
	index      int
	isTemplate bool
}

type CodonChoice struct {
	Sprite
	codon      string
	bases      [3]Nucleobase
	is_dragged bool
}

type tRNA struct {
	CodonChoice
	aminoAcid	Nucleobase
}

type Ribosome struct {
	Sprite
}

type Parallax struct {
	Sprite
	layer float64
}

type InfoPage struct {
	Sprite
	status string
	// Functions: when screen switches, is drawn in btn status. When mouseButtonJustPressed + btn status,
	// changes to pg status. When mouseButtonJustPressed + pg status, changes to button status.
	// update function sets image to status + "_image". Pg image should say "Click to exit."
	// if status = btn, if status = pg
}

type StillImage struct {
	Sprite
}

// Create new sprite with variadic parameters for multiple images
func newSprite(params ...interface{}) Sprite {
	// if 4 parameters passed, two images are needed
	if len(params) == 4 {
		path1 := params[0].(string) // image 1
		path2 := params[1].(string) // image 2
		rect := params[2].(Rectangle)
		scaleW, scaleH := params[3].(float64), params[3].(float64)

		// Store original image from the parameter to use for scaling in fullscreen
		var origImg, _, err1 = ebitenutil.NewImageFromFile(loadFile(path1))
		var origImg2, _, err2 = ebitenutil.NewImageFromFile(loadFile(path2))

		// Check error if image does not exist
		if err1 != nil {
			fmt.Println("Error parsing date:", err)
		}
		if err2 != nil {
			fmt.Println("Error parsing date:", err)
		}

		// Scale original image from parameter based on scaling factors
		var img_1 = scaleImage(origImg, scaleW, scaleH)
		var img_2 = scaleImage(origImg2, scaleW, scaleH)

		// Return Sprite struct
		return Sprite{
			image:       img_1,
			image_2:     img_2,
			rect:        rect,
			scaleW:      scaleW,
			scaleH:      scaleH,
			origImage:   origImg,
			origImage_2: origImg2,
		}

	} else { // if 3 parameters passed, no second image needed.
		path := params[0].(string)
		rect := params[1].(Rectangle)
		scaleW, scaleH := params[2].(float64), params[2].(float64)
		var origImg, _, err1 = ebitenutil.NewImageFromFile(loadFile(path))
		if err1 != nil {
			fmt.Println("Error parsing date:", err)
		}
		var img_1 = scaleImage(origImg, scaleW, scaleH)
		return Sprite{
			image:       img_1,
			image_2:     img_1,
			rect:        rect,
			scaleW:      scaleW,
			scaleH:      scaleH,
			origImage:   origImg,
			origImage_2: origImg,
		}
	}
}

// General function for scaling any image using the parameters for scaling factors
func scaleImage(img *ebiten.Image, scaleFactorW float64, scaleFactorH float64) *ebiten.Image {
	bounds := img.Bounds()
	width := int(float64(bounds.Dx()) * scaleFactorW)
	height := int(float64(bounds.Dy()) * scaleFactorH)
	scaled := ebiten.NewImage(width, height) // Creates empty new image with desired width/height
	ops := &ebiten.DrawImageOptions{}        // Create new DrawImageOptions to resize img from parameter
	ops.GeoM.Scale(scaleFactorW, scaleFactorH)
	scaled.DrawImage(ebiten.NewImageFromImage(img), ops) // Draws resized img onto the empty scaled image
	return scaled                                        // Returns scaled with img drawn onto new bounds
}

// Method of Sprite struct, calls scaleImage() on sprite image using
// sprite geometry (op) and scaling factors
func (s *Sprite) scaleToScreen() {
	s.op = ebiten.GeoM{}
	if ebiten.IsFullscreen() {
		s.image = scaleImage(s.origImage, 1.15*s.scaleW*float64(widthRatio), 1.2*s.scaleH*float64(heightRatio))
		//s.image = scaleImage(s.origImage, s.scaleW*float64(widthRatio), s.scaleH*float64(heightRatio))
	} else {
		s.image = scaleImage(s.origImage, s.scaleW, s.scaleH)
	}
	// NEVER TRY THIS CODE- IT BREAKS THE COMPUTER!!! - s.image = scaleImage(s.origImage, s.scaleW*float64(baseScreenWidth/screenWidth), s.scaleH*float64(baseScreenHeight/screenHeight))
}

func (s Sprite) draw(screen *ebiten.Image, params ...interface{}) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM = s.op
	op.GeoM.Translate(float64(s.rect.pos.x), float64(s.rect.pos.y))
	if len(params) == 0 {
		screen.DrawImage(s.image, op)
	}
	if len(params) == 1 {
		layer := params[0].(float64)
		scaleW := (layer + 0.5) / (layer)
		scaleH := (layer + 0.5) / (layer)
		op.GeoM.Scale(scaleW, scaleH)
		screen.DrawImage(s.image, op)
	}
}

func newStillImage(path string, rect Rectangle) StillImage {
	sprite := newSprite(path, rect, 0.5)
	return StillImage{
		Sprite: sprite,
	}
}

func (s StillImage) draw(screen *ebiten.Image) {
	s.Sprite.draw(screen)
}

func (s StillImage) update(params ...interface{}) {}

func newInfoPage(path1 string, path2 string, rect Rectangle, stat string) InfoPage {
	sprite := newSprite(path1, path2, rect, 1.0)
	return InfoPage{
		Sprite: sprite,
		status: stat,
	}
}

func (i *InfoPage) update(params ...interface{}) {
	var x_c, y_c = ebiten.CursorPosition()
	var b_pos = newVector(x_c, y_c)
	if rect_point_collision(i.rect, b_pos) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if i.status == "btn" {
			i.status = "pg"
			i.Sprite.rect = newRect(0, 0, screenWidth, screenHeight)
		} else {
			i.status = "btn"
			i.Sprite.rect = newRect(850, 0, 165, 165)
		}
		temp_img := i.Sprite.origImage
		i.Sprite.origImage = i.Sprite.origImage_2
		i.Sprite.origImage_2 = temp_img
		i.Sprite.scaleToScreen()
	}
}

func (i InfoPage) draw(screen *ebiten.Image) {
	i.Sprite.draw(screen)
	if i.status == "pg" {
		Purple := color.RGBA{50, 0, 50, 250}
		defaultFont.drawFont(screen, info, 300, 200, color.RGBA(Purple))
	}
}

func newParallax(path string, rect Rectangle, layer float64) Parallax {
	sprite := newSprite(path, rect, (layer+0.5)/(2*layer))
	return Parallax{
		Sprite: sprite,
		layer:  layer,
	}
}

func (p *Parallax) update(params ...interface{}) {
	var x_c, y_c = ebiten.CursorPosition()
	var l = int(p.layer)
	switch scene {
	case "Main Menu":
		p.rect.pos.x = -5 * (x_c + 75) / (6 * l)
		p.rect.pos.y = -5 * (y_c + 100) / (7 * l)
	case "Signal Reception":
		p.rect.pos.x = -6 * (x_c + 100) / (7 * l)
		p.rect.pos.y = -2 * (y_c + 100) / (3 * l)
	case "Signal Transduction":
		p.rect.pos.x = -5 * (x_c + 80) / (7 * l)
		p.rect.pos.y = -3 * (y_c + 100) / (5 * l)	
	case "Translation":
		p.rect.pos.x = -5 * (x_c + 95) / (7 * l)
		p.rect.pos.y = -3 * (y_c + 100) / (5 * l)
	}
	//p.rect.pos.x = (x_c - 625) / (2*l);
	//p.rect.pos.y = (y_c - 375) / (2*l);
}

func (p Parallax) draw(screen *ebiten.Image) {
	p.Sprite.draw(screen, p.layer)
}

func newButton(path string, rect Rectangle, cmd ButtonFunc) Button {
	sprite := newSprite(path, rect, 1.0)
	return Button{
		Sprite: sprite,
		cmd:    cmd,
	}
}

func (b *Button) update(params ...interface{}) {
	if len(params) > 0 {
		g, ok := params[0].(*Game)
		if !ok {
			return
		}
		var x_c, y_c = ebiten.CursorPosition()
		var b_pos = newVector(x_c, y_c)
		if rect_point_collision(b.rect, b_pos) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			b.cmd(g)
		}
	}
}

func (b Button) draw(screen *ebiten.Image) {
	b.Sprite.draw(screen)
}

func newVolButton(path string, rect Rectangle, cmd ButtonFunc, player audio.Player) VolButton {
	btn := newButton(path, rect, cmd)
	return VolButton{
		Button: btn,
		player: player,
		status: "ON",
	}
}

func (v *VolButton) update(params ...interface{}) {
	v.Button.update(params...)
	if v.status == "ON" && !v.player.IsPlaying() {
		v.player.Rewind()
		v.player.Play()
	}
}

func (v VolButton) draw(screen *ebiten.Image) {
	v.Button.draw(screen)
}

func (v *VolButton) Toggle(g *Game) {
	if v.player.IsPlaying() {
		v.SwitchVol("OFF")
	} else {
		v.SwitchVol("ON")
	}
}

func (v *VolButton) SwitchVol(onOff string) {
	v.status = onOff
	if v.status == "OFF" {
		v.player.Pause()
		sprite := newSprite("volButtonOff.png", v.rect, 1.0)
		v.Sprite = sprite
	} else {
		v.player.Play()
		sprite := newSprite("volButtonOn.png", v.rect, 1.0)
		v.Sprite = sprite
	}
}

func newSignal(path string, rect Rectangle) Signal {
	sprite := newSprite(path, rect, 0.5)

	return Signal{
		Sprite:     sprite,
		is_dragged: false,
	}
}

func (s *Signal) update(params ...interface{}) {
	x_c, y_c := ebiten.CursorPosition()
	var b_pos = newVector(x_c, y_c)

	if !s.is_dragged {
		if rect_point_collision(s.rect, b_pos) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			s.is_dragged = true
		}
	} else if s.is_dragged {
		if s.rect.pos.y <= receptionStruct.receptorB.rect.pos.y && b_pos.y <= receptionStruct.receptorB.rect.pos.y-25 {
			s.Sprite.drag(true, true, b_pos)
		} else {
			s.Sprite.drag(true, false, b_pos)
		}
	}
}

func (s *Sprite) drag(x_drag, y_drag bool, b_pos Vector) {
	if x_drag {
		s.rect.pos.x = b_pos.x - (s.rect.width / 2)
	}
	if y_drag {
		s.rect.pos.y = b_pos.y - (s.rect.height / 2)
	}
}

func (s *Signal) bind(r *Receptor) {
	s.is_dragged = false
	if r.receptorType == "receptorA" || r.receptorType == "receptorD" {
		s.rect.pos.x, s.rect.pos.y = r.rect.pos.x+80, r.rect.pos.y
	} else {
		s.rect.pos.x, s.rect.pos.y = r.rect.pos.x+60, r.rect.pos.y
	}
}

func (s Signal) draw(screen *ebiten.Image) {
	s.Sprite.draw(screen)
}

func newReceptor(path1 string, path2 string, rect Rectangle, rtype string) Receptor {
	sprite := newSprite(path1, path2, rect, 0.52)
	return Receptor{
		Sprite:             sprite,
		is_touching_signal: false,
		receptorType:       rtype,
	}
}

func (r Receptor) draw(screen *ebiten.Image) {
	r.Sprite.draw(screen)
}

func (r *Receptor) update(params ...interface{}) {
	var x_c, y_c = ebiten.CursorPosition()
	switch r.receptorType {
	case "receptorA":
		r.rect.pos.x = ((-5 * (x_c + 100) / (9 * 1)) + (screenWidth * 1 / 7)) * screenWidth / baseScreenWidth
		r.rect.pos.y = ((-1 * (y_c + 100) / (4 * 1)) + 450) * screenHeight / baseScreenHeight
	case "receptorB":
		r.rect.pos.x = ((-5 * (x_c + 100) / (9 * 1)) + (screenWidth * 4 / 7)) * screenWidth / baseScreenWidth
		r.rect.pos.y = ((-1 * (y_c + 100) / (4 * 1)) + 400) * screenHeight / baseScreenHeight
	case "receptorC":
		r.rect.pos.x = ((-5 * (x_c + 100) / (9 * 1)) + (screenWidth * 7 / 7)) * screenWidth / baseScreenWidth
		r.rect.pos.y = ((-1 * (y_c + 100) / (4 * 1)) + 400) * screenHeight / baseScreenHeight
	case "receptorD":
		r.rect.pos.x = ((-5 * (x_c + 100) / (9 * 1)) + (screenWidth * 9 / 7)) * screenWidth / baseScreenWidth
		r.rect.pos.y = ((-1 * (y_c + 100) / (4 * 1)) + 450) * screenHeight / baseScreenHeight
	}
	if aabb_collision(receptionStruct.signal.rect, r.rect) {
		r.is_touching_signal = true
	} else {
		r.is_touching_signal = false
	}
}

func (r *Receptor) animate() {
	r.Sprite.origImage = r.Sprite.origImage_2
	r.Sprite.scaleToScreen()
}

func newKinase(path1 string, path2 string, rect Rectangle, ktype string) Kinase {
	sprite := newSprite(path1, path2, rect, 0.52)
	return Kinase{
		Sprite:        sprite,
		is_moving:     false,
		is_clicked_on: false,
		delta:         3,
		kinaseType:    ktype,
	}
}

func (k *Kinase) update(params ...interface{}) {
	var x_c, y_c = ebiten.CursorPosition()
	var b_pos = newVector(x_c, y_c)
	if strings.Contains(k.kinaseType, "temp_tk1") {
		if !k.is_moving {
			var x_c, y_c = ebiten.CursorPosition()
			switch k.kinaseType {
			case "temp_tk1A":
				k.rect.pos.x = ((-5 * (x_c + 100) / (9 * 1)) + (screenWidth * 1 / 7)) * screenWidth / baseScreenWidth
				k.rect.pos.y = ((-1 * (y_c + 100) / (5 * 1)) + 650) * screenHeight / baseScreenHeight
			case "temp_tk1B":
				k.rect.pos.x = ((-5 * (x_c + 100) / (9 * 1)) + (screenWidth * 4 / 7)) * screenWidth / baseScreenWidth
				k.rect.pos.y = ((-1 * (y_c + 100) / (5 * 1)) + 600) * screenHeight / baseScreenHeight
			case "temp_tk1C":
				k.rect.pos.x = ((-5 * (x_c + 100) / (9 * 1)) + (screenWidth * 7 / 7)) * screenWidth / baseScreenWidth
				k.rect.pos.y = ((-1 * (y_c + 100) / (5 * 1)) + 600) * screenHeight / baseScreenHeight
			case "temp_tk1D":
				k.rect.pos.x = ((-5 * (x_c + 100) / (9 * 1)) + (screenWidth * 9 / 7)) * screenWidth / baseScreenWidth
				k.rect.pos.y = ((-1 * (y_c + 100) / (5 * 1)) + 650) * screenHeight / baseScreenHeight
			}
		} else if k.is_moving {
			if k.rect.pos.y <= screenHeight {
				k.descend()
			}
		}
	} else if !k.is_clicked_on && k.is_moving {
		if k.rect.pos.y <= 400*(screenHeight/750) && k.kinaseType == "tk2" {
			k.descend()
		} else if k.rect.pos.y <= 50*(screenHeight/750) && k.kinaseType == "tk1" {
			k.descend()
		} else {
			if ebiten.IsFullscreen() {
				k.rect.pos.x += k.delta * int(widthRatio)
			} else {
				k.rect.pos.x += k.delta
			}
		}
	}
	if k.kinaseType == "tk1" {
		if rect_point_collision(k.rect, b_pos) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && aabb_collision(k.rect, transductionStruct.tk2.rect) {
			k.is_clicked_on = true
		}
	} else if k.kinaseType == "tk2" {
		if rect_point_collision(k.rect, b_pos) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && aabb_collision(k.rect, transductionStruct.tfa.rect) {
			k.is_clicked_on = true
		}
	}
	if k.rect.pos.x+k.rect.width >= screenWidth {
		k.delta = -3
	} else if k.rect.pos.x <= 0 {
		k.delta = 3
	}
}

func (k Kinase) draw(screen *ebiten.Image) {
	k.Sprite.draw(screen)
}

func (k *Kinase) activate() {
	if k.kinaseType == "tk2" {
		k.rect.pos.y -= 3 * (screenHeight / baseScreenHeight)
	}
	if strings.Contains(k.kinaseType, "temp_tk1") && !k.is_moving {
		k.rect.pos.y -= 3 * (screenHeight / baseScreenHeight)
	}
	k.animate()
	k.is_moving = true
}

func (k *Kinase) descend() {
	if ebiten.IsFullscreen() {
		k.rect.pos.y += 3
	} else {
		k.rect.pos.y += 2
	}
}

func (k *Kinase) animate() {
	k.Sprite.origImage = k.Sprite.origImage_2
	k.Sprite.scaleToScreen()
}

func (t *TFA) activate() {
	if t.tfaType == "tfa1" {
		t.rect.pos.y -= 3 * (screenHeight / 750)
	}
	t.animate()
	t.is_active = true
}

func newTFA(path1 string, path2 string, rect Rectangle, tfaType string) TFA {
	sprite := newSprite(path1, path2, rect, 0.52)
	return TFA{
		Sprite:    sprite,
		is_active: false,
		tfaType:   tfaType,
	}
}

func (t *TFA) update(params ...interface{}) {
	if t.is_active {
		if t.rect.pos.y <= screenHeight && t.tfaType == "tfa1" {
			t.rect.pos.y += 2 * (screenHeight / 750)
		}
		if t.tfaType == "tfa2" {
			rnaPolymPos := transcriptionStruct.rnaPolymerase.rect.pos
			if rnaPolymPos.x >= 80 {
				t.rect.pos.x = rnaPolymPos.x + 60
				t.rect.pos.y = rnaPolymPos.y + 115
			} else if t.rect.pos.y <= 450 {
				t.rect.pos.y += 4 * (screenHeight / 750)
				t.rect.pos.x -= 2 * (screenWidth / 1250)
			}
		}
	}
}

func (t *TFA) animate() {
	t.Sprite.origImage = t.Sprite.origImage_2
	t.Sprite.scaleToScreen()
}

func (t TFA) draw(screen *ebiten.Image) {
	t.Sprite.draw(screen)
}

func newRNAPolymerase(path string, rect Rectangle) RNAPolymerase {
	sprite := newSprite(path, rect, 0.5)
	return RNAPolymerase{
		Sprite: sprite,
		next:   false,
	}
}

func (r *RNAPolymerase) update(params ...interface{}) {
	if len(params) > 0 {
		//g, ok := params[0].(*Game)
		tfaPosY := transcriptionStruct.temp_tfa.rect.pos.y
		//if !ok {
		//	return
		//}
		if tfaPosY >= 420 {
			if r.rect.pos.x <= 80 {
				r.rect.pos.y += 2 * (screenHeight / 750)
				r.rect.pos.x += 4 * (screenWidth / 1250)
			}
		}
		// Checks if current DNA codon is complete
		if r.next {
			if currentFrag == 5 && r.rect.pos.x < screenWidth+50 {
				r.rect.pos.x += 5 * (screenWidth / 1250)
				r.rect.pos.y += 3 * (screenHeight / 750)
			} else if r.rect.pos.x < (160 * (currentFrag + 1)) {
				r.rect.pos.x += 5 * (screenWidth / 1250)
			} else {
				transcriptionStruct.DNA[currentFrag].is_complete = false
				reset = true
				r.next = false
			}
		}
	}
}

func (r RNAPolymerase) draw(screen *ebiten.Image) {
	r.Sprite.draw(screen)
}

func newTranscript(path string, rect Rectangle, codon string, isRNA bool) Transcript {
	sprite := newSprite(path, rect, 0.5)
	return Transcript{
		Sprite: sprite,
		codon:  codon,
		isRNA:  isRNA,
	}
}

func (transcr Transcript) draw(screen *ebiten.Image) {
	transcr.Sprite.draw(screen)
}

func (transcr *Transcript) update(params ...interface{}) {
	if transcr.isRNA {
		if currentFrag < 5 {
			transcr.rect.pos.x = transcriptionStruct.rnaPolymerase.rect.pos.x - 750
		} else if transcriptionStruct.rnaPolymerase.rect.pos.x > 1000 {
			if transcr.rect.pos.y > -600 {
				transcr.rect.pos.y -= 4
				transcr.rect.pos.x += 2
			}
		}
	}
}

func newTemplate(path string, rect Rectangle, codon string, fragment int) Template {
	sprite := newSprite(path, rect, 0.5)
	return Template{
		Sprite:      sprite,
		codon:       codon,
		fragment:    fragment,
		is_complete: false,
	}
}

func (temp Template) draw(screen *ebiten.Image) {
	temp.Sprite.draw(screen)
}

func (temp Template) update(params ...interface{}) {}

func nextDNACodon() {
	if currentFrag < 5 {
		currentFrag++
		transcriptionStruct.rnaPolymerase.next = true
	}
}

func nextMRNACodon(g *Game) {
	if mrna_ptr < 4 {
		mrna_ptr++
		mrna[mrna_ptr].is_complete = false
		reset = true
	} else {
		ToMenu(g)
		reset = false
	}
}

func newCodonChoice(path string, rect Rectangle, codon string) CodonChoice {
	sprite := newSprite(path, rect, 0.5)
	var bases [3]Nucleobase
	for x := 0; x < len(codon); x++ {
		bases[x] = newNucleobase(string(codon[x]), newRect(100+sprite.rect.pos.x+(50*x), sprite.rect.pos.y+500, 65, 150), 0, false)
	}
	return CodonChoice{
		Sprite:     sprite,
		codon:      codon,
		is_dragged: false,
		bases:      bases,
	}
}

func (c *CodonChoice) update(params ...interface{}) {
	frag := params[0].(*Template)
	var x_c, y_c = ebiten.CursorPosition()
	var b_pos = newVector(x_c, y_c)
	if rect_point_collision(c.Sprite.rect, b_pos) {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			c.is_dragged = true
		} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			c.is_dragged = false
			if len(params) == 2 {
				if aabb_collision(c.rect, translationStruct.ribosome.rect) && c.codon == transcribe(frag.codon) {
					frag.is_complete = true
				}
			} else if len(params) == 1 {
				if aabb_collision(c.rect, transcriptionStruct.rnaPolymerase.rect) && c.codon == transcribe(frag.codon) {
					frag.is_complete = true
				}
			}
		}
	}
	if c.is_dragged {
		c.Sprite.drag(true, true, b_pos)
	}
	for x := 0; x < len(c.bases); x++ {
		c.bases[x].baseType = string(c.codon[x])
		if len(params) == 2 {
			c.bases[x].rect.pos.x = 75 + c.Sprite.rect.pos.x + (50 * x)
			c.bases[x].rect.pos.y = c.Sprite.rect.pos.y + 300
		} else {
			c.bases[x].rect.pos.x = 100 + c.Sprite.rect.pos.x + (50 * x)
			c.bases[x].rect.pos.y = c.Sprite.rect.pos.y + 125
		}
		switch c.bases[x].baseType {
		case "A":
		c.bases[x].Sprite.image = adenine.image
		case "T":
		c.bases[x].Sprite.image = thymine.image
		case "G":
		c.bases[x].Sprite.image = guanine.image
		case "C":
		c.bases[x].Sprite.image = cytosine.image
		case "U":
		c.bases[x].Sprite.image = uracil.image
		}
	}
}

func (c CodonChoice) draw(screen *ebiten.Image) {
	c.Sprite.draw(screen)
	for _, base := range c.bases {
		base.draw(screen)
	}
}

func (c *CodonChoice) reset(index, y_pos int, newBases string) {
	c.rect.pos = newVector(spots[index], y_pos)
	c.codon = newBases
}

func newTRNA(path string, rect Rectangle, codon string, amino string) tRNA {
	codonChoice := newCodonChoice(path, rect, transcribe(codon))
	aminoAcid := newNucleobase(amino, codonChoice.rect, 1, true)
	return tRNA{
		CodonChoice: codonChoice,
		aminoAcid: 	 aminoAcid,
	}
}

func (t *tRNA) update(params ...interface{}) {
	t.CodonChoice.update(params[0], true)
	t.aminoAcid.rect.pos.x = t.rect.pos.x+25
	t.aminoAcid.rect.pos.y = t.rect.pos.y-25
}

func (t tRNA) draw(screen *ebiten.Image) {
	t.CodonChoice.draw(screen)
	t.aminoAcid.draw(screen)
}

func (t *tRNA) reset(index, y_pos int, newBases string, newAminoAcid string) {
	t.rect.pos = newVector(spots[index], y_pos)
	t.codon = transcribe(newBases)
	t.aminoAcid.baseType = newAminoAcid
	if t.aminoAcid.baseType == "STOP" {
		t.aminoAcid.Sprite.image = stop.image 
	} else {
		t.aminoAcid.Sprite.image = aminoAcid.image
	}
}

func newRibosome(path string, rect Rectangle) Ribosome {
	sprite := newSprite(path, rect, 0.5)
	return Ribosome{
		Sprite: sprite,
	}
}

// Updates movement of ribosome
func (ribo *Ribosome) update(params ...interface{}) {
	if len(params) > 0 {
		g, ok := params[0].(*Game)
		if !ok {
			return
		}
		if ribo.rect.pos.x <= 40 {
			ribo.rect.pos.y += 2 * (screenHeight / 750)
			ribo.rect.pos.x += 4 * (screenWidth / 1250)
		}
		// Checks if current mRNA codon is complete
		if mrna[mrna_ptr].is_complete {
			if mrna_ptr == 4 && ribo.rect.pos.x < screenWidth+50 {
				ribo.rect.pos.x += 5 * (screenWidth / 1250)
				ribo.rect.pos.y += 3 * (screenHeight / 750)
			} else if ribo.rect.pos.x < (160 * (mrna_ptr + 1)) {
				ribo.rect.pos.x += 5 * (screenWidth / 1250)
			} else {
				nextMRNACodon(g)
			}
		}
	}
}

func (ribo Ribosome) draw(screen *ebiten.Image) {
	ribo.Sprite.draw(screen)
}

func newNucleobase(btype string, rect Rectangle, index int, isTemp bool) Nucleobase {
	var path string
	if btype != "A" && btype != "T" && btype != "G" && btype != "C" && btype != "U" && btype != "N/A" && btype != "STOP" {
		path = nucleobaseImages["amino"]
	} else {path = nucleobaseImages[btype]}
	sprite := newSprite(path, rect, 0.48)
	if !isTemp {
		sprite.op.Rotate(-3.14)
	}
	return Nucleobase{
		Sprite:     sprite,
		baseType:   btype,
		index:      index,
		isTemplate: isTemp,
	}
}

var nucleobaseImages = map[string]string{
	"A":   "adenine.png",
	"T":   "thymine.png",
	"G":   "guanine.png",
	"C":   "cytosine.png",
	"U":   "uracil.png",
	"N/A": "empty.png",
	"STOP": "empty.png",
	"amino": "aminoAcid.png", // This is CS. So I can break the rules of biology. :)
}

func (n Nucleobase) draw(screen *ebiten.Image) {
	n.Sprite.draw(screen)
	if !n.isTemplate {
		if n.baseType != "N/A" {
			codonFont.drawFont(screen, n.baseType, n.rect.pos.x-50, n.rect.pos.y-50, color.Black)
		}
	} else {
		if n.baseType != "A" && n.baseType != "T" && n.baseType != "G" && n.baseType != "C" && n.baseType != "U" {
			codonFont.drawFont(screen, n.baseType, n.rect.pos.x, n.rect.pos.y+50, color.Black)
		} else {codonFont.drawFont(screen, n.baseType, n.rect.pos.x, n.rect.pos.y+100, color.Black)}
	}

}

func (n *Nucleobase) update(params ...interface{}) {
	n.rect.pos.x = (675 + transcriptionStruct.RNA[currentFrag].rect.pos.x + (50 * n.index)) - 150*(currentFrag-1)
	n.rect.pos.y = (transcriptionStruct.RNA[5].rect.pos.y + 400 + (25 * n.index)) - 75*(currentFrag-1)
}

func (n Nucleobase) String() string {
	return n.baseType
}
