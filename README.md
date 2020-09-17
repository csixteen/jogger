# jogger

Simple library that allows you to execute a command from your code, given a string representing its name and optional parameters. It writes to `stdout` and `stderr` the progress of the execution (in real-time). You may optionally suppress the output or change the default timeout (currently 30 seconds) after which the process is killed.

# Examples

```
package main

import "github.com/csixteen/jogger"

func main() {
        _, _, err := jogger.Run("ls", []string{"-la"})
        if err != nil {
                panic(err)
        }

        // suppressing output

        _, _, err := jogger.Run("ls", []string{"-la"}, jogger.NoOutput())
        if err != nil {
                panic(err)
        }
}
```

