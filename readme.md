# Go Fix Error Hook

A new way of error handling.

Before throwing an error, you can first attempt to fix it,
and let any other part of the program try and do something about it.
This means if an error is caused by something externally,
a module (for example) could automatically fix the error.

## Installation

```shell
go get github.com/tkdeng/go-fix-error-hook
```

## Usage

```go

// in this example, we will use an "out of memory" error
var Error_OOM error = errors.New("out of memory")

// example for: fix.Try
func decompress(str string) (string, error) {
  // do stuff...
  str, err := "some large str...", Error_OOM

  // call before checking for nil error
  fix.Try(&err, func() error {
    //note: this method will first run through the fix.Hook list,
    // and clear the cache (as shown in the `fix.Hook` example).

    // retry stuff...
    str, err = "some large string", nil

    if err != nil {
      return err // fix failed (try next hook if same error)
    }
    return nil // nil = fixed error
  })

  if err != nil {
    return "", err
  }
  return str, nil
}

// example for: fix.Hook
var Cache = map[string]string{}

func init(){
  fix.Hook(Error_OOM, func() bool {
    // this method will be called, when fix.Try is called.
    // note: if a previous hook in the list fixes the problem first,
    // and the error has changed after running the retry callback (in fix.Try),
    // this method will be skipped.

    if _, ok := Cache["old data"]; !ok { // if "old data" not in cache, return false
      return false // fix changed nothing (skip to next hook)
    }

    delete(Cache, "old data") // remove "old data" from cache
    return true // fix did something (run the retry callback in fix.Try)
  })
}

```
