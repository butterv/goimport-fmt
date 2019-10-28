# goimport-fmt
The import paths are classified into the following three types and sorted in ascending order.
1. Standard package
2. Third-party package
3. Own project package

A blank line is automatically entered between each type.

## Installation

```
$ go get github.com/istsh/goimport-fmt
```

## Usage

```
$ goimport-fmt -filepath path/to/file.go -ownproject github.com/my-repository/my-project
```

## Example
### Before
```go
package main

import (
	"github.com/my-repository/my-project/log"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

func main() {
	str := "12345678"
	
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "parse failed"))
	}
	
	fmt.Printf("str: %s, i: %d\n", str, i)
	log.Infof("str: %s, i: %d\n", str, i)
}
```

### After
```go
package main

import (
	"fmt"
	"strconv"
	
	"github.com/pkg/errors"
	
	"github.com/my-repository/my-project/log"
)

func main() {
	str := "12345678"
	
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "parse failed"))
	}
	
	fmt.Printf("str: %s, i: %d\n", str, i)
	log.Infof("str: %s, i: %d\n", str, i)
}
```

## Author

istsh