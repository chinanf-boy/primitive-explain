
## model

`primitive/primitive/model.go`

``` go
type Model struct {
	Sw, Sh     int
	Scale      float64
	Background Color
	Target     *image.RGBA
	Current    *image.RGBA
	Context    *gg.Context
	Score      float64
	Shapes     []Shape
	Colors     []Color
	Scores     []float64
	Workers    []*Worker
}

```

### NewModel

``` go
// main.go 
// input: 重调大小的图片, bg: 颜色, OutputSize: 输出图片大小, Workers: 多少个工作人员
	model := primitive.NewModel(input, bg, OutputSize, Workers)

```

``` go
func NewModel(target image.Image, background Color, size, numWorkers int) *Model {
	w := target.Bounds().Size().X
	h := target.Bounds().Size().Y
	aspect := float64(w) / float64(h)
	var sw, sh int
	var scale float64
	if aspect >= 1 { // 宽高比例 大比小
		sw = size
		sh = int(float64(size) / aspect)
		scale = float64(size) / float64(w)
	} else {
		sw = int(float64(size) * aspect)
		sh = size
		scale = float64(size) / float64(h)
	}

	model := &Model{}
	model.Sw = sw // 宽
	model.Sh = sh // 高
	model.Scale = scale // 比例
	model.Background = background // 背景颜色
	model.Target = imageToRGBA(target) // 获得个颜色
	model.Current = uniformRGBA(target.Bounds(), background.NRGBA()) // 混合
	model.Score = differenceFull(model.Target, model.Current) // 差多少
	model.Context = model.newContext() // 拿到 画画 的形状
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(model.Target)
		model.Workers = append(model.Workers, worker)
	}
	return model
}

```

- [x] [imageToRGBA](./util.md#imagetorgba)

> 由图片变个颜色返回
 

- [x] [uniformRGBA](./util.md#uniformrgba)

> 颜色换个图片
 

- [ ] [differenceFull](./core.md#differencefull)

- [x] [model.newContext](#newcontext)

- [ ] [NewWorker](./worker.md#newworker)

---

### newContext

由 `gg` 主导的 图片内容

``` go
func (model *Model) newContext() *gg.Context {
	dc := gg.NewContext(model.Sw, model.Sh)
	dc.Scale(model.Scale, model.Scale)
	dc.Translate(0.5, 0.5)
	dc.SetColor(model.Background.NRGBA())
	dc.Clear()
	return dc
}
```

- [gg](#gg)

> 2d 画画库

---

#### gg

- [x] gg例子 

Go Graphics - 使用简单的API在Go中进行2D渲染。

```
primitive-explain gg
```

``` go
// primitive-explain/examples/gg.go
func UseGG() {
	dc := gg.NewContext(1000, 1000)
	dc.DrawCircle(500, 500, 400)
	dc.SetRGB(0, 0, 0)
	dc.Fill()
	dc.SavePNG("out.png")
}

```

