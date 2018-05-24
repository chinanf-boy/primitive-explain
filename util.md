## util

`primitive/primitive/util.go`

工具函数🔧

### LoadImage

``` go
// main.go 使用
input, err := primitive.LoadImage(Input) // 加载图片
```

``` go
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path) // 开
	if err != nil {
		return nil, err
	}
	defer file.Close() // 函数结束, 文件关闭
	im, _, err := image.Decode(file) // 转
	return im, err
}
```

### AverageImageColor

返回, 图片的平局颜色

``` go
// main.go 使用
bg = primitive.MakeColor(primitive.AverageImageColor(input))

```

- imageToRGBA

``` go
func AverageImageColor(im image.Image) color.NRGBA {
	rgba := imageToRGBA(im)
	size := rgba.Bounds().Size()
	w, h := size.X, size.Y
	var r, g, b int
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := rgba.RGBAAt(x, y)
			r += int(c.R)
			g += int(c.G)
			b += int(c.B)
		}
	}
	r /= w * h
	g /= w * h
	b /= w * h
	return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
}
```

- color.NRGBA

``` go
// import
	"image/color"
```

### imageToRGBA

由图片变个颜色返回

``` go
func imageToRGBA(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)
	return dst
}
```

- image.NewRGBA

``` go
// import
	"image"
	"image/draw"
	
```

> 试试 `go run main.go img2rgba`

### uniformRGBA

原图片颜色, 再涂上 背景色

``` go
func uniformRGBA(r image.Rectangle, c color.Color) *image.RGBA {
	im := image.NewRGBA(r)
	draw.Draw(im, im.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)
	return im
}
```