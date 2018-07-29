package phly_img

import (
	"image"
	"mime"
)

const (
	MimeTypeImagePhly = "image/phly"
)

var (
	txttype = mime.TypeByExtension(".txt")
)

// --------------------------------
// PHLY-IMAGE

// PhlyImage combines a raw Go image with meta info.
type PhlyImage struct {
	Img        image.Image
	SourceFile string
}
