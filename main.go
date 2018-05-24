package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"strconv"

	"github.com/chinanf-boy/primitive-explain/examples"
)

var (
	Configs shapeConfigArray
	Alpha   int
	Mode    int
	Repeat  int
)

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
	fmt.Println("flag -n ", value)
	fmt.Println("Set running")
	n, _ := strconv.ParseInt(value, 0, 0)
	*i = append(*i, shapeConfig{int(n), Mode, Alpha, Repeat})
	return nil
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	im, _, err := image.Decode(file)
	return im, err
}

func init() {
	flag.Var(&Configs, "n", "number of primitives")
	Mode = 1
	Alpha = 1
	Repeat = 1
}

func main() {
	flag.Parse()
	exampleIndex := flag.Args()

	if exampleIndex[0] == "gg" {
		examples.UseGG()
		fmt.Println("examples gg create new png {out.png}")
		os.Exit(0)
	}
	if exampleIndex[0] == "img2rgba" {

		input, _ := loadImage("lenna.png")

		rgbas := examples.ImageToRGBA(input)

		fmt.Println("so many RGBA", rgbas)
		os.Exit(0)
	}
	for _, config := range Configs {
		fmt.Println("new config", config)
	}
}
