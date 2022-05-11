package screen

import "testbot/cmp"

type Screen struct {
	Name       string
	Id         string
	Access     bool
	Components []cmp.Component
}
