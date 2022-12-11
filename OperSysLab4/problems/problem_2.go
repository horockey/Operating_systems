package problems

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"
)

type Problem2 struct {
	figuresCount int
	width        int
	height       int
	colors       []color.Color
}

type Problem2Args struct {
	FiguresCount int
	Width        int
	Height       int
}

func (p *Problem2) Init(args interface{}) {
	if argsStruct, ok := args.(Problem2Args); !ok {
		log.Fatal("casting to problem 2 args")
	} else {
		p.figuresCount = argsStruct.FiguresCount
		p.width = argsStruct.Width
		p.height = argsStruct.Height
	}
	p.colors = []color.Color{
		color.RGBA{0, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
		color.RGBA{0, 255, 255, 255},
		color.RGBA{255, 0, 255, 255},
		color.RGBA{255, 255, 0, 255},
	}
}

func (p *Problem2) Run() {
	startChan := make(chan struct{})
	go p.drawCanvas(startChan)
	fmt.Println("Enter any string to trigger...")
	var s string
	fmt.Fscan(os.Stdin, &s)
	startChan <- struct{}{}
	time.Sleep(time.Second)
}

func (p *Problem2) drawCanvas(startChan chan struct{}) {
	<-startChan
	canvas := image.NewRGBA(image.Rect(0, 0, 600, 600))
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < p.figuresCount; i++ {
		x0 := (rand.Int() % canvas.Rect.Dx()) + canvas.Rect.Min.X
		y0 := (rand.Int() % canvas.Rect.Dy()) + canvas.Rect.Min.Y
		w, h := (rand.Int()%100)+1, (rand.Int()%100)+1
		col := p.colors[i%len(p.colors)]
		p.drawRect(canvas, x0, y0, w, h, col)
	}
	out, err := os.Create(
		fmt.Sprintf("D:\\Repos\\Operating_systems\\OperSysLab4\\assets\\2\\result_%s.png", time.Now().Format("2006-01-02_15_04_05")),
	)
	fatalOnErr("creating image file: ", err)
	defer out.Close()
	png.Encode(out, canvas)
}

func (p *Problem2) drawRect(img *image.RGBA, x0, y0, w, h int, col color.Color) {
	if x0 < img.Rect.Min.X || x0 > img.Rect.Max.X || y0 < img.Rect.Min.Y || y0 > img.Rect.Max.Y {
		return
	}
	for x := x0; x < x0+w && x < img.Rect.Max.X; x++ {
		img.Set(x, y0, col)
		img.Set(x, y0+h, col)
	}
	for y := y0; y < y0+h && y < img.Rect.Max.Y; y++ {
		img.Set(x0, y, col)
		img.Set(x0+w, y, col)
	}
}

func fatalOnErr(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
