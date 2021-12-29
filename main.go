package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/globalsign/publicsuffix"
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
	if err := publicsuffix.Update(); err != nil {
		return err
	}
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
		} else {
			binary.Read(bytes.NewBuffer(record.Name.To4()), binary.BigEndian, &ipv4_int)

			var suffix, _ = publicsuffix.PublicSuffix(record.Value)

			no_tld := strings.TrimRight(record.Value, suffix)
			dots := strings.Split(no_tld, ".")
			fmt.Fprintln(writer, strconv.FormatUint(uint64(ipv4_int), 10)+","+dots[len(dots)-1])
		}
	}
	return nil
}
