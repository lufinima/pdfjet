package examples;

import java.io.*;
import java.util.*;

import com.pdfjet.*;


/**
 *  Example_75.java
 *
 *  Draw the Canadian flag using a Path object that contains both lines
 *  and curve segments. Every curve segment must have exactly 2 control points.
 */
public class Example_75 {

    public Example_75() throws Exception {

        PDF pdf = new PDF(
                new BufferedOutputStream(
                        new FileOutputStream("Example_75.pdf")));

        Page page = new Page(pdf, Letter.PORTRAIT);

        Box flag = new Box(0f, page.getHeight() - 32f, 64f, 32f);

        Path path = new Path();

        path.add(new Point(13.0f,  0.0f));
        path.add(new Point(15.5f,  4.5f));

        path.add(new Point(18.0f,  3.5f));
        path.add(new Point(15.5f, 13.5f, Point.CONTROL_POINT));
        path.add(new Point(15.5f, 13.5f, Point.CONTROL_POINT));
        path.add(new Point(20.5f,  7.5f));

        path.add(new Point(21.0f,  9.5f));
        path.add(new Point(25.0f,  9.0f));
        path.add(new Point(24.0f, 13.0f));
        path.add(new Point(25.5f, 14.0f));
        path.add(new Point(19.0f, 19.0f));
        path.add(new Point(20.0f, 21.5f));
        path.add(new Point(13.5f, 20.5f));
        path.add(new Point(13.5f, 27.0f));
        path.add(new Point(12.5f, 27.0f));
        path.add(new Point(12.5f, 20.5f));
        path.add(new Point( 6.0f, 21.5f));
        path.add(new Point( 7.0f, 19.0f));
        path.add(new Point( 0.5f, 14.0f));
        path.add(new Point( 2.0f, 13.0f));
        path.add(new Point( 1.0f,  9.0f));
        path.add(new Point( 5.0f,  9.5f));

        path.add(new Point( 5.5f,  7.5f));
        path.add(new Point(10.5f, 13.5f, Point.CONTROL_POINT));
        path.add(new Point(10.5f, 13.5f, Point.CONTROL_POINT));
        path.add(new Point( 8.0f,  3.5f));

        path.add(new Point(10.5f,  4.5f));
        path.setClosePath(true);
        path.setColor(Color.red);
        path.setFillShape(true);
        path.placeIn(flag, 19.0f, 3.0f);

        path.drawOn(page);

        Box box = new Box();
        box.setSize(16.0f, 32.0f);
        box.setColor(Color.red);
        box.setFillShape(true);
        box.placeIn(flag, 0.0f, 0.0f);
        box.drawOn(page);
        box.placeIn(flag, 48.0f, 0.0f);
        box.drawOn(page);

        pdf.complete();
    }


    public static void main(String[] args) throws Exception {
        long t0 = System.currentTimeMillis();
        new Example_75();
        long t1 = System.currentTimeMillis();
        System.out.println("Example_75 => " + (t1 - t0));
    }

}   // End of Example_75.java
