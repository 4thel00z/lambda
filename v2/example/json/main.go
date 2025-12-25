package main

import (
	"fmt"

	λ "github.com/4thel00z/lambda/v2"
)

type MagicSpell struct {
	Name  string `json:"name"`
	Power int    `json:"power"`
}

func main() {
	// This assumes magic.json exists in the working directory.
	spell := λ.FromJSON[MagicSpell](λ.Open("magic.json").Slurp()).Must()
	fmt.Println(spell.Name, spell.Power)
}


