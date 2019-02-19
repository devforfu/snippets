// Examples taken from:
//  https://github.com/adonovan/gopl.io/blob/master/ch8/thumbnail/thumbnail.go
package main

import (
    "fmt"
    "github.com/mitchellh/go-homedir"
    "image"
    "image/jpeg"
    "image/png"
    "io"
    "log"
    "os"
    "path"
    "path/filepath"
    "strings"
)

func main() {
    log.SetFlags(0)
    home, _ := homedir.Dir()
    imagePath := path.Join(home, "Documents", "Resume", "pic.png")
    output, err := ImageFile(imagePath)
    if err != nil { log.Fatalf("%s", err) }
    log.Printf("Thumbnail file: %s", output)
}

const (
    Unknown = iota
    ImageFormatJPEG
    ImageFormatPNG
)

func Format(ext string) int {
    switch ext {
    case ".jpg", ".jpeg":
        return ImageFormatJPEG
    case ".png":
        return ImageFormatPNG
    default:
        return Unknown
    }
}

func ImageFile(infile string) (string, error) {
    ext := filepath.Ext(infile)
    outfile := strings.TrimSuffix(infile, ext) + ".thumb" + ext
    return outfile, ThumbnailFile(outfile, infile, Format(ext))
}

func ThumbnailFile(outfile, infile string, format int) (err error) {
    in, err := os.Open(infile)
    if err != nil { return }
    defer in.Close()

    out, err := os.Create(outfile)
    if err != nil { return }

    if err := ImageStream(out, in, format); err != nil {
        out.Close()
        return fmt.Errorf("scaling %s to %s: %s", infile, outfile, err)
    }

    return out.Close()
}

func ImageStream(w io.Writer, r io.Reader, format int) error {
    src, _, err := image.Decode(r)
    if err != nil { return err }
    dst := ThumbnailImage(src)
    switch format {
    case ImageFormatJPEG:
        return jpeg.Encode(w, dst, nil)
    case ImageFormatPNG:
        return png.Encode(w, dst)
    default:
        return fmt.Errorf("unknown image format")
    }
}

func ThumbnailImage(src image.Image) image.Image {
    xs := src.Bounds().Size().X
    ys := src.Bounds().Size().Y

    width, height := 128, 128
    if aspect := float64(xs)/float64(ys); aspect < 1.0 {
        width = int(128 * aspect)
    } else {
        height = int(128 / aspect)
    }

    xScale := float64(xs)/float64(width)
    yScale := float64(ys)/float64(height)
    dst := image.NewRGBA(image.Rect(0, 0, width, height))

    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            xSrc := int(float64(x)*xScale)
            ySrc := int(float64(y)*yScale)
            dst.Set(x, y, src.At(xSrc, ySrc))
        }
    }

    return dst
}