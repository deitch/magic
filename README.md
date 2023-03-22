# magic

A native go implementation of the [magic file format](http://www.darwinsys.com/file/).

It also includes a CLI that functions similarly to the `file` command.

There is no reason to use this repository's CLI over the original `file` command, unless you cannot get the original `file` command to work on your system. Its primary purpose is as a testing tool for the library.

## Usage

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/deitch/magic"
)

func main() {
    f, err := os.Open(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    // GetType takes an io.ReaderAt
    info, err := magic.GetType(f)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(info)
}
```