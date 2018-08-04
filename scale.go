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
	scale_imginput  = "in"
	scale_imgoutput = "out"
)

// scale node scales images.
type scale struct {
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

func (n *scale) Describe() phly.NodeDescr {
	descr := phly.NodeDescr{Id: "phly/img/scale", Name: "Scale Image", Purpose: "Resize images."}
	descr.Cfgs = append(descr.Cfgs, phly.CfgDescr{Name: "width", Purpose: "The width of the final image. Allows variables `${w}` (source width) and `${h}` (source height) and arithmetic expressions (i.e. \"(${w} * 0.5) + 10\")."})
	descr.Cfgs = append(descr.Cfgs, phly.CfgDescr{Name: "height", Purpose: "The height of the final image. Allows variables `${w}` (source width) and `${h}` (source height) and arithmetic expressions (i.e. \"(${w} * 0.5) + 10\")."})
	descr.InputPins = append(descr.InputPins, phly.PinDescr{Name: scale_imginput, Purpose: "Image input."})
	descr.OutputPins = append(descr.OutputPins, phly.PinDescr{Name: scale_imgoutput, Purpose: "The resized images."})
	return descr
}

func (n *scale) Instantiate(args phly.InstantiateArgs, cfg interface{}) (phly.Node, error) {
	return &scale{}, nil
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
			err = phly.MergeErrors(err, n.scaleImage(args, img, dstpage))
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

func (n *scale) scaleImage(args phly.RunArgs, img *PhlyImage, page *phly.Page) error {
	if img == nil || img.Img == nil {
		return phly.BadRequestErr
	}

	srcr := img.Img.Bounds()
	dstsize, err := n.makeSize(args, srcr.Size())
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

func (n *scale) makeSize(args phly.RunArgs, srcsize image.Point) (image.Point, error) {
	x, err := n.makeDimension(srcsize, n.Width, args.ClaValue("width"))
	if err != nil {
		return image.Point{}, err
	}
	y, err := n.makeDimension(srcsize, n.Height, args.ClaValue("height"))
	if err != nil {
		return image.Point{}, err
	}
	return image.Point{x, y}, nil
}

func (n *scale) makeDimension(srcsize image.Point, str, cla string) (int, error) {
	// Make input strings for evaluation
	if cla != "" {
		str = cla
	}
	xstr := strconv.Itoa(srcsize.X)
	ystr := strconv.Itoa(srcsize.Y)
	str = strings.Replace(str, "${w}", xstr, -1)
	str = strings.Replace(str, "${h}", ystr, -1)

	// Evaluate
	fs := token.NewFileSet()
	tv, err := types.Eval(fs, nil, token.NoPos, str)
	if err != nil {
		return 0, err
	}

	// Extract
	v := constant.ToInt(tv.Value)
	if v.Kind() != constant.Int {
		return 0, errors.New("Unparseable scale " + str)
	}
	i, _ := constant.Int64Val(v)
	if i < 1 {
		return 0, errors.New("Unparseable scale " + str)
	}

	return int(i), nil
}
