# mqtt_sync

simple tool that will sync one mqtt server to another. 

## usage

    -d string
        source broker connection string (default "tcp://127.0.0.1:1883")
    -debug
        turn on debug output
    -dp string
        destination broker password
    -du string
        destination broker username
    -p string
        destination topic prefix (e.g. /foo)
    -s string
        source broker connection string (default "tcp://127.0.0.1:1883")
    -sp string
        source broker password
    -su string
        source broker username
    -t string
        source topic (default "#")
        
## dockerize

you can build a mini container (<6mb) of this tool simply by running make. it will build the go app using go:alpine (so you do not need to have go installed) and after that build a docker image that you can run like this:

    docker run --rm --name mqtt_sync nonsenz/mqtt_sync -s tcp://source_server:1883 -d tcp://destination_server:1883 -t /foo/# -p /bar

## todo

- buffer payload if destination is down
- tls
- tests :-P