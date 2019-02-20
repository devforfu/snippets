// Examples taken from:
//  https://github.com/adonovan/gopl.io/blob/master/ch8/thumbnail/thumbnail.go
package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/mitchellh/go-homedir"
    "image"
    "image/jpeg"
    "image/png"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path"
    "path/filepath"
    "strings"
)

func init() {
    log.SetFlags(0)
}

func main() {
    conf := parseArgs()
    client := UnsplashClient{conf["access_key"], conf["secret_key"]}
    client.DownloadRandomPhotos(conf["output"], 3)
    //output, err := ImageFile(imagePath)
    //if err != nil { log.Fatalf("%s", err) }
    //log.Printf("Thumbnail file: %s", output)
}

func parseArgs() map[string]string {
    confPath := flag.String(
        "-conf",
        path.Join("config", "unsplash.json"),
        "path to file with Unsplash API config")

    outputPath := flag.String(
        "-out",
        path.Join(MustHomeDir(), "Unsplash"),
        "path to the output folder")

    flag.Parse()
    data, err := ioutil.ReadFile(*confPath)
    if err != nil { log.Fatal(err) }

    config := make(map[string]string)
    json.Unmarshal(data, &config)
    config["output"] = *outputPath
    return config
}

func MustHomeDir() string {
    home, err := homedir.Dir()
    if err != nil { log.Fatal(err) }
    return home
}

const (
    UnsplashBaseURL = "https://api.unsplash.com/"
    UnsplashRandomPhotos = UnsplashBaseURL + "photos/random"
)

type UnsplashClient struct {
    AccessKey, SecretKey string
}

type UnsplashResult struct {
    Id string `json:"id"`
    URLs map[string]string `json:"urls"`
}

func (c *UnsplashClient) DownloadRandomPhotos(dirname string, count int) (err error) {
    log.Printf("Downloading %d image(s) into folder %s\n", count, dirname)
    result, err := c.GetRandomPhotos(count)
    if err != nil { return }

    log.Println("Creating folder")
    err = os.MkdirAll(dirname, os.ModePerm)
    if err != nil { return }

    n := len(result)
    for i, item := range result {
        fmt.Printf("Downloading image %d of %d...\r", i+1, n)
        imageURL := item.URLs["regular"]
        img, err := DownloadImage(imageURL)
        if err != nil { return err }

        filename := path.Join(dirname, fmt.Sprintf("%s.jpeg", item.Id))
        fmt.Printf("Saving downloaded image into file: %s\n", filename)
        file, err := os.Create(filename)
        if err != nil { return err }

        err = jpeg.Encode(file, img, nil)
        if err != nil { return err }

        file.Close()
    }

    return nil
}

func (c *UnsplashClient) GetRandomPhotos(count int) (result []UnsplashResult, err error) {
    client := http.Client{}
    url := fmt.Sprintf("%s?client_id=%s&count=%d", UnsplashRandomPhotos, c.AccessKey, count, )
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", c.SecretKey))

    resp, err := client.Do(req)
    if err != nil { return }
    defer resp.Body.Close()

    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil { return }

    return result, nil
}

func DownloadImage(url string) (image image.Image, err error) {
    resp, err := http.Get(url)
    if err != nil { return }
    defer resp.Body.Close()
    image, err = jpeg.Decode(resp.Body)
    if err != nil { return }
    return image, nil
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