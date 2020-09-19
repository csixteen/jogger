# jogger

Simple library that allows you to execute a command from your code, given a string representing its name and optional parameters. It writes to `stdout` and `stderr` the progress of the execution (in real-time). You may optionally suppress the output or change the default timeout (currently 30 seconds) after which the process is killed.

# Limitations

I haven't tested this on Windows, but it probably won't work (read the *Overview* [here](https://golang.org/pkg/os/exec/) for more details).

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

# Options

There are currently two options that you can pass to `jogger.Run`:
- `jogger.NoOutput()` - in case you want to suppress the output
- `jogger.Timeout(d time.Duration)` - in case you want to change the default timout after which the process gets killed.

# Contributing

If you find any issues or want to add any functionality, feel free to submit a Pull-request.

# LICENSE
[MIT](https://github.com/csixteen/jogger/blob/master/LICENSE), as always.
