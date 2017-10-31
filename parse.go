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

func addrFormat(parts []string) string {
    addr := ""
    if strings.Contains(parts[1], ".") {
        // IPv4
        fromTo := strings.Split(parts[1], "-")
        addr = fromTo[0] + ";" + fromTo[1]
    } else {
        // IPv6
        ipv6 := strings.Join(parts[1:], ":")
        addr = ipv6 + ";" + ipv6
    }
    addr += ";"

    return addr
}

// Parse(fh *os.File, ip6 bool)
// Parses the input ripe database for IPv4 ranges and corresponding countries.
// Also includes IPv6 ranges if ip6 param is set to true.
// Status messages are written to stderr, and data to stdout.
func Parse(fh *os.File, ip6 bool) {
    started := time.Now()
    fmt.Fprintf(os.Stderr, "Starting output at: %s\n", started.String())

    scanner := bufio.NewScanner(fh)
    isBlock := false
    str := ""
    spaces, _ := regexp.Compile(" +")

    hasInetnum, _ := regexp.Compile("^inetnum:")
    if ip6 {
        hasInetnum, _ = regexp.Compile("^inet6?num:")
    }

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
                str += addrFormat(parts)
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
    ip6 := flag.Bool("ip6", false, "Include IPv6 ranges as well, defaults to false")

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
    Parse(f, *ip6)
}
