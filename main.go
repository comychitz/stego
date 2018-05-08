package main

import (
    "os"
    "fmt"
    "stego/stego"
    goopt "github.com/droundy/goopt"
)
var usage = func() string {
    return `  options:
    -i image    target image
    -m msg      hide the message \"msg\" into the target image
    -r outfile  read the image, outputting the message into outfile
    -h          print this message`
}

func main() {
    //
    // read command line arguments
    //
    var msg = goopt.String([]string{"-m"}, "", "specify msg to hide")
    var image = goopt.String([]string{"-i"}, "", "target image")
    var outfile = goopt.String([]string{"-r"}, "", "read the image, outputting msg into outfile")
    goopt.Summary = "stego [options]"
    goopt.Help = usage
    goopt.Parse(nil)

    if(len(*image) == 0 ||
      (len(*outfile) == 0 && len(*msg) == 0) ||
      (len(*outfile) > 0 && len(*msg) > 0)) {
          fmt.Println(goopt.Usage())
          os.Exit(1)
    }

    var outmsg string
    var ret int = 0
    if(len(*outfile) > 0) {
        //
        // read msg from image and put into outfile
        //
        ret = stego.Read(*image, &outmsg)
        if ret == 0 {

            // TODO - save msg into outfile

        }
    } else {
        //
        // write secret msg into image
        //
        ret = stego.Hide(*msg, *image)
    }
    os.Exit(ret)
}
