package phly_img

import (
	"errors"
	"github.com/hackborn/phly"
	"go/constant"
	"go/token"
	"go/types"
	"golang.org/x/image/draw"
	"image"
	"strconv"
	"strings"
)

const (
	scale_imginput  = "0"
	scale_imgoutput = "0"
)

// scale node scales images.
type scale struct {
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`

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
	dstsize, err := n.makeSize(srcr.Size())
	if err != nil {
		return err
	}
	scaler := draw.BiLinear.NewScaler(dstsize.X, dstsize.Y, srcr.Size().X, srcr.Size().Y)

	dstr := image.Rect(0, 0, dstsize.X, dstsize.Y)
	dst := image.NewRGBA(dstr)
	ops := draw.Options{}
	scaler.Scale(dst, dstr, img.Img, srcr, draw.Over, &ops)

	page.AddItem(&PhlyImage{Img: dst, SourceFile: img.SourceFile})
	return nil
}

func (n *scale) makeSize(srcsize image.Point) (image.Point, error) {
	// Make input strings for evaluation
	xstr := strconv.Itoa(srcsize.X)
	ystr := strconv.Itoa(srcsize.Y)
	wstr := n.Width
	wstr = strings.Replace(wstr, "${w}", xstr, -1)
	wstr = strings.Replace(wstr, "${h}", ystr, -1)
	hstr := n.Height
	hstr = strings.Replace(hstr, "${w}", xstr, -1)
	hstr = strings.Replace(hstr, "${h}", ystr, -1)

	// Evaluate
	fs := token.NewFileSet()
	wtv, err := types.Eval(fs, nil, token.NoPos, wstr)
	if err != nil {
		return image.Point{}, err
	}
	htv, err := types.Eval(fs, nil, token.NoPos, hstr)
	if err != nil {
		return image.Point{}, err
	}

	// Extract
	wv := constant.ToInt(wtv.Value)
	hv := constant.ToInt(htv.Value)
	if wv.Kind() != constant.Int || hv.Kind() != constant.Int {
		return image.Point{}, errors.New("Unparseable scale " + wstr + " or " + hstr)
	}
	wi, _ := constant.Int64Val(wv)
	hi, _ := constant.Int64Val(hv)
	if wi < 1 || hi < 1 {
		return image.Point{}, errors.New("Unparseable scale " + wstr + " or " + hstr)
	}

	return image.Point{int(wi), int(hi)}, nil
}

func (n *scale) Instantiate(cfg interface{}) (phly.Node, error) {
	return &scale{}, nil
}
