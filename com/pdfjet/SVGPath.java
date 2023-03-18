package com.pdfjet;

import java.util.List;

public class SVGPath {
    String data;                // The SVG path data
    List<PathOp> operations;    // The PDF path operations
    int fill = Color.black;     // The fill color or -1 (don't fill)
    int stroke = -1;            // The stroke color or -1 (don't stroke)
    float strokeWidth;          // The stroke width
}