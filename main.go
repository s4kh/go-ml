package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
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

	var w, b float64
	w = 1
	b = 1
	// create linear regression
	l, err := plotter.NewLine(plotter.XYs{
		{X: 0, Y: b}, {X: 14, Y: 14*w + b},
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
