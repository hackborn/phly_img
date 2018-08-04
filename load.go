package phly_img

import (
	"github.com/hackborn/phly"
	"image"
	"os"
)

const (
	file_stringinput = "file"
	file_imgoutput   = "out"
)

// load struct node images from a filename.
type load struct {
}

func (n *load) Describe() phly.NodeDescr {
	descr := phly.NodeDescr{Id: "phly/img/load", Name: "Load Image"}
	descr.InputPins = append(descr.InputPins, phly.PinDescr{Name: file_stringinput, Purpose: "Supply file names to load."})
	descr.OutputPins = append(descr.OutputPins, phly.PinDescr{Name: file_imgoutput, Purpose: "The loaded images, one for each file name input."})
	return descr
}

func (n *load) Instantiate(args phly.InstantiateArgs, cfg interface{}) (phly.Node, error) {
	return &load{}, nil
}

func (n *load) Run(args phly.RunArgs, input, output phly.Pins) error {
	var err error
	for _, doc := range input.Get(file_stringinput) {
		err = phly.MergeErrors(err, n.runDoc(args, doc, output))
	}
	return err
}

// runDoc() iterates the docs, pages and items, translating each filename into an image.
func (n *load) runDoc(args phly.RunArgs, srcdoc *phly.Doc, output phly.Pins) error {
	if srcdoc == nil {
		return phly.MissingDocErr
	}
	var err error
	dstdoc := &phly.Doc{MimeType: MimeTypeImagePhly}
	for _, page := range srcdoc.Pages {
		dstpage := &phly.Page{}
		for _, _file := range page.Items {
			file, ok := _file.(string)
			if !ok {
				return phly.BadRequestErr
			}
			err = phly.MergeErrors(err, n.loadFile(args.Filename(file), dstpage))
		}
		if len(dstpage.Items) > 0 {
			dstdoc.AddPage(dstpage)
		}
	}
	if len(dstdoc.Pages) > 0 {
		output.Add(file_imgoutput, dstdoc)
	}
	return err
}

func (n *load) loadFile(name string, page *phly.Page) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	page.AddItem(&PhlyImage{Img: img, SourceFile: name})
	return nil
}
