package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SceneSwapFunc func(*Game)

type Button struct {
	image *ebiten.Image
	rect  Rectangle
	cmd   SceneSwapFunc
}

type Signal struct {
	image      *ebiten.Image
	rect       Rectangle
	is_dragged bool
	signalType string
}

type Receptor struct {
	image              *ebiten.Image
	rect               Rectangle
	is_touching_signal bool
	receptorType       string
}

type Kinase struct {
	image         *ebiten.Image
	rect          Rectangle
	is_clicked_on bool
}

func newButton(path string, rect Rectangle, cmd SceneSwapFunc) Button {
	var btn_image, _, err = ebitenutil.NewImageFromFile(path)

	if err != nil {
		fmt.Println("Error parsing date:", err)
	}
	return Button{
		image: btn_image,
		rect:  rect,
		cmd:   cmd,
	}
}

func newSignal(path string, rect Rectangle) Signal {
	var sig_image, _, err = ebitenutil.NewImageFromFile(path)

	if err != nil {
		fmt.Println("Error parsing date:", err)
	}
	return Signal{
		image:      sig_image,
		rect:       rect,
		is_dragged: false,
	}
}

func newKinase(path string, rect Rectangle) Kinase {
	var kin_image, _, err = ebitenutil.NewImageFromFile(path)
	if err != nil {
		fmt.Println("Error parsing date:", err)
	}
	return Kinase{
		image:         kin_image,
		rect:          rect,
		is_clicked_on: false,
	}
}

func (b Button) on_click(g *Game) {
	var x_c, y_c = ebiten.CursorPosition()
	var b_pos = newVector(x_c, y_c)
	if rect_point_collision(b.rect, b_pos) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		b.cmd(g)
	}
}

func (b Button) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.rect.pos.x), float64(b.rect.pos.y))
	screen.DrawImage(b.image, op)
}

func (s *Signal) on_click(g *Game) {
	var x_c, y_c = ebiten.CursorPosition()
	var b_pos = newVector(x_c, y_c)

	if !s.is_dragged {
		if rect_point_collision(s.rect, b_pos) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			s.is_dragged = true
		}
	} else {
		s.is_dragged = false
	}

	if s.is_dragged {
		s.rect.pos = newVector(b_pos.x-s.rect.width/2, b_pos.y-s.rect.height/2)
	}
}

func (s Signal) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.rect.pos.x), float64(s.rect.pos.y))
	screen.DrawImage(s.image, op)
}

func newReceptor(path string, rect Rectangle, rtype string) Receptor {
	var rec_image, _, err = ebitenutil.NewImageFromFile(path)

	if err != nil {
		fmt.Println("Error parsing date:", err)
	}
	return Receptor{
		image:              rec_image,
		rect:               rect,
		is_touching_signal: false,
		receptorType:       rtype,
	}
}

func (r Receptor) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(r.rect.pos.x), float64(r.rect.pos.y))
	screen.DrawImage(r.image, op)
}

func (r *Receptor) update() {
	if aabb_collision(signal.rect, r.rect) {
		r.is_touching_signal = true
	} else {
		r.is_touching_signal = false
	}
}

func (k *Kinase) update() {
	var x_c, y_c = ebiten.CursorPosition()
	var b_pos = newVector(x_c, y_c)

	if !k.is_clicked_on {
		k.rect.pos = newVector(k.rect.pos.x+1, k.rect.pos.y+1)
	}
	if rect_point_collision(k.rect, b_pos) {
		k.is_clicked_on = true
	}
}
