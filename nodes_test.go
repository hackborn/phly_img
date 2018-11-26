package phly_img

/*
import (
	"bytes"
	"flag"
	"fmt"
	"github.com/hackborn/phly"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

// --------------------------------
// SCALE

func TestScale1(t *testing.T) {
	n := &scale{Width: "${srcw}*0.5", Height: "${srch}*0.5"}
	runImgTest(t, n, []string{"dog.jpg"}, scale_imginput, []string{"dog_half.golden"})
}

// --------------------------------
// SUPPORT

// runImgTest() takes a node and file input and compares it to file output.
func runImgTest(t *testing.T, n phly.Node, src []string, pinname string, dst []string) {
	input := makePinInput(t, src, pinname)
	output, err := phly.RunNodeTest(t, n, input)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if output == nil {
		fmt.Println("No response")
		t.Fail()
	}
	docs := output.Get(scale_imgoutput)
	if len(docs) != len(dst) {
		fmt.Println("Output mismatch")
		t.Fail()
	}
	if len(docs) == 0 {
		return
	}
	idx := 0
	for _, doc := range docs {
		for _, page := range doc.Pages {
			for _, item := range page.Items {
				img, ok := item.(*PhlyImage)
				if !ok {
					t.Log("Invalid image result")
					t.Fail()
				}
				dstfn := filepath.Join("testdata", dst[idx])
				idx++

				if *update {
					fmt.Println("update golden", dstfn)
					t.Log("update golden file " + dstfn)
					if err = testSaveImageToFile(img, dstfn); err != nil {
						t.Fatalf("failed to update golden file: %s", err)
					}
				}
				// Compare
				prvbytes, err := ioutil.ReadFile(dstfn)
				if err != nil {
					t.Fatalf("failed reading .golden: %s", err)
				}
				newbytes, err := testSaveImageToBuffer(img)
				if err != nil {
					t.Fatalf("failed writing buffer: %s", err)
				}

				if !bytes.Equal(newbytes.Bytes(), prvbytes) {
					t.Errorf("Images do not match")
				}
				idx++
			}
		}
	}
}

func makePinInput(t *testing.T, src []string, pinname string) phly.Pins {
	doc := &phly.Doc{MimeType: MimeTypeImagePhly}
	page := &phly.Page{}
	for _, s := range src {
		fn := filepath.Join("testdata", s)
		err := loadFile(fn, page)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
	}
	doc.AddPage(page)

	pins := phly.NewPins()
	pins.Add(pinname, doc)
	return pins
}

func testSaveImageToBuffer(img *PhlyImage) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := testSaveImageTo(img, buf)
	return buf, err
}

func testSaveImageToFile(img *PhlyImage, fn string) error {
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	return testSaveImageTo(img, f)
}

func testSaveImageTo(img *PhlyImage, w io.Writer) error {
	if img == nil || img.Img == nil || w == nil {
		return phly.BadRequestErr
	}
	return png.Encode(w, img.Img)
}
*/
