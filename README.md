# monitor

Monitor is a command-line utility that can check whether private Asheron's
Call servers are online. A single server can be checked,

  ```
  ./monitor check play.coldeve.online:9000
  ```

Or the entire public list at https://github.com/acresources/serverslist can
be checked,

  ```
  ./monitor list
  ```

## Installation

0. Set up Go
1. Check out the repo and cd into the directory
2. `go build`
3. `./monitor`
