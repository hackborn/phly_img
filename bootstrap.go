package phly_img

import (
	"github.com/hackborn/phly"
	_ "image/gif"  // Initialize the gif codec
	_ "image/jpeg" // Initialize the jpg codec
	_ "image/png"  // Initialize the png codec
)

func init() {
	phly.Register(&load{})
	phly.Register(&save{})
	phly.Register(&scale{})
}
