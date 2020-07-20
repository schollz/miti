package main

import (
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	mutateNotes := []string{"C", "B", "B2", "B3", "D", "D3", "D4", "G", "G3", "G4"}
	lines := []string{
		"C3 E G C E C G E",
		"C3 E G C E C G E",
		"C3 E G C E C G E",
		"C3 E G C E C G E",
		"B2 E G B E B G E",
		"B2 E G B E B G E",
		"B2 E G B E B G E",
		"B2 E G B E B G E",
	}
	notes := [][]string{}
	originalNotes := [][]string{}
	for _, line := range lines {
		notes = append(notes, strings.Fields(line))
		originalNotes = append(originalNotes, strings.Fields(line))
	}

	first := true
	for {
		newFile := `pattern a
tempo 240
instruments op-1
legato 90
C3EG-
C3EG-
C3EG-
C3EG
E3GBE-
E3GBE
E3GBD-
E3GBD
	
instruments nts-1	
legato 1
`
		time.Sleep(2 * time.Second)
		for i, line := range notes {
			for j := range line {
				if rand.Float64() < 0.1 && j > 0 && !first {
					notes[i][j] = mutateNotes[rand.Intn(len(mutateNotes))]
				}
				if rand.Float64() < 0.2 && j > 0 {
					notes[i][j] = originalNotes[i][j]
				}
				newFile += notes[i][j] + " "
			}
			newFile += "\n"
		}
		ioutil.WriteFile("arp2.miti", []byte(newFile), 0644)
		if first {
			time.Sleep(5 * time.Second)
		}
		first = false
	}
}
