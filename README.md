# Unvcode for Golang

好耶，是 [幼女Code](https://github.com/RimoChan/unvcode) 的 Golang 版本！

## 原理和效果

参见原 repo，我才不会再说一遍的（

## 示例

```go
package main

import (
    "fmt"
    "github.com/XiaoMengXinX/go-unvcode"
    "io/ioutil"
)

func main() {
    font, _ := ioutil.ReadFile("NotoSerifSC-SemiBold.otf")
    unv, err := unvcode.New(font)

    /*
        unv.SkipAscii = true
        unv.Mse = 0.1
    */
    
    if err != nil {
        panic(err)
    }
    fmt.Println(unv.Unvcode("Librian幼女娱乐中心开业了，注册即送色图！"))
}
```

新建 Unvcode 对象时需传入字体，example 中的字体文件可以在 [RimoChan/unvcode](https://github.com/RimoChan/unvcode) 找到

可选的参数是 `*unv.SkipAscii = true, *unv.Mse = 0.1`，前者是跳过 ASCII 字符，后者是字符相似度的阈值。（详见源代码）

## 输出

```
Librian幼⼥娱乐㆗⼼开业了，注册即送⾊图！
[-1 -1 -1 -1 -1 -1 -1 -1 0 -1 -1 0.01099964000000089 0 -1 -1 0 -1 -1 -1 -1 -1 0 -1 -1]
```

## 基准

### Python

```python
from unvcode import unvcode
import time

time_start = int(time.time() * 1000)
print("Benckmark starts at:", time_start)

for num in range(0, 23333):
    a, b = unvcode('不许自慰！')

time_done = int(time.time() * 1000)
print("Benckmark done at:", time_done)
print("Duration:", time_done - time_start, "ms")
```

```
> python3 .\test.py
Benckmark starts at: 1644239667488
Benckmark done at: 1644239673265
Duration: 5777 ms
```

### Golang

```go
package main

import (
    "fmt"
    "github.com/XiaoMengXinX/go-unvcode"
    "io/ioutil"
    "time"
)

func main() {
    font, _ := ioutil.ReadFile("NotoSerifSC-SemiBold.otf")
    unv, err := unvcode.New(font)
    if err != nil {
        panic(err)
    }

    timeStart := time.Now().UnixNano() / 1e6
    fmt.Printf("Benchmark starts at: %d\n", timeStart)

    for i := 0; i < 23333; i++ {
        _, _ = unv.Unvcode("不许自慰！")
    }

    timeDone := time.Now().UnixNano() / 1e6
    fmt.Printf("Benchmark done at: %d\n", timeDone)
    fmt.Printf("Duration: %d ms\n", timeDone-timeStart)
}
```

```
> go build -ldflags='-w -s'
> .\test.exe
Benchmark starts at: 1644239353755
Benchmark done at: 1644239391471
Duration: 37716 ms
```

我也不知道怎么优化了，欢迎大佬们的 pr（

## 感谢

[unvcode](https://github.com/RimoChan/unvcode)
