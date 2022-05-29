package richdiff

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/r3labs/diff/v3"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

var xorKeyOffset int64
var richBuffer []byte
var jsonResults []Result

type Result struct {
	CompilerPatchLevel  int `diff:"compilerPatchLevel"`
	ProductID           int `diff:"productID"`
	Count               int `diff:"count"`
	MSInternalName      string `diff:"msInternalName"`
	VisualStudioRelease string `diff:"visualStudioRelease"`
}

type Results struct {
	Results       []Result
	DecryptedRich []byte
	ByteImage     image.Image
}

func (r Results) RichTable() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Compiler Patch Level", "Product ID", "Count", "MS Internal Name", "Visual Studio Release"})
	for _, v := range r.Results {
		t.AppendRow([]interface{}{v.CompilerPatchLevel, v.ProductID, v.Count, v.MSInternalName, v.VisualStudioRelease})
	}
	t.Render()
}

func (r Results) Sort() {
	for i := 0; i < len(r.Results); i++ {
		for j := i + 1; j < len(r.Results); j++ {
			if r.Results[i].ProductID > r.Results[j].ProductID {
				r.Results[i], r.Results[j] = r.Results[j], r.Results[i]
			}
		}
	}
}

func (r Results) String() string {
	jsonBytes, _ := json.Marshal(r)
	return string(jsonBytes)
}

func (r Results) WriteJsonToFile(filename string) error {
	file, _ := os.Create(filename)
	defer file.Close()
	err := json.NewEncoder(file).Encode(r)
	if err != nil {
		return err
	}
	return nil
}

func (r Results) WritePngToFile(filename string) error {
	file, _ := os.Create(filename)
	defer file.Close()
	err := png.Encode(file, r.ByteImage)
	if err != nil {
		return err
	}
	return nil
}

func (r Results) DiffResults(or Results) ([]diff.Change, int, error) {
	// Create a diff between the two results then uses  Sørensen–Dice coefficient to determine if the results are similar
	changelog, err := diff.Diff(r.Results, or.Results)
	if err != nil {
		return nil, -1, err
	}
	if len(r.Results) >= len(or.Results) {
		return changelog, dsc(len(r.Results), len(or.Results), len(r.Results)-len(changelog)), nil
	} else {
		return changelog, dsc(len(r.Results), len(or.Results), len(or.Results)-len(changelog)), nil
	}

}

func RichExtraction(fileBuffer []byte) (Results, error) {
	for i := 1; i <= 200; i++ {
		dwordBuffer := fileBuffer[(0x80+(i)*4) : 0x80+(i)*4+4]
		if string(dwordBuffer) == "Rich" {
			xorKeyOffset = int64(0x80 + ((i + 1) * 4))
			richSize := 0x80 + (i * 4) - 0x80
			tempBuffer := fileBuffer[0x80 : 0x80+richSize]
			richBuffer = tempBuffer
			break
		}
	}
	if len(richBuffer) == 0 {
		return Results{}, errors.New("Rich Header not found")
	}
	xorBuffer := fileBuffer[xorKeyOffset : xorKeyOffset+4]
	xordBytes := make([]byte, len(richBuffer))
	for i, v := range richBuffer {
		xordBytes[i] = v ^ xorBuffer[i%4]
	}
	for i := 17; i < len(xordBytes); i = i + 8 {
		arrayBuffer := xordBytes[i-1 : i+7]
		idDisplay := "unlisted - " + strconv.Itoa(int(binary.LittleEndian.Uint16(arrayBuffer[2:4])))
		if prodList[int(binary.LittleEndian.Uint16(arrayBuffer[2:4]))] != "" {
			idDisplay = prodList[int(binary.LittleEndian.Uint16(arrayBuffer[2:4]))]
		}
		productVersion := binary.LittleEndian.Uint16(arrayBuffer[0:2])
		productID := binary.LittleEndian.Uint16(arrayBuffer[2:4])
		vsProduct, vsVersion := vs_version(int(productID))
		count := binary.LittleEndian.Uint32((arrayBuffer[4:8]))
		jsonResults = append(jsonResults, Result{
			CompilerPatchLevel:  int(productVersion),
			ProductID:           int(productID),
			Count:               int(count),
			MSInternalName:      idDisplay,
			VisualStudioRelease: vsProduct + " " + vsVersion,
		})
	}
	r := Results{Results: jsonResults, DecryptedRich: xordBytes}
	err := richToImg(r)
	if err != nil {
		return r, err
	}

	return r, nil
}

func RichFileExtraction(filename string) (Results, error) {
	fileBuffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return Results{}, err
	}
	return RichExtraction(fileBuffer)
}

func richToImg(r Results) error {
	xordBytes := r.DecryptedRich
	byteSize := len(xordBytes)
	sideLength := math.Sqrt(float64(byteSize))
	imageSideLength := int(math.Ceil(sideLength))
	sizeDiff := (imageSideLength * imageSideLength) - byteSize
	if sizeDiff > 0 {
		xordBytes = append(xordBytes, make([]byte, sizeDiff)...)
	}
	img := image.NewRGBA(image.Rect(0, 0, imageSideLength, imageSideLength))

	for i := 0; i < imageSideLength; i++ {
		for j := 0; j < imageSideLength; j++ {
			img.Set(i, j, color.RGBA{
				uint8(xordBytes[j*imageSideLength+i]),
				uint8(255 - xordBytes[j*imageSideLength+i]),
				uint8(255 - xordBytes[j*imageSideLength+i]),
				uint8(xordBytes[j*imageSideLength+i]),
			})
		}
	}
	r.ByteImage = img
	return nil
}


