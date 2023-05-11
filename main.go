package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

var iterations int

func main() {
	flag.IntVar(&iterations, "n", 1000, "number of iterations")
	flag.Parse()
	xys, err := readData("data.txt")
	if err != nil {
		log.Fatalf("could not read data.txt: %v", err)
	}

	err = plotData("out.png", xys)
	if err != nil {
		log.Fatalf("could not plot data: %v", err)
	}
	_ = xys
}

func plotData(path string, xys plotter.XYs) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create out.png: %v", err)
	}

	// create scatter with all data points
	p := plot.New()
	s, err := plotter.NewScatter(xys)
	if err != nil {
		return fmt.Errorf("could create scatter: %v", err)
	}
	s.GlyphStyle.Shape = draw.CrossGlyph{}
	s.Color = color.RGBA{R: 255, A: 255}
	p.Add(s)

	// create linear regression
	// y = w*x+b x - is the feature

	w, b := linearRegression(xys, 0.01)

	l, err := plotter.NewLine(plotter.XYs{
		{X: 0, Y: 0*w + b}, {X: 20, Y: 20*w + b},
	})
	if err != nil {
		return fmt.Errorf("could create new line: %v", err)
	}
	p.Add(l)

	wt, err := p.WriterTo(256, 256, "png")
	if err != nil {
		return fmt.Errorf("could not create writer: %v", err)
	}
	_, err = wt.WriteTo(f)
	if err != nil {
		return fmt.Errorf("could not write to out.png: %v", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close file(%s): %v", path, err)
	}
	return nil
}

func readData(path string) (plotter.XYs, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	var xys plotter.XYs
	s := bufio.NewScanner(f)
	for s.Scan() {
		var x, y float64
		_, err := fmt.Sscanf(s.Text(), "%f,%f", &x, &y)
		if err != nil {
			log.Printf("discarding incorrect data point %q:%v", s.Text(), err)
		}
		xys = append(xys, struct{ X, Y float64 }{x, y})
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("could not scan: %v", err)
	}

	return xys, nil
}

func linearRegression(xys plotter.XYs, alpha float64) (w, b float64) {

	for i := 0; i < iterations; i++ {
		dw, db := computeGradientDescent(xys, w, b)
		w += -dw * alpha
		b += -db * alpha
		// fmt.Printf("grad(%.2f, %.2f) = (%.2f, %.2f)\n", w, b, dw, db)
		fmt.Printf("cost(%.2f, %.2f) = %.2f\n", w, b, computeCost(xys, w, b))
	}

	return w, b
}

func computeCost(xys plotter.XYs, w, b float64) float64 {
	// cost = 1/n * sum((y-(w*x+b))^2)
	s := 0.0
	for _, xy := range xys {
		d := xy.Y - (xy.X*w + b)
		s += d * d
	}

	return s / float64(len(xys))
}

func computeGradientDescent(xys plotter.XYs, w, b float64) (dw, db float64) {
	// cost = 1/n * sum((y-(w*x+b))^2)
	// cost/dw = 2/N * sum(-x * (y-(w*x+c)))
	// cost/dc = 2/N * sum(-(y-(w*x+c)))
	for _, xy := range xys {
		d := xy.Y - (xy.X*w + b)
		dw += -xy.X * d
		db += -d
	}
	n := float64(len(xys))

	return 2 / n * dw, 2 / n * db
}
