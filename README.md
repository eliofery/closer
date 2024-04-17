# closer : 🚪 provides a clean exit mechanism for your Go application.

We shouldn't hard-break connections to various parts of our application. 
When stopping an application, it is important to smoothly close connections to active connections, 
be it a database or something else. 
Because you may encounter troubles in the form of data loss and corruption.

This small package was written to solve this problem.

## Installation

```bash
go get github.com/eliofery/closer
```

## Usage

```go
func main() {
    // Create new instance of Closer.
    clr := closer.New()
    
    // Success close something.
    clr.Add(func() error {
        fmt.Print("Hang on! I'm closing some DBs, wiping some trails..")
        time.Sleep(3 * time.Second)
        fmt.Println("  Done.")
        
        return nil
    })
    
    // Error close something.
    clr.Add(func() error {
        return errors.New("error close something")
    })
    
    // Waiting for the application to exit, for example when pressing Ctrl+C in the terminal.
    clr.Wait()
    
    // After closing, wait for all functions to close.
    if err := clr.Close(); err != nil {
        log.Printf("\n%v", err)
    }
}
```

## Inspiration

I was inspired to write this package by the [xlab](https://github.com/xlab/closer).
