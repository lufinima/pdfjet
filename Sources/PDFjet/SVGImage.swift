/**
 *  SVGImage.swift
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
import Foundation

/**
 * Used to embed SVG images in the PDF document.
 */
public class SVGImage {
    var x: Float = 0.0  // location x
    var y: Float = 0.0  // location y
    var w: Float = 0.0  // SVG width
    var h: Float = 0.0  // SVG height
    var pdfPathOps: [PathOp]?

    var color: Int32 = Color.black
    var penWidth: Float = 0.3
    var fillPath: Bool = true

    var uri: String?
    var key: String?
    var language: String?
    var actualText: String = Single.space
    var altDescription: String = Single.space

    /**
     * Used to embed SVG images in the PDF document.
     *
     * @param stream the input stream.
     * @throws Exception  if exception occurred.
     */
    public init(_ stream: InputStream) {
        stream.open()
        defer {
            stream.close()
        }
        var paths = [String]()
        let buffer = UnsafeMutablePointer<UInt8>.allocate(capacity: 1)
        defer {
            buffer.deallocate()
        }
        var scalars = [UnicodeScalar]()
        while stream.hasBytesAvailable {
            let read = stream.read(buffer, maxLength: 1)
            if read > 0 {
                scalars.append(UnicodeScalar(buffer[0]))
            }
        }
        var buf = String()
        var token = false
        var param: String?
        for scalar in scalars {
            if !token && buf.hasSuffix(" width=") {
                token = true
                param = "width"
                buf = ""
            } else if !token && buf.hasSuffix(" height=") {
                token = true
                param = "height"
                buf = ""
            } else if !token && buf.hasSuffix("<path d=") {
                token = true
                param = "path"
                buf = ""
            } else if !token && buf.hasSuffix(" fill=") {
                token = true
                param = "fill"
                buf = ""
            } else if !token && buf.hasSuffix(" stroke=") {
                token = true
                param = "stroke"
                buf = ""
            } else if !token && buf.hasSuffix(" stroke-width=") {
                token = true
                param = "stroke-width"
                buf = ""
            } else if token && scalar == UnicodeScalar("\"") {
                token = false
                if param == "width" {
                    w = Float(buf)!
                } else if param == "height" {
                    h = Float(buf)!
                } else if param == "path" {
                    paths.append(buf)
                } else if param == "fill" {
                    color = mapColorNameToValue(buf)
                }
                buf = ""
            } else {
                buf.append(String(scalar))
            }
        }
        let svgPathOps: [PathOp] = SVG.getSVGPathOps(paths)
        pdfPathOps = SVG.getPDFPathOps(svgPathOps)
    }

    func mapColorNameToValue(_ colorName: String) -> Int32 {
        var color = Color.black
        let mirror = Mirror(reflecting: ColorCSS())
        mirror.children.forEach { child in
            if child.label! == colorName {
                color = child.value as! Int32
            }
        }
        return color
    }

    public func getPDFPathOps() -> [PathOp] {
        return self.pdfPathOps!
    }

    /**
     *  Sets the location of this SVG on the page.
     *
     *  @param x the x coordinate of the top left corner of this box when drawn on the page.
     *  @param y the y coordinate of the top left corner of this box when drawn on the page.
     *  @return this SVG object.
     */
    public func setLocation(_ x: Float, _ y: Float) {
        self.x = x
        self.y = y
    }

    /**
     *  Sets the fill path flag to true or false.
     *
     *  @param fillPath if true fills that SVG path, strokes otherwise.
     */
    public func setFillPath(_ fillPath: Bool) {
        self.fillPath = fillPath
    }

    /**
     *  Sets the size of this box.
     *
     *  @param w the width of this box.
     *  @param h the height of this box.
     */
    public func setSize(_ w: Float, _ h: Float) {
        self.w = w
        self.h = h
    }

    public func setPenWidth(_ w: Float) {
        self.w = w
    }

    public func setHeight(_ h: Float) {
        self.h = h
    }

    public func getPenWidth() -> Float {
        return self.w
    }

    public func getHeight() -> Float {
        return self.h
    }

    public func drawOn(_ page: Page) -> [Float] {
        page.addBMC(StructElem.P, language, actualText, altDescription)
        page.setPenWidth(penWidth)
        if fillPath {
            page.setBrushColor(color)
        }
        else {
            page.setPenColor(color)
        }
        for op in pdfPathOps! {
            if op.cmd == "M" {
                page.moveTo(op.x + x, op.y + y)
            } else if op.cmd == "L" {
                page.lineTo(op.x + x, op.y + y)
            } else if op.cmd == "C" {
                page.curveTo(
                    op.x1 + x, op.y1 + y,
                    op.x2 + x, op.y2 + y,
                    op.x + x, op.y + y)
            } else if op.cmd == "Z" {
                if !fillPath {
                    page.closePath()
                }
            }
        }
        if fillPath {
            page.fillPath()
        }
        page.addEMC()

        if (uri != nil || key != nil) {
            page.addAnnotation(Annotation(
                    uri,
                    key,    // The destination name
                    x,
                    y,
                    x + w,
                    y + h,
                    language,
                    actualText,
                    altDescription))
        }
        return [self.x + self.w, self.y + self.h]
    }
}   // End of SVGImage.java
