package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

type RDNS struct {
	Timestamp string `json:"timestamp"`
	Name      net.IP `json:"name"`
	Value     string `json:"value"`
	Type      string `json:"type"`
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	t := time.Now()
	defer func() {
		fmt.Printf("non-suffix time: %v\n", time.Since(t))
	}()

	file, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)

	// writer := bufio.NewWriterSize(os.Stdout, 4096)
	writer := bufio.NewWriter(io.Discard)

	var ipv4_int uint32
	var record RDNS

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return err
		}

		ipv4 := record.Name.To4()
		ipv4_int = (uint32(ipv4[0]) << 24) + (uint32(ipv4[1]) << 16) + (uint32(ipv4[2]) << 8) + (uint32(ipv4[3]))
		// ipv4_int = (uint32(record.Name[12+0]) << 24) + (uint32(record.Name[12+1]) << 16) + (uint32(record.Name[12+2]) << 8) + (uint32(record.Name[12+3]))

		var suffix, _ = publicsuffix.PublicSuffix(record.Value)

		no_tld := strings.TrimRight(record.Value, "."+suffix)
		li := strings.LastIndex(no_tld, ".") // there's strings.Cut in go 1.18 (to be released in Feb 2022)
		domain := no_tld[li+1:]
		fmt.Fprintf(writer, "%v,%v\n", ipv4_int, domain)
	}
	return nil
}
