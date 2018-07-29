package phly_img

import (
	"github.com/hackborn/phly"
)

// scale node scales images.
type scale struct {
	Abs Sizei `json:"abs,omitempty"`
	Rel Sizef `json:"rel,omitempty"`
}

func (n *scale) Run(args phly.RunArgs, input, output phly.Pins) error {
	return nil
}

func (n *scale) Instantiate(cfg interface{}) (phly.Node, error) {
	return &scale{}, nil
}
