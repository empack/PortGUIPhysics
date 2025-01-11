package data

import (
	"fyne.io/fyne/v2/data/binding"
)

type Parameter struct {
	check      binding.Bool
	name       binding.String
	defaultVal binding.Float
	val        binding.Float
	min        binding.Float
	max        binding.Float
}
