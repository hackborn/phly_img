package phly_img

import (
	"github.com/hackborn/phly"
	"image"
	"os"
)

const (
	load_stringinput = "file"
	load_imgoutput   = "out"
)

// load struct loads an image for each valid filename.
type load struct {
}

func (n *load) Describe() phly.NodeDescr {
	descr := phly.NodeDescr{Id: "phly/img/load", Name: "Load Image"}
	descr.InputPins = append(descr.InputPins, phly.PinDescr{Name: load_stringinput, Purpose: "Supply file names to load."})
	descr.OutputPins = append(descr.OutputPins, phly.PinDescr{Name: load_imgoutput, Purpose: "The loaded images, one for each file name input."})
	return descr
}

func (n *load) Instantiate(args phly.InstantiateArgs, cfg interface{}) (phly.Node, error) {
	return &load{}, nil
}

func (n *load) Process(args phly.ProcessArgs, stage phly.NodeStage, input phly.Pins, output phly.NodeOutput) error {
	var err error
	doc := &phly.Doc{MimeType: MimeTypeImagePhly}
	phly.WalkStringItems(input, load_stringinput, func(channel string, src *phly.Doc, index int, item string) {
		err = phly.MergeErrors(err, loadFile(args.Filename(item), doc))
	})
	if len(doc.Items) > 0 {
		output.SendPins(phly.PinBuilder{}.Add(load_imgoutput, doc).Pins())
	}
	// Run once and we're done
	output.SendMsg(phly.MsgFromStop(nil))
	return err
}

func (n *load) StopNode(args phly.StoppedArgs) error {
	return nil
}

// loadFile() loads the given filename to an image, placing it in the doc.
func loadFile(filename string, dst *phly.Doc) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	dst.AppendItem(&PhlyImage{Img: img, SourceFile: filename})
	return nil
}
