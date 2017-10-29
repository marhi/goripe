/*
MIT License

Copyright (c) 2017 Tomaž Marhat

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
 */

package main

import (
    "log"
    "bufio"
    "fmt"
    "strings"
    "time"
    "regexp"
    "flag"
    "os"
)


// Parse(fh *os.File)
// Parses the input ripe database for IPv4 ranges and corresponding countries.
// Status messages are written to stderr, and data to stdout
func Parse(fh *os.File) {
    started := time.Now()
    fmt.Fprintf(os.Stderr, "Starting output at: %s\n", started.String())

    scanner := bufio.NewScanner(fh)
    isBlock := false
    str := ""
    spaces, _ := regexp.Compile(" +")

    // Add inet6num if you need IPv6 too.
    hasInetnum, _ := regexp.Compile("^inetnum:")
    isInet := false


    for scanner.Scan() {
        text := scanner.Text()

        if hasInetnum.MatchString(text) {
            isBlock = true
            isInet = true
        } else {
            isInet = false
        }

        isCountry := strings.Contains(text, "country:")

        if isBlock && (isInet || isCountry) {
            simple := spaces.ReplaceAllString(text, "")
            parts := strings.Split(simple, ":")

            if isInet {
                fromTo := strings.Split(parts[1], "-")
                str += fromTo[0] + ";" + fromTo[1] + ";"
            } else {
                str += parts[1]
            }
        }

        if text == "" {
            isBlock = false

            if str != "" {
                fmt.Println(str)
                str = ""
            }
        }

    }

    fmt.Fprintf(os.Stderr, "Ending output at: %s, took: %s\n", time.Now().String(), time.Now().Sub(started).String())
}

func main() {
    fName := flag.String("in", "ripe.db", "Input file to parse, defauls to ripe.db")
    flag.Parse()

    if len(*fName) <= 0 {
        flag.Usage()
        os.Exit(6)
    }

    f, err := os.Open(*fName)
    if err != nil {
        log.Fatalln(err)
    }
    defer f.Close()
    Parse(f)
}
