# workman

`go.work` cli manager.

## How to install

- method 1: Use the `go` command `go install github.com/dashengyeah/workman@latest`.

- method 2: Go to the [release](https://github.com/dashengyeah/workman/releases) page.

## How to use

just run `workman` in any module inside a `go.work` defined workspace, and you'll see the command line UI interactives.

> BTW, the command itself already support json output and arg-based operations. Run `workman -h` to see the usage.
```
Usage of workman:
  -l    list go.work info
  -n    no command line ui
  -u string
        update go.work (default "{}")
```
> I'm working on a `vscode-go-work-manager` plugin too. But not sure when to finish, ah...🎃
