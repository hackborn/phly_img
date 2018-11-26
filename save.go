package phly_img

import (
	"fmt"
	"github.com/hackborn/phly"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

const (
	save_imginput  = "in"
	save_imgoutput = "out"
)

// save struct saves image items.
type save struct {
	File string `json:"file,omitempty"`
}

func (n *save) Describe() phly.NodeDescr {
	descr := phly.NodeDescr{Id: "phly/img/save", Name: "Save Image", Purpose: "Save images to a file."}
	descr.Cfgs = append(descr.Cfgs, phly.CfgDescr{Name: "file", Purpose: "The name of the saved file. Allows variables `${src}` (the source file path), `${srcdir}` (the source directory), `${srcbase}` (the source filename base, minus the extension) and  `${srcext}` (the source extension)."})
	descr.InputPins = append(descr.InputPins, phly.PinDescr{Name: save_imginput, Purpose: "Image input."})
	descr.OutputPins = append(descr.OutputPins, phly.PinDescr{Name: save_imgoutput, Purpose: "Image output. All input items are provided, even if the save failed."})
	return descr
}

func (n *save) Instantiate(args phly.InstantiateArgs, cfg interface{}) (phly.Node, error) {
	return &save{}, nil
}

func (n *save) Process(args phly.ProcessArgs, stage phly.NodeStage, input phly.Pins, output phly.NodeOutput) error {
	var err error
	doc := &phly.Doc{MimeType: MimeTypeImagePhly}
	phly.WalkItems(input, save_imginput, func(channel string, src *phly.Doc, index int, _item interface{}) {
		if item, ok := _item.(*PhlyImage); ok {
			err = phly.MergeErrors(err, n.saveImage(args, item, doc))
		}
	})
	if len(doc.Items) > 0 {
		output.SendPins(phly.PinBuilder{}.Add(save_imgoutput, doc).Pins())
	}
	// Run once and we're done
	output.SendMsg(phly.MsgFromStop(nil))
	return err
}

func (n *save) StopNode(args phly.StoppedArgs) error {
	return nil
}

func (n *save) saveImage(args phly.ProcessArgs, img *PhlyImage, doc *phly.Doc) error {
	if img == nil || img.Img == nil {
		return phly.NewBadRequestError("phly/img/save on invalid image")
	}

	dstname, err := n.makeFilename(args, img)
	if err != nil {
		return err
	}
	err = n.makeDir(dstname)
	if err != nil {
		return nil
	}
	f, err := os.Create(dstname)
	if err != nil {
		return err
	}
	defer f.Close()
	doc.AppendItem(img)
	return png.Encode(f, img.Img)
}

// makeFilename() applies my variables to make the filename.
func (n *save) makeFilename(args phly.ProcessArgs, img *PhlyImage) (string, error) {
	// Start with getting the necessary pieces from the source
	src := img.SourceFile
	srcdir := filepath.Dir(src)
	srcext := filepath.Ext(src)
	srcbase := strings.TrimSuffix(filepath.Base(src), srcext)

	return args.Env().ReplaceVars(n.File, "${srcpath}", src, "${srcdir}", srcdir, "${srcbase}", srcbase, "${srcext}", srcext), nil
}

// makeDir() makes the dir if it doesn't exist.
func (n *save) makeDir(filename string) error {
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModeDir)
	}
	return nil
}
