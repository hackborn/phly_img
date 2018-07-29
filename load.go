package phly_img

import (
	"fmt"
	"github.com/hackborn/phly"
)

// load struct node images from a filename.
type load struct {
	File string `json:"file,omitempty"`
}

func (n *load) Run(args phly.RunArgs, input, output phly.Pins) error {
	fn := n.File
	fmt.Println("load file", fn)
	doc := &phly.Doc{}
	output.Add("0", doc)
	return nil
}

func (n *load) Instantiate(cfg interface{}) (phly.Node, error) {
	return &load{}, nil
}
