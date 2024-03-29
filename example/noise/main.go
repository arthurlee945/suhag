package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/arthurlee945/suhag"
	"github.com/arthurlee945/suhag/example/utility"
	"github.com/arthurlee945/suhag/noise"
	"github.com/arthurlee945/suhag/vec"
	"github.com/fzipp/canvas"
)

func main() {

	utility.StartCanvas(700, 700, "Noise View", runCanvas)
}

func runCanvas(ctx *canvas.Context) {
	ctx.SetFillStyle(color.RGBA{0x08, 0x08, 0x08, 0xff})

	engine := NewNoiseView(ctx.CanvasWidth(), ctx.CanvasHeight())

	for {
		select {
		case event := <-ctx.Events():
			if _, ok := event.(canvas.CloseEvent); ok {
				return
			}
			engine.Handle(event)
		default:
			engine.Draw(ctx)
			ctx.Flush()
			time.Sleep(5 * time.Millisecond)
		}
	}
}

type NoiseView struct {
	size   *vec.Vec2
	offset *vec.Vec2
	noise  *noise.Noise
}

func NewNoiseView(canvasWidth, canvasHeight int) *NoiseView {
	noiseview := &NoiseView{
		size:   &vec.Vec2{float64(canvasWidth), float64(canvasHeight)},
		offset: &vec.Vec2{0, 0},
		noise:  noise.NewNoise(noise.WithSeededPermutation(8, noise.PERMUTATION_SIZE)),
	}
	return noiseview
}

func (nv *NoiseView) Draw(ctx *canvas.Context) {
	nv.draw2D(ctx)
}

func (nv *NoiseView) draw1D(ctx *canvas.Context) {
	xoff := nv.offset[0]
	ctx.ClearRect(0, 0, float64(nv.size[0]), float64(nv.size[1]))
	ctx.BeginPath()
	for x := range int(nv.size[0]) {
		y, err := suhag.Map(nv.noise.Run(xoff, 0, 0), 0, 1, 0, float64(nv.size[1]))
		if err != nil {
			fmt.Println(err.Error())
		}
		ctx.LineTo(float64(x), y)
		ctx.Stroke()
		xoff += 0.01
	}
	nv.offset[0] += 0.01
}

func (nv *NoiseView) draw2D(ctx *canvas.Context) {
	xoff := nv.offset[0]
	ctx.ClearRect(0, 0, float64(nv.size[0]), float64(nv.size[1]))
	rgbaImg := image.NewRGBA(image.Rect(0, 0, int(nv.size[0]), int(nv.size[1])))
	for x := range int(nv.size[0]) {
		yoff := 0.0
		for y := range int(nv.size[1]) {
			brightness, err := suhag.Map(nv.noise.Run(xoff, yoff, nv.offset[1]), 0, 1, 0, 255)
			if err != nil {
				fmt.Println(err.Error())
			}
			rgbaImg.Set(x, y, color.RGBA{75, 0, 130, uint8(brightness)})
			yoff += 0.01
		}
		xoff += 0.01
	}
	ctx.DrawImage(ctx.CreateImageData(rgbaImg), 0, 0)
	nv.offset[1] += 0.01
}

func (nv *NoiseView) Handle(evt canvas.Event) {}

// func (nv *NoiseView) draw2D(ctx *canvas.Context) {
// 	xoff := 0.0
// 	wg := sync.WaitGroup{}
// 	ctx.ClearRect(0, 0, float64(nv.x), float64(nv.y))
// 	rgbaImg := image.NewRGBA(image.Rect(0, 0, nv.x, nv.y))
// 	for x := range nv.x {
// 		yoff := 0.0
// 		for y := range nv.y {
// 			wg.Add(1)
// 			go func() {
// 				brightness, err := suhag.Map(nv.noise.Run(xoff, yoff, 0), 0, 1, 0, 255)
// 				fmt.Println(xoff, yoff)
// 				if err != nil {
// 					fmt.Println(err.Error())
// 				}
// 				rgbaImg.Set(x, y, color.RGBA{255, 0, 0, uint8(brightness)})
// 				wg.Done()
// 			}()
// 			yoff += 0.01
// 		}
// 		xoff += 0.01
// 	}
// 	wg.Wait()
// 	ctx.DrawImage(ctx.CreateImageData(rgbaImg), 0, 0)
// }
