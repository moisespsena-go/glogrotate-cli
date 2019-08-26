# glogrotation-cli

    ╔═╗┌─┐┬  ╦  ┌─┐┌─┐  ╦═╗┌─┐┌┬┐┌─┐┌┬┐┬┌─┐┌┐┌
    ║ ╦│ ││  ║  │ ││ ┬  ╠╦╝│ │ │ ├─┤ │ ││ ││││
    ╚═╝└─┘o  ╩═╝└─┘└─┘  ╩╚═└─┘ ┴ ┴ ┴ ┴ ┴└─┘┘└┘
                 Go! Log Rotation
    
    Home Page: https://github.com/moisespsena-go/glogrotation-cli
    Author: Moisés P. Sena

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
$ glogrotation -h
```
   
    ╔═╗┌─┐┬  ╦  ┌─┐┌─┐  ╦═╗┌─┐┌┬┐┌─┐┌┬┐┬┌─┐┌┐┌
    ║ ╦│ ││  ║  │ ││ ┬  ╠╦╝│ │ │ ├─┤ │ ││ ││││
    ╚═╝└─┘o  ╩═╝└─┘└─┘  ╩╚═└─┘ ┴ ┴ ┴ ┴ ┴└─┘┘└┘
                 Go! Log Rotation
    
    Home Page: https://github.com/moisespsena-go/glogrotation-cli
    Author: Moisés P. Sena
     
    Starts file writer rotation reads IN and writes to OUT.
    
    EXAMPLES
        NOTE: duration as minutely
    
        A. Basic example
            $ my_program | glogrotation -d m -o program.log
            $ my_program 2>&1 | glogrotation -d m -o program.log
        
        B. Input is STDIN, UDP, TCP and HTTP server
            main terminal:
                $ echo message from stdin | 
                    glogrotation -d m -o program.log -i +udp:localhost:5678+tcp:localhost:5679
    
            secondary terminal:
                a. send message from UDP client:
                    $ echo "message from UDP client" >/dev/udp/localhost/5678
    
                b. send message from TCP client:
                    $ echo "message from TCP client " >/dev/udp/localhost/5679
    
                c. send message from HTTP client:
                    $ curl -X POST -d "message from HTTP client" http://localhost:5680
    
        C. Input is STDIN and UDP server
            main terminal:
                $ (while true; do date; sleep 3; done) | 
                    glogrotation -d m -o program.log -i +udp:localhost:5678
    
            secondary terminal - send message from UDP client:
                $ echo "date from UDP client: "$(date) >/dev/udp/localhost/5678
    
    IN:
        Accept multiple inputs of STDIN, UDP and TCP servers.
        NOTE: Use plus char to join multiple values.
              The first plus char, combines with STDIN.
    
        SERVERS:
            UDP: udp:ADDR, udp4:ADDR, udp6:ADDR ('udp:' is alias of 'udp4:')
                Max message size is 1024 bytes.
    
                Example:
                    udp:localhost:5678
                    udp4:localhost:5678
                    udp:[::1]:5678
                    udp6:[::1]:5678
    
            TCP: tcp:ADDR ('tcp:' is alias of 'tcp4:')
                Example:
                    tcp:localhost:5679
                    tcp4:localhost:5679
                    tcp:[::1]:5679
                    tcp6:[::1]:5679
    
            HTTP: http:ADDR ('http:' is alias of 'http4:')
                - Accept HTTP POST method and copy all request body.
                - Accept WebSocket INPUT on "/" and copy all message body.
    
                Example:
                    http:localhost:5680
                    http4:localhost:5680
                    http:[::1]:5680
                    http6:[::1]:5680
        
        Examples:
            1. Multiple servers
                udp:localhost:5678+tcp:localhost:5679+http:localhost:5680
            2. Multiple servers with STDIN
                +udp:localhost:5678+tcp:localhost:5679+http:localhost:5680
    
    ENV VARIABLES:
        GLOGROTATION_OUT, GLOGROTATION_IN
        GLOGROTATION_HISTORY_DIR, GLOGROTATION_HISTORY_PATH, GLOGROTATION_HISTORY_COUNT 
        GLOGROTATION_DURATION, GLOGROTATION_MAX_SIZE  
        GLOGROTATION_DIR_MODE, GLOGROTATION_FILE_MODE
        GLOGROTATION_SILENT
    
        SET ENV variables to set default flag values.
    
        Usage example:
            Set duration as minutely and enable silent mode:
            $ export GLOGROTATION_DURATION=m
            $ export GLOGROTATION_SILENT=true
            
            run first program as background:
            $ my_first_program | glogrotation -d m -o first_program.log &
    
            run second program:
            $ my_second_program | glogrotation -d m -o second_program.log		
        
    TIME FORMAT:
        %Y - Year. (example: 2006)
        %M - Month with left zero pad. (examples: 01, 12)
        %D - Day with left zero pad. (examples: 01, 31)
        %h - Hour with left zero pad. (examples: 00, 05, 23)
        %m - Minute with left zero pad. (examples: 00, 05, 59)
        %s - Second with left zero pad. (examples: 00, 05, 59)
        %Z - Time Zone. If not set, uses UTC time. (examples: +0700, -0330)
    
    Usage:
      glogrotation [flags]
      glogrotation [command]
    
    Available Commands:
      follower    tail with follower OUT file
      help        Help about any command
      version     Show binary version
    
    Flags:
          --config string         config file
      -M, --dir-mode int          directory perms (default 0750)
      -d, --duration string       rotates every DURATION. Accepted values: Y - yearly, M - monthly, W - weekly, D - daily, h - hourly, m - minutely (default "M")
      -m, --file-mode int         file perms (default 0640)
      -h, --help                  help for glogrotation
      -C, --history-count int     Max history log count
      -c, --history-dir string    history root directory (default "OUT.history")
      -p, --history-path string   dynamic direcotry path inside ROOT DIR using TIME FORMAT (default "%Y/%M")
      -i, --in -                  the INPUT file. - (hyphen char) is STDIN. See INPUT section for details (default "-")
      -S, --max-size string       Forces rotation if current log size is greather then MAX_SIZE. Values in bytes. Examples: 100, 100K, 50M, 1G, 1T (default "50M")
      -o, --out string            the OUTPUT file
          --print                 print current config
          --silent                disable tee to STDOUT
          --udp-max-bs int16      max UDP server buffer size. It's int16 value (default 2048)

## Author
[Moises P. Sena](https://github.com/moisespsena)
