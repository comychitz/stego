package main

import (
    "os"
    "fmt"
    "stego/stego"
)

func usage() {
    fmt.Println("usage: stego [options] path/to/image")
    fmt.Println("")
    fmt.Println("options:")
    fmt.Println(" -m msg     hide the message \"msg\" into the image")
    fmt.Println(" -r         read the message inside the image")
    fmt.Println(" -h         prints this help message")
    fmt.Println("")
    os.Exit(1)
}

func main() {
    var msg, image string

    //
    // read command line arguments
    //
    // TODO



    //
    // add message into decoded image
    //
    // TODO

    //
    // encode & save new image
    //
    // TODO


    os.Exit(0)
}
