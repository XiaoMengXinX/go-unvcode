# 【幼女 Code For Golang】好耶

好耶，是 [【幼女Code】](https://github.com/RimoChan/unvcode) 的 Golang 版本！

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
    if err != nil {
        panic(err)
    }
    fmt.Println(unv.Parse("Librian幼女娱乐中心开业了，注册即送色图！"))
}
```

新建 Unvcode 实例时需传入字体，example 中的字体文件以以在 [RimoChan/unvcode](https://github.com/RimoChan/unvcode) 找到

可选的参数是 `Unvcode.SkipAscii bool = true, Unvcode.Mse float64 = 0.1`，前者是跳过 ASCII 字符，后者是字符相似度的阈值。

## 输出

```
Librian幼⼥娱乐㆗⼼开业了，注册即送⾊图！
[-1 -1 -1 -1 -1 -1 -1 -1 0 -1 -1 0.011000740074007401 0 -1 -1 0 -1 -1 -1 -1 -1 0 -1 -1]
```

## 基准

### Python

```python
from unvcode import unvcode
import time

t1 = time.time()
for num in range(0,23333):
    a, b = unvcode('不许自慰！')
t2 = time.time()
print((t2-t1)*1000)
```

```
> python3 .\test.py
5306.095600128174
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

	for i := 0; i < 23333; i++ {
		_, _ = unv.Parse("不许自慰！")
	}

	fmt.Println(time.Now().UnixNano()/1e6 - timeStart)
}

```

```
> go build -ldflags='-w -s'
> .\test.exe
45312
```

我也不知道怎么优化了，欢迎大佬们的 pr（

## 感谢

[unvcode](https://github.com/RimoChan/unvcode)
