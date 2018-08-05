package phly_img

import (
	"github.com/hackborn/phly"
	_ "image/gif"  // Initialize the gif codec
	_ "image/jpeg" // Initialize the jpg codec
	_ "image/png"  // Initialize the png codec
)

func init() {
	// Register vars
	phly.RegisterVar("srcw", "Source image width", nil)
	phly.RegisterVar("srch", "Source image height", nil)
	phly.RegisterVar("srcpath", "Source file path", nil)
	phly.RegisterVar("srcdir", "Source file path directory", nil)
	phly.RegisterVar("srcbase", "Source file name (excluding extension)", nil)
	phly.RegisterVar("srcext", "Source file extension", nil)

	// Register nodes
	phly.Register(&load{})
	phly.Register(&save{})
	phly.Register(&scale{})
}
