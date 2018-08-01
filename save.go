package phly_img

import (
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

func (n *save) Run(args phly.RunArgs, input, output phly.Pins) error {
	var err error
	for _, doc := range input.Get(save_imginput) {
		err = phly.MergeErrors(err, n.runDoc(args, doc, output))
	}
	return err
}

// runDoc() iterates the docs, pages and items, translating each filename into an image.
func (n *save) runDoc(args phly.RunArgs, srcdoc *phly.Doc, output phly.Pins) error {
	if srcdoc == nil {
		return phly.MissingDocErr
	}
	var err error
	dstdoc := &phly.Doc{MimeType: MimeTypeImagePhly}
	for _, page := range srcdoc.Pages {
		dstpage := &phly.Page{}
		for _, _img := range page.Items {
			img, ok := _img.(*PhlyImage)
			if !ok {
				return phly.BadRequestErr
			}
			err = phly.MergeErrors(err, n.saveImage(args, img))
			// Always just pass through to output
			dstpage.AddItem(img)
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

func (n *save) saveImage(args phly.RunArgs, img *PhlyImage) error {
	if img == nil || img.Img == nil {
		return phly.BadRequestErr
	}

	dstname, err := n.makeFilename(args, img)
	if err != nil {
		return err
	}
	f, err := os.Create(dstname)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img.Img)
}

func (n *save) makeFilename(args phly.RunArgs, img *PhlyImage) (string, error) {
	// Start with getting the necessary pieces from the source
	src := img.SourceFile
	srcdir := filepath.Dir(src) + string(filepath.Separator)
	srcext := filepath.Ext(src)
	srcbase := strings.TrimSuffix(filepath.Base(src), srcext)

	filename := n.File
	filename = strings.Replace(filename, "${src}", src, -1)
	filename = strings.Replace(filename, "${srcdir}", srcdir, -1)
	filename = strings.Replace(filename, "${srcbase}", srcbase, -1)
	filename = strings.Replace(filename, "${srcext}", srcext, -1)
	return filename, nil
}

func (n *save) Instantiate(cfg interface{}) (phly.Node, error) {
	return &save{}, nil
}
