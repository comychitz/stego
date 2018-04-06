# stego
A steganography tool written in Go for hiding a message within an image.

## usage
```
$ stego myImage.jpeg "This my super secret message"
```

## approach
Each image is simply a bitmap of pixels, each pixel having an RGB (red, green,
and blue) color within it (and an alpha, for opaqueness of the pixel), ranging
from 0-255. If we modify the values of the pixels in an elegant way, we can
include any type of information within the image itself. Go has a convenient 
image library for users, which we will take advantage of in this project. 

Hiding our super secret message within the image without producing noise (a
noticeable difference) is essential; if we introduce too much noise, we defeat 
the whole purpose of *hiding*. To keep things simple, we will steal the
least significant bit from each color within each pixel. Thus, the max length 
of our secret message will depend on the size of our image. For example, if we 
had a 500x400 image, the maximum length of our message will be 75,000
characters/bytes (500 * 400 * 3 / 8 = 75,000).

