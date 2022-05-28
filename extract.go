package richdiff

import (
	"encoding/binary"
	"errors"
	"image"
	"io/ioutil"
	"math"
	"strconv"
)

var xorKeyOffset int64
var richBuffer []byte
var jsonResults []Result

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
		return nil, errors.New("Rich Header not found")
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
	return jsonResults, nil
}

func RichFileExtraction(filename string) (Results, error) {
	fileBuffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return RichExtraction(fileBuffer)
}

func RichToImg(fileBuffer []byte) (image.Image, error) {
	for i := 1; i <= 200; i++ {
		dwordBuffer := fileBuffer[(0x80 + (i)*4) : 0x80+(i)*4+4]
		if string(dwordBuffer) == "Rich" {
			xorKeyOffset = int64(0x80 + ((i + 1) * 4))
			richSize := 0x80 + (i * 4) - 0x80
			tempBuffer := fileBuffer[0x80 : 0x80+richSize]
			richBuffer = tempBuffer
			break
		}
	}
		byteSize := len(richBuffer)
		sideLength := math.Sqrt(float64(byteSize))
		imageSideLength := int(math.Ceil(sideLength))
		sizeDiff := (imageSideLength * imageSideLength) - byteSize
		if sizeDiff > 0 {
			richBuffer = append(richBuffer, make([]byte, sizeDiff)...)
		}
		img := image.NewRGBA(image.Rect(0, 0, imageSideLength, imageSideLength))
		for i := 0; i < byteSize; i++ {
			img.Pix[i] = richBuffer[i]
		}
	return img, nil
	}
	func RichFileToImg(filename string) (image.Image, error) {
		fileBuffer, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		return RichToImg(fileBuffer)
	}




