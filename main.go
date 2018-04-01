package main

import (
    "os"
    "fmt"
    "stego/stego"
)

func usage() {
    fmt.Println("Usage: stego path/to/image \"your secret message\"")
    os.Exit(1)
}

func main() {
    var msg, image string

    if(len(os.Args) != 3) {
        usage()
    }
    image = os.Args[1]
    msg = os.Args[2]

    fmt.Printf("Hiding msg: '%s' into '%s'...\n", msg, image)

    var ret int = -1
    ret = stego.Hide(msg, image)

    if(ret != 0) {
        fmt.Println("Problems hiding message!")
        os.Exit(1)
    } else {
        fmt.Println("Done!")
        os.Exit(0)
    }
}
