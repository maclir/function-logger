# maclir/function-logger

A wrapper to contain bioler plate code for logging on google cloud functions.

## Installations

```sh
go get -u github.com/maclir/function-logger
```

## Usage

```go
package something

import (
...
	"cloud.google.com/go/logging"
	flogger "github.com/maclir/function-logger"
)

// reused between function invocations.
var logger *logging.Logger

func init() {
	var err error
	logger, err = flogger.New()
	if err != nil {
		panic(err)
	}
}

func Something(w http.ResponseWriter, r *http.Request) {
	defer logger.Flush()

	...

	logger.Log(logging.Info, "This is a log")

	...
}
```
