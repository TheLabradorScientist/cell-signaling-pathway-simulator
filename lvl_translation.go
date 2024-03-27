package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	mrna_ptr = 0
	mrna     	[5]Template
	protein  	[5]Transcript	
	mRNAbases   [15]Nucleobase
)

type TranslationLevel struct {
	// TRANSLATION SPRITES
	protoCytoBg_2     StillImage
	cytoBg_2          Parallax
	cytoNuc_2         Parallax
	ribosome          Ribosome
	mRNA			  [5]Template
	rightTrna         tRNA
	wrongTrna1        tRNA
	wrongTrna2        tRNA
	infoButton        InfoPage
	otherToMenuButton Button
	message           string
}

var translationStruct *TranslationLevel

func newTranslationLevel(g *Game) {
	if len(g.translationSprites) == 0 {
		translationStruct = &TranslationLevel{
			protoCytoBg_2: newStillImage("CytoBg2.png", newRect(0, 0, 1250, 750)),
			cytoBg_2:      newParallax("ParallaxCyto2.png", newRect(100, 100, 1250, 750), 4),
			cytoNuc_2:     newParallax("ParallaxCyto2.5.png", newRect(100, 100, 1250, 750), 3),

			ribosome: newRibosome("ribosome.png", newRect(-200, 50, 300, 330)),
			mRNA: mrna,
			message: 
				"FINALLY, BACK TO THE CYTOPLASM! \n" +
				"Drag the tRNA with the corresponding \n" +
				"amino acid to the ribosome to \n" +
				"synthesize your protein!!!!",
		}

		translationStruct.rightTrna = newTRNA("tRNA.png", newRect(100, 450, 140, 200), mrna[0].codon, translate(mrna[0].codon))
		randomCodon1 := randomRNACodon(translationStruct.rightTrna.codon)
		translationStruct.wrongTrna1 = newTRNA("tRNA.png", newRect(400, 450, 140, 200), randomCodon1, translate(randomCodon1))
		randomCodon2 := randomRNACodon(translationStruct.rightTrna.codon)
		translationStruct.wrongTrna2 = newTRNA("tRNA.png", newRect(700, 450, 140, 200), randomCodon2, translate(randomCodon2))
		translationStruct.infoButton = infoButton
		translationStruct.otherToMenuButton = otherToMenuButton

		g.translationSprites = []GUI{
			&translationStruct.protoCytoBg_2, &translationStruct.cytoBg_2, &translationStruct.cytoNuc_2,
			&translationStruct.ribosome, &translationStruct.mRNA[0],
			&translationStruct.rightTrna, &translationStruct.wrongTrna1,
			&translationStruct.wrongTrna2, &translationStruct.otherToMenuButton,
			&translationStruct.infoButton, 
		}
	}
	g.stateMachine.state = translationStruct
}

func (t *TranslationLevel) Init(g *Game) {
	mrna_ptr = 0
	for x := 0; x < len(mRNAbases); x++ {
		base := string(mrna[x/3].codon[x%3])
		posX := mrna[2].rect.pos.x + (50 * x)
		posY := mrna[2].rect.pos.y
		mRNAbases[x] = newNucleobase(base, newRect(posX, posY, 65, 150), x, true)
	}
	t.ResetChoices()
	g.state_array = g.translationSprites
}

func (t *TranslationLevel) ResetChoices() {
	curr := &mrna[mrna_ptr]
	rand.Shuffle(len(spots), func(i, j int) {spots[i], spots[j] = spots[j], spots[i]})
	t.rightTrna.reset(0, 450, curr.codon, translate(curr.codon))
	randomCodon1 := randomRNACodon(t.rightTrna.codon)
	t.wrongTrna1.reset(1, 450, randomCodon1, translate(randomCodon1))
	randomCodon2 := randomRNACodon(t.rightTrna.codon)
	t.wrongTrna2.reset(2, 450, randomCodon2, translate(randomCodon2))
	reset = false
}

func (t *TranslationLevel) Update(g *Game) {
	t.cytoBg_2.update()
	t.cytoNuc_2.update()
	t.otherToMenuButton.update(g)
	t.infoButton.update()

	curr := &mrna[mrna_ptr]

	if reset {t.ResetChoices()}

	t.rightTrna.update(curr)
	t.wrongTrna1.update(curr)
	t.wrongTrna2.update(curr)

	t.ribosome.update(g)
}

func (t *TranslationLevel) Draw(g *Game, screen *ebiten.Image) {
	t.protoCytoBg_2.draw(screen)
	t.cytoBg_2.draw(screen)
	t.cytoNuc_2.draw(screen)

	t.mRNA[0].draw(screen)

	t.otherToMenuButton.draw(screen)

	if t.ribosome.rect.pos.x >= 42 {
		// Draw amino acids before ribosome moves without drawing amino acid for STOP.
		for x := 0; x <= mrna_ptr; x++ {
			if x < 4 {
				protein[x].draw(screen)
				codonFont.drawFont(screen, protein[x].codon, protein[x].rect.pos.x, protein[x].rect.pos.y+25, color.Black)
			}
		}
	}

	t.ribosome.draw(screen)

	t.rightTrna.draw(screen)
	t.wrongTrna1.draw(screen)
	t.wrongTrna2.draw(screen)

	if mrna_ptr != -1 {
		if mrna_ptr == 0 {
			mRNAbases[0].draw(screen)
			mRNAbases[1].draw(screen)
			mRNAbases[2].draw(screen)
		} else {
			for y := 0; y < 3; y++ {
				mRNAbases[(mrna_ptr*3)+y].draw(screen)
			}
		}
	}

	defaultFont.drawFont(screen, t.message, 75, 50, color.Black)

	t.infoButton.draw(screen)
}
