package phly_img

import (
	"github.com/hackborn/phly"
	"golang.org/x/image/draw"
	"image"
)

const (
	scale_imginput  = "0"
	scale_imgoutput = "0"
)

// scale node scales images.
type scale struct {
	Abs Sizei `json:"abs,omitempty"`
	Rel Sizef `json:"rel,omitempty"`
}

func (n *scale) Run(args phly.RunArgs, input, output phly.Pins) error {
	var err error
	for _, doc := range input.Get(scale_imginput) {
		err = phly.MergeErrors(err, n.runDoc(args, doc, output))
	}
	return err
}

// runDoc() iterates the docs, pages and items, scaling each image.
func (n *scale) runDoc(args phly.RunArgs, srcdoc *phly.Doc, output phly.Pins) error {
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
			err = phly.MergeErrors(err, n.scaleImage(img, dstpage))
		}
		if len(dstpage.Items) > 0 {
			dstdoc.AddPage(dstpage)
		}
	}
	if len(dstdoc.Pages) > 0 {
		output.Add(scale_imgoutput, dstdoc)
	}
	return err
}

func (n *scale) scaleImage(img *PhlyImage, page *phly.Page) error {
	if img == nil || img.Img == nil {
		return phly.BadRequestErr
	}

	srcr := img.Img.Bounds()
	scaler := draw.BiLinear.NewScaler(200, 200, srcr.Size().X, srcr.Size().Y)

	dstr := image.Rect(0, 0, 200, 200)
	dst := image.NewRGBA(dstr)
	ops := draw.Options{}
	scaler.Scale(dst, dstr, img.Img, srcr, draw.Over, &ops)

	page.AddItem(&PhlyImage{Img: dst, SourceFile: img.SourceFile})
	return nil
}

func (n *scale) Instantiate(cfg interface{}) (phly.Node, error) {
	return &scale{}, nil
}
