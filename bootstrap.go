package phly_img

import (
	"github.com/hackborn/phly"
	_ "image/gif"  // Initialize the gif codec
	_ "image/jpeg" // Initialize the jpg codec
	_ "image/png"  // Initialize the png codec
)

func init() {
	phly.Register("phly/img/load", &load{})
	phly.Register("phly/img/save", &save{})
	phly.Register("phly/img/scale", &scale{})
}
