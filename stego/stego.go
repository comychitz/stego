package stego

import (
    "fmt"
    "image"
    _ "image/png"
    _ "image/jpeg"
    "os"
)

/**
 * open the image specified by the user
 */
func openImage(imagePath string) (image.Image, error) {
    var m image.Image
    reader, err := os.Open(imagePath)
    if err != nil {
        fmt.Printf("Error opening file: %s\n", err)
        return m, err
    }
    defer reader.Close()
    m, t, err := image.Decode(reader)
    if err != nil {
        fmt.Printf("Error decoding image: %s\n", err)
        return m, err
    }
    fmt.Printf("%s has dimensions %dx%d and is of type '%s'\n", imagePath,
               m.Bounds().Max.X-m.Bounds().Min.X,
               m.Bounds().Max.Y-m.Bounds().Min.Y, t)
    return m, nil
}

/**
 * calculate the max message size that can be hidden in the picture, returning
 * the value in number of bytes
 */
func calcMaxMsgSize(i image.Image) int  {
    return (i.Bounds().Max.X-i.Bounds().Min.X)*(i.Bounds().Max.Y-i.Bounds().Min.Y)*3/8
}

func Hide(msg string, imagePath string) int {

    m, err := openImage(imagePath)
    if err != nil {
        fmt.Printf("Error opening image: %s\n", imagePath, err)
        return 1
    }

    fmt.Printf("Max length of hidden message for this image: %d bytes\n",
               calcMaxMsgSize(m));

    // TODO
    //
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

