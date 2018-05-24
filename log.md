## primitive log

`primitive/primitive/log.go`

primitive 内部的日志等级设置

### 定义 日志等级

``` go
package primitive

import "fmt"

var LogLevel int 
// 定义 日志等级

```

> 定义 日志等级 被使用在`main.go`

``` go
// main.go
	if V {
		primitive.LogLevel = 1
	}
	if VV {
		primitive.LogLevel = 2
	}
```

---

### 等级输出

```go
func Log(level int, format string, a ...interface{}) {
	if LogLevel >= level { // 简单的比大小
		fmt.Printf(format, a...)
	}
}

func v(format string, a ...interface{}) {
	Log(1, format, a...)
}

func vv(format string, a ...interface{}) {
	Log(2, "  "+format, a...)
}

func vvv(format string, a ...interface{}) {
	Log(3, "    "+format, a...)
}

```

> 等级输出 被使用在 多处 比如 `main.go`

``` go
// main.go 读图片的信息
	primitive.Log(1, "reading %s\n", Input)

```