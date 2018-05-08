package stego

import (
    "fmt"
    "image"
    _ "image/png"
    _ "image/jpeg"
    "image/color"
    "os"
)

const (
    PARTIAL int = 1
    CHAR    int = 2
    ETX     int = 3
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

var char byte = 0
var count uint = 0
/**
 * decode the image pixels into a character
 * if a full byte hasn't been parsed yet, return PARTIAL
 */
func decode(r, g, b, a uint32, c* byte)  int {

    colors := []uint32 {r,g,b}
    var leftover bool = false
    for _, c := range colors {

        if c & 1 > 0 {
            if count <= 8 {
                char += 1<<count
            } else {
                leftover = true
            }
        }
        count++
    }
    if count < 9 {
        return PARTIAL
    }
    *c = char
    char = 0
    if(leftover) {
        char = 1
    }
    leftover = false
    count = 0

    if int(*c) == ETX {
        return ETX
    }
    return CHAR
}

func Read(imagePath string, str *string) int {

    m, err := openImage(imagePath)
    if err != nil {
        fmt.Printf("Error opening image: %s\n", imagePath, err)
        return 1
    }

    // iterate over bitmap, decoding the message until "etx" value is seen
    var msg string
    Outer:
        for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
            for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {

                var c byte = 0
                r, g, b, a := m.At(x,y).RGBA()
                ret := decode(r, g, b, a, &c)

                if ret == ETX {
                    break Outer
                } else if ret == CHAR {
                    msg += string(c)
                }
            }
        }

    fmt.Printf("Read hidden msg (%d chars) from image: \n%s\n", len(msg), msg)

    return 0
}

type Changeable interface {
        Set(x, y int, c color.Color)
}

/**
 * encode a character of the msg to hide into the image
 */
var x int = -1
var y int = -1
var count2 int = 0
func encode(char byte, m image.Image) int {

    if x < 0 {
        x = m.Bounds().Min.X
        y = m.Bounds().Min.Y
    }
    for i := uint32(0); i < 8; i++ {

        r, g, b, a := m.At(x,y).RGBA()
        colors := []uint32 {r,g,b}

        if uint32(char) & 1<<i > 0 {
            colors[count2] = colors[count2]|1
        } else {
            colors[count2] = colors[count2]&0xFFFE
        }

        var rgba color.RGBA
        rgba.R = uint8(colors[0]/a)
        rgba.G = uint8(colors[1]/a)
        rgba.B = uint8(colors[2]/a)
        rgba.A = uint8(a)

        if img, ok := m.(Changeable); ok {
            img.Set(x, y, rgba)
        } else {
            fmt.Println("Unable to modify image")
            return 1
        }

        count2++
        if count2 > 3 {
            // move to the next pixel
            if x+1 > m.Bounds().Max.X {
                y += 1
                x=0
            } else {
                x += 1
            }
            count2 = 0
        }
    }
    return 0
}

func Hide(msg string, imagePath string) int {

    m, err := openImage(imagePath)
    if err != nil {
        fmt.Printf("Error opening image: %s\n", imagePath, err)
        return 1
    }

    maxSize := calcMaxMsgSize(m)
    fmt.Printf("Max length of hidden message for this image: %d bytes\n", maxSize)

    if len(msg) >= maxSize {
        fmt.Printf("Message (%d bytes) can't fit in this image\n", len(msg))
        return 1
    }

    // iterate over message, decode each byte into bits and add into image. 
    // put "etx" ascii value  at end of msg to indicate completion
    for _, s := range(msg) {
        if encode(byte(s), &m) != 0 {
            return 1
        }
    }
    encode(byte(ETX), &m)

    return 0
}

