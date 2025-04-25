[![PkgGoDev](https://pkg.go.dev/badge/github.com/s0rg/quadtree)](https://pkg.go.dev/github.com/s0rg/quadtree)
[![License](https://img.shields.io/badge/license-MIT%20License-blue.svg)](https://github.com/s0rg/quadtree/blob/master/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/s0rg/quadtree)](go.mod)
[![Tag](https://img.shields.io/github/v/tag/s0rg/quadtree?sort=semver)](https://github.com/s0rg/quadtree/tags)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

[![CI](https://github.com/s0rg/quadtree/workflows/ci/badge.svg)](https://github.com/s0rg/quadtree/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/s0rg/quadtree)](https://goreportcard.com/report/github.com/s0rg/quadtree)
[![Maintainability](https://qlty.sh/badges/ae1e81be-3dc3-4cf6-ae43-cd4b1344e765/maintainability.svg)](https://qlty.sh/gh/s0rg/projects/quadtree)
[![Code Coverage](https://qlty.sh/badges/ae1e81be-3dc3-4cf6-ae43-cd4b1344e765/test_coverage.svg)](https://qlty.sh/gh/s0rg/projects/quadtree)
![Issues](https://img.shields.io/github/issues/s0rg/quadtree)

# quadtree

[Quadtree](https://en.wikipedia.org/wiki/Quadtree) for golang.

# features

- generic
- heavy optimized
- zero-alloc
- 100% test coverage

# example
```go
import (
    "log"

    "github.com/s0rg/quadtree"
)

func main() {
    // width, height and max depth for new tree
    tree := quadtree.New[int](100.0, 100.0, 4)

    // add some points
    tree.Add(10.0, 10.0, 5.0, 5.0, 1)
    tree.Add(15.0, 20.0, 10.0, 10.0, 2)
    tree.Add(40.0, 10.0, 4.0, 4.0, 3)
    tree.Add(90.0, 90.0, 5.0, 8.0, 4)

    val, ok := tree.Get(9.0, 9.0, 11.0, 11.0)
    if !ok {
        log.Fatal("not found")
    }

    log.Println(val) // should print 1

    const (
        distance = 20.0
        count = 2
    )

    tree.KNearest(80.0, 80.0, distance, count, func(x, y, w, h float64, val int) {
        log.Printf("(%f, %f, %f, %f) = %d", x, y, w, h, val)
    })

    // output: (90.000000, 90.000000, 5.000000, 8.000000) = 4
}
```

# benchmark

```
goos: linux
goarch: amd64
pkg: github.com/s0rg/quadtree
cpu: AMD Ryzen 5 5500U with Radeon Graphics
BenchmarkNode/Insert-12             14974236      71.07 ns/op       249 B/op          0 allocs/op
BenchmarkNode/Del-12                6415672       188.3 ns/op         0 B/op          0 allocs/op
BenchmarkNode/Search-12             21702474      51.83 ns/op         0 B/op          0 allocs/op
BenchmarkTree/Add-12                18840514      67.83 ns/op       241 B/op          0 allocs/op
BenchmarkTree/Get-12                21204722      55.46 ns/op         0 B/op          0 allocs/op
BenchmarkTree/Move-12               8061322       147.5 ns/op         0 B/op          0 allocs/op
BenchmarkTree/ForEach-12            18723290      58.60 ns/op         0 B/op          0 allocs/op
BenchmarkTree/KNearest-12           3595956       324.7 ns/op         0 B/op          0 allocs/op
BenchmarkTree/Del-12                6234123       193.1 ns/op         0 B/op          0 allocs/op
PASS
ok      github.com/s0rg/quadtree    12.666s
```
