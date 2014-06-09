# go-otama

Go-bindings for [otama](https://github.com/nagadomi/otama) CBIR Engine Library.

[![Build Status](https://travis-ci.org/hhatto/go-otama.png?branch=master)](https://travis-ci.org/hhatto/go-otama)

## Usage
```go
package main

import (
    "bytes"
    "fmt"
    "github.com/hhatto/go-otama"
    "os"
    "path/filepath"
)

func main() {
    // setup
    ids := make(map[string]string)
    pwd, _ := os.Getwd()
    os.Mkdir("data", 0755)
    o := new(otama.Otama)
    o.Open("test.conf")
    o.CreateDatabase()

    // insert & commit
    buf := bytes.NewBufferString(pwd)
    buf.WriteString("/image")
    filepath.Walk(buf.String(), func(path string, info os.FileInfo, err error) error {
        if info == nil || info.IsDir() {
            return nil
        }

        id, err := o.Insert(path)
        if err != nil {
            fmt.Println("otama Insert() error", err)
            return nil
        }

        ids[id] = path
        return nil
    })
    err := o.Pull()
    if err != nil {
        fmt.Println("otama Pull() error", err)
        os.Exit(1)
    }

    // search
    targetFile := bytes.NewBufferString(pwd)
    targetFile.WriteString("/image/lena.jpg")
    results, err := o.Search(10, targetFile.String())
    if err != nil {
        fmt.Println("otama Search() error", err)
        os.Exit(1)
    }

    // print search result
    for result := range results {
        fmt.Println(fmt.Sprintf("key=%s, sim=%0.3f, file=%s",
            results[result].Id, results[result].Similarity,
            ids[results[result].Id]))
    }
}
```

## License

  * GPLv3

