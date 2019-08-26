# glogrotation-cli
Command Line Go! Log Rotation Tool for https://github.com/moisespsena-go/glogrotation

## INSTALLATION

### Binary Download

See to [release page](https://github.com/moisespsena-go/glogrotation-cli/releases).


### Go! auto build

```bash
go get -u github.com/moisespsena-go/glogrotation-cli/glogrotation
```

Executable installed on $GOPATH/bin/glogrotation

### Build from source

```bash
cd $GOPATH/src/github.com/moisespsena-go/glogrotation-cli/glogrotation
```

#### Using Makefile

requires [goreleaser](https://goreleaser.com/).

```bash
make spt
```

See `./dist` directory to show all executables.

#### Default build

```bash
go build main.go
```

## USAGE

```bash
glogrotion -h
```

## Author
[Moises P. Sena](https://github.com/moisespsena)
