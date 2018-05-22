package stego

import (
    "fmt"
    "image"
    "image/draw"
    "image/jpeg"
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

func readPixel(x int, y int, *m image.Image) (r, g, b, a uint8) {

    r, g, b, a := outImg.At(x,y).RGBA()
    r = uint8(r)
    g = uint8(g)
    b = uint8(b)
    a = uint8(a)
    fmt.Printf("Pixel [%d,%d] is: R:%d G:%d B:%d \n", x, y, r, g, b, a)
    return r, g, b, a
}

/**
 * return a decoder for decoding the image pixel into a character
 * if a full byte hasn't been parsed yet, return PARTIAL
 */
 func decoder() func(r, g, b uint8, c*byte) int {
     var char byte = 0
     var count uint = 0
     return func(r, g, b uint8, c* byte)  int {
         colors := []uint8 {r, g, b}
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
         fmt.Println("byte: ", *c)
         return CHAR
     }
 }

/**
 * return function for encoding a character of the msg to hide into the image
 */
func encoder(m image.Image) func(byte, *image.RGBA) int {
    x := m.Bounds().Min.X
    y := m.Bounds().Min.Y
    index := uint32(0)
    return func(char byte, m *image.RGBA) int {
        fmt.Printf("Hiding char: %c %d %b\n", char, char, char)
        for i := uint32(0); i < 8; i+=3 {

            r, g, b, a = readPixel(x, y, m)
            colors := []uint8 {r, g, b}

            for j := index; j < 3 && i+j < 8; j++ {
                if uint32(char)&(1<<(7-(i+j))) > 0 {
                    colors[j] = colors[j]|1
                } else {
                    colors[j] = colors[j]&0xFE
                }
            }
            var rgba color.RGBA
            rgba.R = colors[0]
            rgba.G = colors[1]
            rgba.B = colors[2]
            rgba.A = a
            m.SetRGBA(x, y, rgba)
            fmt.Printf("Set pixel [%d,%d] to: R:%d G:%d B:%d A:%d\n\n",
                       x, y, rgba.R, rgba.G, rgba.B, rgba.A)
            if i+j == 8 {
                index = uint32(0)
                if x+1 > m.Bounds().Max.X {
                    y += 1
                    x = m.Bounds().Min.X
                } else {
                    x += 1
                }
            }
        }
        return 0
    }
}

func Hide(msg string, imagePath string) int {

    m, err := openImage(imagePath)
    if err != nil {
        fmt.Printf("Error opening image: %s (%d)\n", imagePath, err)
        return 1
    }
    if len(msg) >= calcMaxMsgSize(m) {
        fmt.Printf("Message (%d bytes) can't fit in this image\n", len(msg))
        return 1
    }
    fmt.Printf("Hidding message of length %d (%d bits): %s\n", len(msg),
               len(msg)*8, msg)

    outImg := image.NewRGBA(m.Bounds())
    draw.Draw(outImg, m.Bounds(), m, image.Point{}, draw.Over)

    encode := encoder(m)
    for _, s := range(msg) {
        if encode(byte(s), outImg) != 0 {
            return 1
        }
    }
    encode(byte(ETX), outImg)

    //
    // TODO parse image name and add -2 to end of it
    //
    var savePath string = imagePath + "-2"
    w, err := os.Create(savePath)
    if err != nil {
        fmt.Printf("Error opening for saving (%s)\n", err)
        return 1
    }
    defer w.Close()
    err =  jpeg.Encode(w, outImg, nil)
    if err != nil {
        fmt.Printf("Error saving image (%s)\n", err)
        return 1
    }
    fmt.Printf("Saved image %s\n", savePath)
    return 0
}

func Read(imagePath string, str *string) int {
    m, err := openImage(imagePath)
    if err != nil {
        fmt.Printf("Error opening image: %s\n", imagePath, err)
        return 1
    }
    img := image.NewRGBA(m.Bounds())
    draw.Draw(img, m.Bounds(), m, image.Point{}, draw.Over)

    decode := decoder()
    Outer:
        for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
            for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {

                var c byte = 0
                r, g, b, _ := readPixel(x, y, img)
                ret := decode(r, g, b, &c)

                if ret == ETX {
                    fmt.Println("got end of text: ", c)
                    break Outer
                } else if ret == CHAR {
                    fmt.Println("read byte: ", c)
                    *str += string(c)
                }
            }
        }
    fmt.Printf("Read hidden msg (%d chars) from image: \n%s\n", len(*str), *str)
    return 0
}

