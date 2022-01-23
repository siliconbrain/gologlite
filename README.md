# A lightweight logging interface for Go

## Feature highlights
* Minimal, non-opinionated interface
* **B**ring **Y**our **O**wn **B**ackend

## Usage

Add the module as a dependency to your project
```
go get github.com/siliconbrain/gologlite
```

As `gologlite` contains only a minimal API for logging, it's recommended to *create your own custom `log` package* in your project and:
* use alises to re-publish necessary types and functions of `gologlite` and any other modules/packages extending functionality herein
* implement your own helpers and utilities to best adapt to the needs of your project

The rest of your project should import this custom `log` package.

### Example
```go
package log

import (
    "fmt"
    "io"
    "time"

    "github.com/siliconbrain/gologlite/log"
)

type Fields = log.FieldMap
var Event = log.Event

// V is for verbosity
type V int

func (v V) ForEachField(fn func(name string, value interface{}) (stop bool)) {
    fn("verbosity", int(v))
}

func (v V) LookupFieldByName(name string) (value interface{}, found bool) {
    if name == "verbosity" {
        return int(v), true
    }
    return nil, false
}

type WriterTarget struct {
    Out io.Writer
}

func (t WriterTarget) Record(message string, fields log.FieldSet) {
    fmt.Fprintf(t.Out, "%s %s %#v", time.Now().Format(time.RFC3339), message, log.CollapseFieldSets(fields))
}

```

```go
package myproject

import (
    "os"

    "local.dev/myproject/log"
)

func main() {
    logger := log.WriterTarget{Out: os.Stdout}
    // ...
    log.Event(logger, "computation finished", log.V(1), log.Fields{"answer": 42})
}
```
