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

func main() {

    fName := flag.String("in", "", "Input file to parse")
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

    started := time.Now()
    fmt.Fprintf(os.Stderr, "Starting output at: %s\n", started.String())

    scanner := bufio.NewScanner(f)
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

