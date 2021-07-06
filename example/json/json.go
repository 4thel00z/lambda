package main

import (
	"fmt"
	λ "github.com/4thel00z/lambda/v1"
	"strings"
)

type MagicSpell struct {
	Name        string  `json:"name"`
	AttackPower float64 `json:"attack_power"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
}

func main() {
	var (
		m MagicSpell
	)
	λ.Open("magic.json").Slurp().JSON(&m).Catch(λ.Die)

	fmt.Println(strings.Join([]string{m.Name, m.Type, fmt.Sprintf("%f", m.AttackPower), m.Description}, "\n"))

	// ToJSON() detects if the current value is a pointer or not
	fmt.Println(λ.WrapValue(m).ToJSON().UnwrapString())
	// Works even if you use the pointer operator again
	fmt.Println(λ.WrapValue(&m).ToJSON().UnwrapString())
}
