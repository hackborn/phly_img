package phly_img

import (
	"fmt"
	"github.com/hackborn/phly"
	"github.com/micro-go/parse"
	"golang.org/x/image/draw"
	"image"
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

func (n *scale) Process(args phly.ProcessArgs, stage phly.NodeStage, input phly.Pins, output phly.NodeOutput) error {
	var err error
	doc := &phly.Doc{MimeType: MimeTypeImagePhly}
	phly.WalkItems(input, scale_imginput, func(channel string, src *phly.Doc, index int, _item interface{}) {
		if item, ok := _item.(*PhlyImage); ok {
			err = phly.MergeErrors(err, n.scaleImage(args, item, doc))
		}
	})
	fmt.Println("scale done, size", len(doc.Items))
	if len(doc.Items) > 0 {
		output.SendPins(phly.PinBuilder{}.Add(scale_imgoutput, doc).Pins())
	}
	// Run once and we're done
	output.SendMsg(phly.MsgFromStop(nil))
	return err
}

func (n *scale) StopNode(args phly.StoppedArgs) error {
	return nil
}

func (n *scale) scaleImage(args phly.ProcessArgs, img *PhlyImage, doc *phly.Doc) error {
	if img == nil || img.Img == nil {
		return phly.NewBadRequestError("phly/img/scale on invalid image")
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

	doc.AppendItem(&PhlyImage{Img: dst, SourceFile: img.SourceFile})
	return nil
}

func (n *scale) makeSize(args phly.ProcessArgs, srcsize image.Point) (image.Point, error) {
	x, err := n.makeDimension(args, srcsize, n.Width, args.ClaValue("width"))
	if err != nil {
		return image.Point{}, err
	}
	y, err := n.makeDimension(args, srcsize, n.Height, args.ClaValue("height"))
	if err != nil {
		return image.Point{}, err
	}
	return image.Point{x, y}, nil
}

func (n *scale) makeDimension(args phly.ProcessArgs, srcsize image.Point, str, cla string) (int, error) {
	// Make input strings for evaluation
	if cla != "" {
		str = cla
	}
	str = args.Env().ReplaceVars(str, "${srcw}", srcsize.X, "${srch}", srcsize.Y)
	return parse.SolveInt(str)
}
