package main

import (
	"bufio"
	"color"
	"fmt"
	"io/ioutil"
	"letter"
	"log"
	"os"
	"pdfjet"
	"strconv"
	"strings"
	"time"
)

// Example01 --  TODO:
func Example01(mode string) {
	file, err := os.Create("Example_01.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	pdf := pdfjet.NewPDF(w, 0)

	file1, err := os.Open("fonts/OpenSans/OpenSans-Regular.ttf.stream")
	if err != nil {
		log.Fatal(err)
	}
	defer file1.Close()
	reader := bufio.NewReader(file1)
	font1 := pdfjet.NewFontStream1(pdf, reader)
	font1.SetSize(12.0)

	file2, err := os.Open("fonts/Droid/DroidSansFallback.ttf.stream")
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()
	reader = bufio.NewReader(file2)
	font2 := pdfjet.NewFontStream1(pdf, reader)
	font2.SetSize(12.0)

	page := pdfjet.NewPage(pdf, letter.Portrait, true)

	textLine := pdfjet.NewTextLine(font1, "Happy New Year!")
	textLine.SetLocation(70.0, 70.0)
	textLine.DrawOn(page)

	textLine = pdfjet.NewTextLine(font1, "С Новым Годом!")
	textLine.SetLocation(70.0, 100.0)
	textLine.DrawOn(page)

	textLine = pdfjet.NewTextLine(font1, "Ευτυχισμένο το Νέο Έτος!")
	textLine.SetLocation(70.0, 130.0)
	textLine.DrawOn(page)

	textLine = pdfjet.NewTextLine(font1, "新年快樂！")
	textLine.SetFallbackFont(font2)
	textLine.SetLocation(300.0, 70.0)
	textLine.DrawOn(page)

	textLine = pdfjet.NewTextLine(font1, "新年快乐！")
	textLine.SetFallbackFont(font2)
	textLine.SetLocation(300.0, 100.0)
	textLine.DrawOn(page)

	textLine = pdfjet.NewTextLine(font1, "明けましておめでとうございます！")
	textLine.SetFallbackFont(font2)
	textLine.SetLocation(300.0, 130.0)
	textLine.DrawOn(page)

	textLine = pdfjet.NewTextLine(font1, "새해 복 많이 받으세요!")
	textLine.SetFallbackFont(font2)
	textLine.SetLocation(300.0, 160.0)
	textLine.DrawOn(page)

	page = pdfjet.NewPage(pdf, letter.Portrait, true)

	lcgText, err := ioutil.ReadFile("data/LCG.txt")
	if err != nil {
		log.Fatal(err)
	}

	paragraphs := make([]*pdfjet.Paragraph, 0)
	lines := strings.Split(string(lcgText), "\n")
	textline := pdfjet.NewTextLine(font1, "")
	for i, line := range lines {
		if line == "" {
			continue
		}
		paragraph := pdfjet.NewParagraph()
		textline = pdfjet.NewTextLine(font1, line)
		paragraph.Add(textline)
		if i == 0 {
			textLine := pdfjet.NewTextLine(font1,
				"Hello, World! This is a test to check if this line will be wrapped around properly.")
			textLine.SetColor(color.Blue)
			textLine.SetUnderline(true)
			paragraph.Add(textLine)

			textLine = pdfjet.NewTextLine(font1, "This is a test!")
			textLine.SetColor(color.Oldgloryred)
			textLine.SetUnderline(true)
			paragraph.Add(textLine)
		}
		paragraphs = append(paragraphs, paragraph)
	}

	text := pdfjet.NewText(paragraphs)
	text.SetLocation(50.0, 50.0)
	text.SetWidth(500.0)
	xy := text.DrawOn(page)

	size := text.GetSize()
	box := pdfjet.NewBox()
	box.SetLocation(50.0, 50.0)
	box.SetSize(size[0], size[1])
	box.DrawOn(page)

	box = pdfjet.NewBox()
	box.SetLocation(xy[0], xy[1])
	box.SetSize(20.0, 20.0)
	box.DrawOn(page)

	points := text.GetBeginParagraphPoints()
	for i, xy := range points {
		textLine := pdfjet.NewTextLine(font1, strconv.Itoa(i+1)+".")
		textLine.SetLocation(xy[0]-20, xy[1])
		textLine.DrawOn(page)
	}

	page = pdfjet.NewPage(pdf, letter.Portrait, true)

	cjkText, err := ioutil.ReadFile("data/CJK.txt")
	if err != nil {
		log.Fatal(err)
	}

	paragraphs = make([]*pdfjet.Paragraph, 0)
	lines = strings.Split(string(cjkText), "\n")
	textline = pdfjet.NewTextLine(font1, "")
	for _, line := range lines {
		if line == "" {
			continue
		}
		paragraph := pdfjet.NewParagraph()
		textline = pdfjet.NewTextLine(font1, line)
		textline.SetFallbackFont(font2)
		paragraph.Add(textline)
		paragraphs = append(paragraphs, paragraph)
	}

	text = pdfjet.NewText(paragraphs)
	text.SetLocation(50.0, 50.0)
	text.SetWidth(500.0)
	text.DrawOn(page)

	pdf.Complete()
}

func main() {
	start := time.Now()
	Example01("stream")
	elapsed := time.Since(start).String()
	fmt.Printf("Example_01 => %s\n", elapsed[:strings.Index(elapsed, ".")])
}
