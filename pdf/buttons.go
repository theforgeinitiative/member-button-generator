package pdf

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"

	"github.com/go-pdf/fpdf"
)

//go:embed button.png
var button []byte

//go:embed CaveatBrush-Regular.json
var fontJSON []byte

//go:embed CaveatBrush-Regular.z
var fontZ []byte

const ButtonsPerRow = 3
const ButtonsPerPage = 12
const RowsPerPage = 4

const ButtonDiameter float64 = 2.75
const LeftMargin float64 = 0.125

const NameOffset float64 = 0.95
const NameHeight float64 = 0.7

const FontColorR = 0
const FontColorG = 92
const FontColorB = 185

func RenderButtons(names []string) ([]byte, error) {
	pdf := fpdf.New("P", "in", "Letter", "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.SetTextColor(FontColorR, FontColorG, FontColorB)
	pdf.AddFontFromBytes("CaveatBrush", "", fontJSON, fontZ)
	pdf.SetFont("CaveatBrush", "", 27)

	// register background image
	buttonReader := bytes.NewReader(button)
	pdf.RegisterImageOptionsReader("bg", fpdf.ImageOptions{ImageType: "png"}, buttonReader)

	var row, col int
	for i, name := range names {
		row = (i / ButtonsPerRow) % RowsPerPage
		col = i % ButtonsPerRow

		// add a page break every 12 buttons
		if i%ButtonsPerPage == 0 {
			pdf.AddPage()
		}

		x := LeftMargin + float64(col)*ButtonDiameter
		y := float64(row) * ButtonDiameter
		pdf.ImageOptions("bg", x, y, ButtonDiameter, ButtonDiameter, false, fpdf.ImageOptions{}, 0, "")

		pdf.SetXY(x, y+NameOffset)
		pdf.CellFormat(ButtonDiameter, NameHeight, name, "", 0, "CM", false, 0, "")
	}

	var b bytes.Buffer
	pdfWriter := bufio.NewWriter(&b)
	err := pdf.Output(pdfWriter)
	if err != nil {
		return nil, fmt.Errorf("failed to generate button pdf: %w", err)
	}
	return b.Bytes(), nil
}
