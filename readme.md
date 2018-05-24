# primitive [![explain](http://llever.com/explain.svg)](https://github.com/chinanf-boy/Source-Explain)

「 用几何图元重现图像。 」


Explanation

> "version": "1.0.0"

[github source](https://github.com/fogleman/primitive)

[中文](./readme.md) | ~~[english](./readme.en.md)~~

---

使用

```
go get -u github.com/fogleman/primitive
primitive -i input.png -o output.png -n 100
```

| 命令 | 默认 | 详述 |
| --- | --- | --- |
| `i` | n/a | 输入文件|
| `o` | n/a | 输出文件|
| `n` | n/a | 形状数量|
| `m` | 1 | 模式：0 =组合，1 =三角形，2 =矩形，3 =椭圆形，4 =圆形，5 =旋转，6 =贝齐尔，7 =旋转圆形，8 =多边形|
| `rep` | 0 | 每次迭代添加N个额外形状，减少搜索次数（主要用于贝塞尔）|
| `nth` | 1 | 保存每个第N帧（仅当输出路径中有％d时）|
| `r` | 256 | 在处理之前将较大的输入图像调整为此大小|
| `s` | 1024 | 输出图像大小|
| `a` | 128 | 颜色alpha（使用`0`让算法为每个形状选择alpha）|
| `bg` | 平均 | 开始背景颜色（十六进制）|
| `j` | 0 | 并行工作者数量（默认使用全部核心）|
| `v` | off | 详细输出|
| `vv` | off | 非常详细的输出|

<img src="https://www.michaelfogleman.com/static/primitive/examples/monalisa.3.2000.gif" width="49%"/> <img src="https://www.michaelfogleman.com/static/primitive/examples/monalisa-original.png" width="49%"/>

---

> 既然是 命令行 , 我们从 项目的`main.go` 开始

先定个示例

``` bash
primitive -i in.png -o out.png -n 100 -v
```

### 本目录

---

### main.go-import

__1.__ 导入

<details>

``` go
package main

import (
	"flag" // 命令行
	"fmt" // 输出
	"log" // 日志
	"math/rand" // 随机
	"os" // 平台
	"path/filepath" // 文件路径
	"runtime" // 运行-cpu数
	"strconv" // 字符串格式
	"strings"
	"time"

	"github.com/fogleman/primitive/primitive"
	"github.com/nfnt/resize" // 纯golang图像调整大小
)

```

</details>

### main.go-var

__2.__ 定义类型

<details>

``` go
var (
	Input      string
	Outputs    flagArray
	Background string
	Configs    shapeConfigArray
	Alpha      int
	InputSize  int
	OutputSize int
	Mode       int
	Workers    int
	Nth        int
	Repeat     int
	V, VV      bool
)

type flagArray []string

func (i *flagArray) String() string {
	return strings.Join(*i, ", ")
}

func (i *flagArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type shapeConfig struct {
	Count  int
	Mode   int
	Alpha  int
	Repeat int
}

type shapeConfigArray []shapeConfig

func (i *shapeConfigArray) String() string {
	return ""
}

func (i *shapeConfigArray) Set(value string) error {
	n, _ := strconv.ParseInt(value, 0, 0)
	*i = append(*i, shapeConfig{int(n), Mode, Alpha, Repeat})
	return nil
}
```

注意 函数`Set`, 当 命令给出 `-n` 输入 100, 触发 `Set` 函数, 进而

补全 `shapeConfigArray`-结构, 试试

- [x] Set 

```
primitive-explain -n 100
```



</details>

### main.go-init

__3.__ 初始化 命令选项

``` go
func init() {
	flag.StringVar(&Input, "i", "", "input image path")
	flag.Var(&Outputs, "o", "output image path")
	flag.Var(&Configs, "n", "number of primitives")
	flag.StringVar(&Background, "bg", "", "background color (hex)")
	flag.IntVar(&Alpha, "a", 128, "alpha value")
	flag.IntVar(&InputSize, "r", 256, "resize large input images to this size")
	flag.IntVar(&OutputSize, "s", 1024, "output image size")
	flag.IntVar(&Mode, "m", 1, "0=combo 1=triangle 2=rect 3=ellipse 4=circle 5=rotatedrect 6=beziers 7=rotatedellipse 8=polygon")
	flag.IntVar(&Workers, "j", 0, "number of parallel workers (default uses all cores)")
	flag.IntVar(&Nth, "nth", 1, "save every Nth frame (put \"%d\" in path)")
	flag.IntVar(&Repeat, "rep", 0, "add N extra shapes per iteration with reduced search")
	flag.BoolVar(&V, "v", false, "verbose")
	flag.BoolVar(&VV, "vv", false, "very verbose")
}
// 正如在 命令行选项 看到的 (变量 命令选项名 默认 详述)

func errorMessage(message string) bool {   // 错误输出
	fmt.Fprintln(os.Stderr, message)
	return false
}

func check(err error) { // 错误检查
	if err != nil {
		log.Fatal(err)
	}
}

```

> init 总在 main 之前运行

---

### main.go-main

__4.__ 重头戏

```go
func main() {
	// parse and validate arguments
	flag.Parse() // 解析命令行
	ok := true
	if Input == "" {
		ok = errorMessage("ERROR: input argument required")
	}
	if len(Outputs) == 0 {
		ok = errorMessage("ERROR: output argument required")
	}
	if len(Configs) == 0 {
		ok = errorMessage("ERROR: number argument required")
	}
	if len(Configs) == 1 {
		Configs[0].Mode = Mode
		Configs[0].Alpha = Alpha
		Configs[0].Repeat = Repeat
	}
	for _, config := range Configs {
		if config.Count < 1 {
			ok = errorMessage("ERROR: number argument must be > 0")
		}
	}
	if !ok {
		fmt.Println("Usage: primitive [OPTIONS] -i input -o output -n count")
		flag.PrintDefaults()
		os.Exit(1) // 错误退出
	}

	// 设置 日志 等级
	if V {
		primitive.LogLevel = 1
	}
	if VV {
		primitive.LogLevel = 2
	}

	// 设置 随机生成器
	rand.Seed(time.Now().UTC().UnixNano())

	// determine worker count
	if Workers < 1 {
		Workers = runtime.NumCPU() // 拿到 CPU 核心
	}
```

- [x] [primitive 日志等级](./log.md)

> 通过简单的比大小, 确定此输出信息是否被显示

---

#### 4.1 读和重设图片

- [x] [LoadImage](./util.md#loadimage)

``` go
	// read input image
	primitive.Log(1, "reading %s\n", Input)
	input, err := primitive.LoadImage(Input) // 加载图片
	check(err)

	// 如果需要缩小输入图像
	size := uint(InputSize) // 256默认
	if size > 0 {
		input = resize.Thumbnail(size, size, input, resize.Bilinear)
	}

```

- [x] resize

> https://github.com/nfnt/resize


---

#### 4.2 确定背景颜色

- [ ] [MakeColor](./color.md#makecolor)

- [ ] [MakeHexColor](./color.md#makehexcolor)

``` go
	// determine background color
	var bg primitive.Color
	if Background == "" {
		bg = primitive.MakeColor(primitive.AverageImageColor(input))
	} else {
		bg = primitive.MakeHexColor(Background)
    }
    
```

#### 4.3 算法



```go
    // run algorithm
    // 原图 颜色 输出大小 cpu数
	model := primitive.NewModel(input, bg, OutputSize, Workers)
    primitive.Log(1, "%d: t=%.3f, score=%.6f\n", 0, 0.0, model.Score)
    
	start := time.Now() // startting
	frame := 0
	for j, config := range Configs { // 这里一般 length == 1
		primitive.Log(1, "count=%d, mode=%d, alpha=%d, repeat=%d\n",
			config.Count, config.Mode, config.Alpha, config.Repeat)

		for i := 0; i < config.Count; i++ { // 这里才会 循环
			frame++

			// 找到最佳形状并将其添加到模型中
			t := time.Now() // ⏰

			n := model.Step(primitive.ShapeType(config.Mode), config.Alpha, config.Repeat)
			nps := primitive.NumberString(float64(n) / time.Since(t).Seconds())

			elapsed := time.Since(start).Seconds()  // ⏰
			primitive.Log(1, "%d: t=%.3f, score=%.6f, n=%d, n/s=%s\n", frame, elapsed, model.Score, n, nps)

```

- [ ] [Step](./model.md#step)

- [ ] [NumberString](./util.md#numberstring)

---

#### 4.4 每次一个形状写入图片

1. 确定 输出文件的路径

2. 根据不同文件后缀名选择保存方式

``` go	
			// 接着 for i := 0; i < config.Count; i++ {
			// write output image(s)
			for _, output := range Outputs {
				ext := strings.ToLower(filepath.Ext(output))
				percent := strings.Contains(output, "%")
				saveFrames := percent && ext != ".gif"
				saveFrames = saveFrames && frame%Nth == 0
				last := j == len(Configs)-1 && i == config.Count-1
				if saveFrames || last {
					path := output
					if percent {
						path = fmt.Sprintf(output, frame)
					} // 都是为了 做出 路径
					primitive.Log(1, "writing %s\n", path)
					switch ext {
					default:
						check(fmt.Errorf("unrecognized file extension: %s", ext))
					case ".png":
						check(primitive.SavePNG(path, model.Context.Image()))
					case ".jpg", ".jpeg":
						check(primitive.SaveJPG(path, model.Context.Image(), 95))
					case ".svg":
						check(primitive.SaveFile(path, model.SVG()))
					case ".gif":
						frames := model.Frames(0.001)
						check(primitive.SaveGIFImageMagick(path, frames, 50, 250))
					}
				}
			}
		}
	}
}
```

- [ ] [SavePNG](./util.md#savepng)
- [ ] [SaveJPG](./util.md#savejpg)
- [ ] [SaveFile](./util.md#savefile)
- [ ] [SaveGIFImageMagick](./util.md#savegifimagemagick)
