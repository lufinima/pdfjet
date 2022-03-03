package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"github.com/edragoev1/pdfjet/src"
	"github.com/edragoev1/pdfjet/src/color"
	"github.com/edragoev1/pdfjet/src/compliance"
	"github.com/edragoev1/pdfjet/src/corefont"
	"github.com/edragoev1/pdfjet/src/imagetype"
	"github.com/edragoev1/pdfjet/src/letter"
	"strings"
	"time"
)

// Example15 draws the Canadian flag using a Path object that contains both lines
// and curve segments. Every curve segment must have exactly 2 control points.
func Example15() {
	f, err := os.Create("Example_15.pdf")
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)

	pdf := pdfjet.NewPDF(w, compliance.PDF15)

	font := pdfjet.NewCoreFont(pdf, corefont.Helvetica())

	f, err = os.Open("fonts/Droid/DroidSerif-Regular.ttf.stream")
	reader := bufio.NewReader(f)
	font2 := pdfjet.NewFontStream1(pdf, reader)

	f, err = os.Open("fonts/Droid/DroidSansMono.ttf.stream")
	reader = bufio.NewReader(f)
	font3 := pdfjet.NewFontStream1(pdf, reader)

	f, err = os.Open("images/ee-map.png")
	reader = bufio.NewReader(f)
	image := pdfjet.NewImage(pdf, reader, imagetype.PNG)
	f.Close()

	page := pdfjet.NewPage(pdf, letter.Portrait, true)

	flag := pdfjet.NewBoxAt(85.0, 85.0, 64.0, 32.0)

	path := pdfjet.NewPath()

	path.Add(pdfjet.NewPoint(13.0, 0.0))
	path.Add(pdfjet.NewPoint(15.5, 4.5))

	path.Add(pdfjet.NewPoint(18.0, 3.5))
	path.Add(pdfjet.NewControlPoint(15.5, 13.5))
	path.Add(pdfjet.NewControlPoint(15.5, 13.5))
	path.Add(pdfjet.NewPoint(20.5, 7.5))

	path.Add(pdfjet.NewPoint(21.0, 9.5))
	path.Add(pdfjet.NewPoint(25.0, 9.0))
	path.Add(pdfjet.NewPoint(24.0, 13.0))
	path.Add(pdfjet.NewPoint(25.5, 14.0))
	path.Add(pdfjet.NewPoint(19.0, 19.0))
	path.Add(pdfjet.NewPoint(20.0, 21.5))
	path.Add(pdfjet.NewPoint(13.5, 20.5))
	path.Add(pdfjet.NewPoint(13.5, 27.0))
	path.Add(pdfjet.NewPoint(12.5, 27.0))
	path.Add(pdfjet.NewPoint(12.5, 20.5))
	path.Add(pdfjet.NewPoint(6.0, 21.5))
	path.Add(pdfjet.NewPoint(7.0, 19.0))
	path.Add(pdfjet.NewPoint(0.5, 14.0))
	path.Add(pdfjet.NewPoint(2.0, 13.0))
	path.Add(pdfjet.NewPoint(1.0, 9.0))
	path.Add(pdfjet.NewPoint(5.0, 9.5))

	path.Add(pdfjet.NewPoint(5.5, 7.5))
	path.Add(pdfjet.NewControlPoint(10.5, 13.5))
	path.Add(pdfjet.NewControlPoint(10.5, 13.5))
	path.Add(pdfjet.NewPoint(8.0, 3.5))

	path.Add(pdfjet.NewPoint(10.5, 4.5))
	path.SetClosePath(true)
	path.SetColor(color.Red)
	path.SetFillShape(true)
	path.PlaceIn(flag, 19.0, 3.0)

	path.DrawOn(page)

	box := pdfjet.NewBox()
	box.SetSize(16.0, 32.0)
	box.SetColor(color.Red)
	box.SetFillShape(true)
	box.PlaceIn(flag, 0.0, 0.0)
	box.DrawOn(page)
	box.PlaceIn(flag, 48.0, 0.0)
	box.DrawOn(page)

	path.ScaleBy(15.0)
	path.SetFillShape(false)
	xy := path.DrawOn(page)

	box = pdfjet.NewBox()
	box.SetLocation(xy[0], xy[1])
	box.SetSize(20.0, 20.0)
	box.DrawOn(page)

	font.SetSize(24.0)
	textField := pdfjet.NewTextLine(font, "Hello, World!")
	textField.SetLocation(300.0, 300.0)
	textField.SetColor(color.Blanchedalmond)
	textField.DrawOn(page)

	font2.SetSize(24.0)
	textField2 := pdfjet.NewTextLine(font2, "This is great!")
	textField2.SetLocation(400.0, 400.0)
	textField2.SetColor(color.Blue)
	textField2.SetStrikeout(true)
	textField2.SetUnderline(true)
	textField2.DrawOn(page)

	font2.SetSize(14.0)
	textField2 = pdfjet.NewTextLine(font2, "This is great!")
	textField2.SetLocation(400.0, 500.0)
	textField2.SetColor(color.Blue)
	textField2.DrawOn(page)

	font3.SetSize(24.0)
	textField2 = pdfjet.NewTextLine(font3, "This is great!")
	textField2.SetLocation(400.0, 600.0)
	textField2.SetColor(color.Blue)
	textField2.DrawOn(page)

	image.SetLocation(100.0, 500.0)
	image.ScaleBy(0.5)
	image.DrawOn(page)

	pdf.Complete()

	f.Close()
}

func main() {
	start := time.Now()
	Example15()
	elapsed := time.Since(start).String()
	fmt.Printf("Example_15 => %s\n", elapsed[:strings.Index(elapsed, ".")])
}
