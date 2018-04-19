package stego

import "fmt"

func Hide(msg string, image string) int {

    fmt.Println("Hide() has yet to be implemented")
    return 1

    // TODO
    //
    // 1) open image, get bitmap of pixels
    // 2) iterate over message, decode each byte into bits
    //    and add into bitmap accordingly. put "etx" ascii value 
    //    at end of msg to indicate completion
    // 3) save new image
    //


    return 0
}

func Read(image string, outfile string) int {

    fmt.Println("Read() has yet to be implemented")
    return 1

    // TODO
    //
    // 1) open image, get bitmap of pixels
    // 2) iterate over bitmap, decoding the message until "etx"
    //    value is seen.
    // 3) save msg into outfile
    //


    return 0
}
