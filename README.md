[![PkgGoDev](https://pkg.go.dev/badge/github.com/s0rg/quadtree)](https://pkg.go.dev/github.com/s0rg/quadtree)
[![License](https://img.shields.io/badge/license-MIT%20License-blue.svg)](https://github.com/s0rg/quadtree/blob/master/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/s0rg/quadtree)](go.mod)
[![Tag](https://img.shields.io/github/v/tag/s0rg/quadtree?sort=semver)](https://github.com/s0rg/quadtree/tags)

[![Go Report Card](https://goreportcard.com/badge/github.com/s0rg/quadtree)](https://goreportcard.com/report/github.com/s0rg/quadtree)
<!--
[![Maintainability](https://api.codeclimate.com/v1/badges/d6759af1231bf4f60f70/maintainability)](https://codeclimate.com/github/s0rg/quadtree/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/d6759af1231bf4f60f70/test_coverage)](https://codeclimate.com/github/s0rg/quadtree/test_coverage)
-->

# quadtree

Generic quadtree for golang.

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
	tree := quadtree.New[int](100, 100, 4)

    // add some points
	tree.Add(10, 10, 5, 5, 1)
	tree.Add(15, 20, 10, 10, 2)
	tree.Add(40, 10, 4, 4, 3)
	tree.Add(90, 90, 5, 8, 4)

    val, ok := tree.Get(9, 9, 11, 11)
    if !ok {
        log.Fatal("not found")
    }

    log.Println(val) // should print 1

    const (
        searchX = 60.0
        searchY = 60.0
        distance = 35.0
        count = 2
    )

    tree.KNearest(searchX, searchY, distance, count, func(x, y, w, h float64, val int) {
		log.Printf("(%f, %f, %f, %f) = %d", x, y, w, h, val)
	})
}
```

# benchmark
```
pkg: github.com/s0rg/quadtree
cpu: AMD Ryzen 5 5500U with Radeon Graphics
BenchmarkNodeInsert1-12        	4119890	    302.1 ns/op	    139 B/op	      0 allocs/op
BenchmarkNodeInsert10-12       	 314479	     3384 ns/op	   1237 B/op	      0 allocs/op
BenchmarkNodeInsert100-12      	  33385	    33683 ns/op	  14128 B/op	      0 allocs/op
BenchmarkNodeDel10-12          	  38428	    49986 ns/op	      0 B/op	      0 allocs/op
BenchmarkNodeDel100-12         	  10000	   647607 ns/op	      0 B/op	      0 allocs/op
BenchmarkNodeSearch10-12       	 375060	     3071 ns/op	      0 B/op	      0 allocs/op
BenchmarkNodeSearch100-12      	  37975	    31457 ns/op	      0 B/op	      0 allocs/op
BenchmarkTreeKNearest10-12     	 639724	     1842 ns/op	      0 B/op	      0 allocs/op
BenchmarkTreeKNearest100-12    	  64370	    18922 ns/op	      0 B/op	      0 allocs/op
```
