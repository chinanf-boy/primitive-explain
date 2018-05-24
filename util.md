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

###