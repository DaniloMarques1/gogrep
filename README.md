Replicates the grep unix command line tool.

to build run
```shell
go build main.go
mv main gogrep
```

to use
```shell
gogrep "whatever you want to search" directory -R
```

or
```shell
gogrep "whatever you want to search" file
```
