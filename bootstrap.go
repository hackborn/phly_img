package phly_img

import (
	"github.com/hackborn/phly"
	_ "image/gif"  // Initialize the gif codec
	_ "image/jpeg" // Initialize the jpg codec
	_ "image/png"  // Initialize the png codec
)

func init() {
	// Register vars
	phly.RegisterVar("srcw", "Source image width", "")
	phly.RegisterVar("srch", "Source image height", "")
	phly.RegisterVar("srcpath", "Source file path", "")
	phly.RegisterVar("srcdir", "Source file path directory", "")
	phly.RegisterVar("srcbase", "Source file name (excluding extension)", "")
	phly.RegisterVar("srcext", "Source file extension", "")

	// Register nodes
	phly.Register(&load{})
	phly.Register(&save{})
	phly.Register(&scale{})
}
