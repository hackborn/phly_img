package phly_img

import (
	"github.com/hackborn/phly"
)

func init() {
	phly.Register("phly/img/load", &load{})
	phly.Register("phly/img/scale", &scale{})
}
