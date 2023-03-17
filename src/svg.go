package pdfjet

/**
 * svg.go
 *
Copyright 2023 Innovatics Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type SVG struct {
}

func NewSVG() *SVG {
	svg := new(SVG)
	return svg
}

func (svg *SVG) GetSVGPaths(filename string) []string {
	var paths = make([]string, 0)
	contents, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var inPath = false
	var buffer = make([]rune, 0)
	for _, ch := range contents {
		if !inPath && strings.HasSuffix(string(buffer), "<path d=") {
			inPath = true
			buffer = nil
		} else if inPath && ch == '"' {
			inPath = false
			paths = append(paths, string(buffer))
			buffer = nil
		} else {
			paths = append(paths, string(buffer))
		}
	}
	return paths
}

func isCommand(ch rune) bool {
	// Please note:
	// Capital letter commands use absolute coordinates
	// Small letter commands use relative coordinates
	if ch == 'M' || ch == 'm' { // moveto
		return true
	} else if ch == 'L' || ch == 'l' { // lineto
		return true
	} else if ch == 'H' || ch == 'h' { // horizontal lineto
		return true
	} else if ch == 'V' || ch == 'v' { // vertical lineto
		return true
	} else if ch == 'Q' || ch == 'q' { // quadratic curveto
		return true
	} else if ch == 'T' || ch == 't' { // smooth quadratic curveto
		return true
	} else if ch == 'C' || ch == 'c' { // cubic curveto
		return true
	} else if ch == 'S' || ch == 's' { // smooth cubic curveto
		return true
	} else if ch == 'A' || ch == 'a' { // elliptical arc
		return true
	} else if ch == 'Z' || ch == 'z' { // close path
		return true
	}
	return false
}

func (svg *SVG) GetSVGPathOps(paths []string) []*PathOp {
	operations := make([]*PathOp, 0)
	var op *PathOp
	for _, path := range paths {
		// Path example:
		// "M22.65 34h3v-8.3H34v-3h-8.35V14h-3v8.7H14v3h8.65ZM24 44z"
		var buf = strings.Builder{}
		var token = false
		for _, ch := range path {
			if isCommand(ch) { // open path
				if token {
					op.args = append(op.args, buf.String())
					buf.Reset()
				}
				token = false
				op = NewPathOp(ch)
				operations = append(operations, op)
			} else if ch == ' ' || ch == ',' {
				if token {
					op.args = append(op.args, buf.String())
					buf.Reset()
				}
				token = false
			} else if ch == '-' {
				if token {
					op.args = append(op.args, buf.String())
					buf.Reset()
				}
				token = true
				buf.WriteRune(ch)
			} else if ch == '.' {
				if strings.Contains(buf.String(), ".") {
					op.args = append(op.args, buf.String())
					buf.Reset()
				}
				token = true
				buf.WriteRune(ch)
			} else {
				token = true
				buf.WriteRune(ch)
			}
		}
	}
	// printOperations(operations)
	return operations
}

func printOperations(operations []*PathOp) {
	for _, oper := range operations {
		print(string(oper.cmd) + " ")
		for _, arg := range oper.args {
			print(arg + " ")
		}
		println()
	}
}

func (svg *SVG) GetPDFPathOps(list []*PathOp) []*PathOp {
	operations := make([]*PathOp, 0)
	var lastOp *PathOp
	var x0 float32 = 0.0 // Start of subpath
	var y0 float32 = 0.0
	for _, op := range list {
		if op.cmd == 'M' || op.cmd == 'm' {
			for i := 0; i <= len(op.args)-2; i += 2 {
				var pathOp *PathOp
				x, err := strconv.ParseFloat(op.args[i], 32)
				if err != nil {
					log.Fatal(err)
				}
				y, err := strconv.ParseFloat(op.args[i+1], 32)
				if err != nil {
					log.Fatal(err)
				}
				if op.cmd == 'm' && lastOp != nil {
					x += float64(lastOp.x)
					y += float64(lastOp.y)
				}
				if i == 0 {
					x0 = float32(x)
					y0 = float32(y)
					pathOp = NewPathOpXY('M', float32(x), float32(y))
				} else {
					pathOp = NewPathOpXY('L', float32(x), float32(y))
				}
				operations = append(operations, pathOp)
				lastOp = pathOp
			}
		} else if op.cmd == 'L' || op.cmd == 'l' {
			for i := 0; i <= len(op.args)-2; i += 2 {
				var pathOp *PathOp
				x, err := strconv.ParseFloat(op.args[i], 32)
				if err != nil {
					log.Fatal(err)
				}
				y, err := strconv.ParseFloat(op.args[i+1], 32)
				if err != nil {
					log.Fatal(err)
				}
				if op.cmd == 'l' && lastOp != nil {
					x += float64(lastOp.x)
					y += float64(lastOp.y)
				}
				pathOp = NewPathOpXY('L', float32(x), float32(y))
				operations = append(operations, pathOp)
				lastOp = pathOp
			}
		} else if op.cmd == 'H' || op.cmd == 'h' {
			for i := 0; i < len(op.args); i++ {
				var pathOp *PathOp
				x, err := strconv.ParseFloat(op.args[i], 32)
				if err != nil {
					log.Fatal(err)
				}
				if op.cmd == 'h' && lastOp != nil {
					x += float64(lastOp.x)
				}
				pathOp = NewPathOpXY('L', float32(x), lastOp.y)
				operations = append(operations, pathOp)
				lastOp = pathOp
			}
		} else if op.cmd == 'V' || op.cmd == 'v' {
			for i := 0; i < len(op.args); i++ {
				var pathOp *PathOp
				y, err := strconv.ParseFloat(op.args[i], 32)
				if err != nil {
					log.Fatal(err)
				}
				if op.cmd == 'v' && lastOp != nil {
					y += float64(lastOp.y)
				}
				pathOp = NewPathOpXY('L', lastOp.x, float32(y))
				operations = append(operations, pathOp)
				lastOp = pathOp
			}
		} else if op.cmd == 'Q' || op.cmd == 'q' {
			for i := 0; i <= len(op.args)-4; i += 4 {
				pathOp := NewPathOp('C')
				x1, err := strconv.ParseFloat(op.args[i], 32)
				if err != nil {
					log.Fatal(err)
				}
				y1, err := strconv.ParseFloat(op.args[i+1], 32)
				if err != nil {
					log.Fatal(err)
				}
				x, err := strconv.ParseFloat(op.args[i+2], 32)
				if err != nil {
					log.Fatal(err)
				}
				y, err := strconv.ParseFloat(op.args[i+3], 32)
				if err != nil {
					log.Fatal(err)
				}
				if op.cmd == 'q' {
					x1 += float64(lastOp.x)
					y1 += float64(lastOp.y)
					x += float64(lastOp.x)
					y += float64(lastOp.y)
				}
				// Save the original control point
				pathOp.x1q = float32(x1)
				pathOp.y1q = float32(y1)
				// Calculate the coordinates of the cubic control points
				x1c := lastOp.x + (2.0/3.0)*(float32(x1)-lastOp.x)
				y1c := lastOp.y + (2.0/3.0)*(float32(y1)-lastOp.y)
				x2c := float32(x) + (2.0/3.0)*(float32(x1)-float32(x))
				y2c := float32(y) + (2.0/3.0)*(float32(y1)-float32(y))
				pathOp.addCubicPoints(x1c, y1c, x2c, y2c, float32(x), float32(y))
				operations = append(operations, pathOp)
				lastOp = pathOp
			}
		} else if op.cmd == 'T' || op.cmd == 't' {
			for i := 0; i <= len(op.args)-2; i += 2 {
				pathOp := NewPathOp('C')
				x1 := lastOp.x
				y1 := lastOp.y
				if lastOp.cmd == 'C' {
					// Find the reflection control point
					x1 = 2*lastOp.x - lastOp.x1q
					y1 = 2*lastOp.y - lastOp.y1q
				}
				x, err := strconv.ParseFloat(op.args[i], 32)
				if err != nil {
					log.Fatal(err)
				}
				y, err := strconv.ParseFloat(op.args[i+1], 32)
				if err != nil {
					log.Fatal(err)
				}
				if op.cmd == 't' {
					x += float64(lastOp.x)
					y += float64(lastOp.y)
				}
				// Calculate the coordinates of the cubic control points
				x1c := lastOp.x + (2.0/3.0)*(x1-lastOp.x)
				y1c := lastOp.y + (2.0/3.0)*(y1-lastOp.y)
				x2c := float32(x) + (2.0/3.0)*(x1-float32(x))
				y2c := float32(y) + (2.0/3.0)*(y1-float32(y))
				pathOp.addCubicPoints(x1c, y1c, x2c, y2c, float32(x), float32(y))
				operations = append(operations, pathOp)
				lastOp = pathOp
			}
		} else if op.cmd == 'C' || op.cmd == 'c' {
			for i := 0; i <= len(op.args)-6; i += 6 {
				pathOp := NewPathOp('C')
				x1, err := strconv.ParseFloat(op.args[i], 32)
				if err != nil {
					log.Fatal(err)
				}
				y1, err := strconv.ParseFloat(op.args[i+1], 32)
				if err != nil {
					log.Fatal(err)
				}
				x2, err := strconv.ParseFloat(op.args[i+2], 32)
				if err != nil {
					log.Fatal(err)
				}
				y2, err := strconv.ParseFloat(op.args[i+3], 32)
				if err != nil {
					log.Fatal(err)
				}
				x, err := strconv.ParseFloat(op.args[i+4], 32)
				if err != nil {
					log.Fatal(err)
				}
				y, err := strconv.ParseFloat(op.args[i+5], 32)
				if err != nil {
					log.Fatal(err)
				}
				if op.cmd == 'c' {
					x1 += float64(lastOp.x)
					y1 += float64(lastOp.y)
					x2 += float64(lastOp.x)
					y2 += float64(lastOp.y)
					x += float64(lastOp.x)
					y += float64(lastOp.y)
				}
				pathOp.addCubicPoints(
					float32(x1), float32(y1),
					float32(x2), float32(y2),
					float32(x), float32(y))
				operations = append(operations, pathOp)
				lastOp = pathOp
			}
		} else if op.cmd == 'S' || op.cmd == 's' {
			// Smooth Cubic Curve
		} else if op.cmd == 'A' || op.cmd == 'a' {
			// Elliptical Arc
		} else if op.cmd == 'Z' || op.cmd == 'z' {
			pathOp := NewPathOp('Z')
			pathOp.x = x0
			pathOp.y = y0
			operations = append(operations, pathOp)
			lastOp = pathOp
		}
	}
	return operations
}
