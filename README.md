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
### Befor
```go
package main

import (
	"github.com/my-repository/my-project/log"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "12345678"
	
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Printf("hash: %s\n", hash)
	log.Infof("hash: %s\n", hash)
}
```

### After
```go
package main

import (
	"fmt"
	
	"golang.org/x/crypto/bcrypt"
	
	"github.com/my-repository/my-project/log"
)

func main() {
	password := "12345678"
	
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Printf("hash: %s\n", hash)
	log.Infof("hash: %s\n", hash)
}
```

## Author

istsh