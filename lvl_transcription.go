package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	currentFrag = 0
	rna         [5]Transcript
	dna         [5]Template
	spots       = [3]int{350, 650, 950}
)

type TranscriptionLevel struct {
	// TRANSCRIPTION SPRITES
	nucleusBg         StillImage
	temp_tfa          TFA
	rnaPolymerase     RNAPolymerase
	RNA               [6]Transcript
	DNA               [6]Template
	DNAbases          [15]Nucleobase
	origRNAbases      [18]Nucleobase // Dummy list containing positions and bases, accessed by RNA bases
	RNAbases          [18]Nucleobase // List that is actually drawn onto screen and updated.
	rightChoice       CodonChoice
	wrongChoice1      CodonChoice
	wrongChoice2      CodonChoice
	infoButton        InfoPage
	otherToMenuButton Button
	message           string

	// Note to self: when updating DNA image, make the sprite like plasma membrane
	// So it can scroll to the left and show different codons, with bases as separate sprites
	// Also maybe try making RNA with theta and scrolling off to a upper-right diagonal

}

var transcriptionStruct *TranscriptionLevel

func newTranscriptionLevel(g *Game) {
	if len(g.transcriptionSprites) == 0 {
		transcriptionStruct = &TranscriptionLevel{
			nucleusBg: newStillImage("NucleusBg.png", newRect(0, 0, 1250, 750)),

			temp_tfa:      newTFA("inact_TFA.png", "act_TFA.png", newRect(420, -100, 150, 150), "tfa2"),
			rnaPolymerase: newRNAPolymerase("rnaPolym.png", newRect(-400, 100, 340, 265)),
			message: "WELCOME TO THE NUCLEUS! \n" +
				"Drag the complementary RNA codon \n" +
				"to RNA Polymerase to transcribe \n" +
				"a new mRNA molecule!!!",
		}

		copy(transcriptionStruct.DNA[:], dna[:])
		transcriptionStruct.DNA[5] = dna[4]
		copy(transcriptionStruct.RNA[:], rna[:])
		transcriptionStruct.RNA[5] = newTranscript("RNA5.png", transcriptionStruct.RNA[4].rect, transcriptionStruct.RNA[4].codon, true)

		transcriptionStruct.rightChoice = newCodonChoice("codonButton.png", newRect(100, 200, 192, 111), transcribe(dna[0].codon))
		transcriptionStruct.wrongChoice1 = newCodonChoice("codonButton.png", newRect(400, 200, 192, 111), randomRNACodon(transcriptionStruct.rightChoice.codon))
		transcriptionStruct.wrongChoice2 = newCodonChoice("codonButton.png", newRect(700, 200, 192, 111), randomRNACodon(transcriptionStruct.rightChoice.codon))
		transcriptionStruct.infoButton = infoButton
		transcriptionStruct.otherToMenuButton = otherToMenuButton

		g.transcriptionSprites = []GUI{
			&transcriptionStruct.nucleusBg, &transcriptionStruct.temp_tfa,
			&transcriptionStruct.DNA[0],
			&transcriptionStruct.RNA[1], &transcriptionStruct.RNA[2],
			&transcriptionStruct.RNA[3], &transcriptionStruct.RNA[4],
			&transcriptionStruct.rnaPolymerase, &transcriptionStruct.rightChoice,
			&transcriptionStruct.wrongChoice1, &transcriptionStruct.wrongChoice2,
			&transcriptionStruct.otherToMenuButton, &transcriptionStruct.infoButton,
		}
	}
	g.stateMachine.state = transcriptionStruct
}

func (t *TranscriptionLevel) Init(g *Game) {
	currentFrag = 0
	t.origRNAbases[0] = newNucleobase("N/A", newRect(0, 0, 65, 150), 0, false)
	t.origRNAbases[1] = newNucleobase("N/A", newRect(0, 0, 65, 150), 1, false)
	t.origRNAbases[2] = newNucleobase("N/A", newRect(0, 0, 65, 150), 2, false)
	for x := 0; x < len(t.origRNAbases)-3; x++ {
		base := string(t.RNA[x/3].codon[x%3])
		posX := 125 + t.RNA[currentFrag].rect.pos.x + (50 * x)
		posY := 250 + t.RNA[currentFrag].rect.pos.y + 220 - (15 * x)
		t.origRNAbases[x+3] = newNucleobase(base, newRect(posX, posY, 65, 150), x, false)
	}
	t.RNAbases = t.origRNAbases
	for x := 0; x < len(t.DNAbases); x++ {
		base := string(t.DNA[x/3].codon[x%3])
		posX := t.DNA[2].rect.pos.x + (50 * x)
		posY := t.DNA[2].rect.pos.y
		t.DNAbases[x] = newNucleobase(base, newRect(posX, posY, 65, 150), x, true)
	}
	t.ResetChoices()
	g.state_array = g.transcriptionSprites
	t.temp_tfa.activate()
}

func (t *TranscriptionLevel) ResetChoices() {
	curr := &t.DNA[currentFrag]
	rand.Shuffle(len(spots), func(i, j int) { spots[i], spots[j] = spots[j], spots[i] })
	t.rightChoice.reset(0, 600, transcribe(curr.codon))
	t.wrongChoice1.reset(1, 600, randomRNACodon(t.rightChoice.codon))
	t.wrongChoice2.reset(2, 600, randomRNACodon(t.rightChoice.codon))
	for x := 0; x < (currentFrag+1)*3; x++ {
		temp := (currentFrag+1)*3 - 1 - x
		base := t.RNAbases[x]
		base.rect.pos.y = t.origRNAbases[temp].rect.pos.y
		t.RNAbases[x] = base
	}
	reset = false
}

func (t *TranscriptionLevel) Update(g *Game) {
	t.otherToMenuButton.update(g)
	t.RNA[currentFrag].update()
	if currentFrag < 5 {
		t.RNA[currentFrag+1].update()
	}
	t.infoButton.update()
	t.temp_tfa.update()
	t.rnaPolymerase.update(g)

	curr := &t.DNA[currentFrag]

	if reset {
		t.ResetChoices()
	}

	//fmt.Printf("%t\n", dna[currentFrag].is_complete)
	t.rightChoice.update(curr)
	t.wrongChoice1.update(curr)
	t.wrongChoice2.update(curr)

	for i, base := range t.RNAbases {
		base.update()
		t.RNAbases[i] = base
	}
	// Checks if current DNA codon is complete
	if t.DNA[currentFrag].is_complete && !t.rnaPolymerase.next {
		nextDNACodon()
	}
	if t.RNA[5].rect.pos.y <= -600 {
		ToCyto2(g)
		reset = false
	}
}

func (t *TranscriptionLevel) Draw(g *Game, screen *ebiten.Image) {
	t.nucleusBg.draw(screen)

	t.RNA[currentFrag].draw(screen)
	t.rnaPolymerase.draw(screen)
	t.DNA[0].draw(screen)
	t.temp_tfa.draw(screen)
	//codonFont.drawFont(screen, strings.Join(template[0:5], ""), dna[currentFrag].rect.pos.x+300, dna[currentFrag].rect.pos.y, color.Black)

	t.rightChoice.draw(screen)
	t.wrongChoice1.draw(screen)
	t.wrongChoice2.draw(screen)

	for y := 0; y < (currentFrag+1)*3; y++ {
		t.RNAbases[y].draw(screen)
	}

	if currentFrag != -1 {
		if currentFrag == 0 {
			t.DNAbases[0].draw(screen)
			t.DNAbases[1].draw(screen)
			t.DNAbases[2].draw(screen)
		} else if currentFrag < 5 {
			for y := 0; y < 3; y++ {
				t.DNAbases[(currentFrag*3)+y].draw(screen)
			}
		}
	}

	//if currentFrag != -1 {
	//	codonFont.drawFont(screen, dna[currentFrag].codon, dna[0].rect.pos.x+500, dna[0].rect.pos.y+50, color.Black)
	//}

	defaultFont.drawFont(screen, t.message, 75, 50, color.Black)

	t.otherToMenuButton.draw(screen)

	t.infoButton.draw(screen)
}
