package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/globalsign/publicsuffix"
)

type RDNS struct {
	Timestamp string `json:"timestamp"`
	Name      net.IP `json:"name"`
	Value     string `json:"value"`
	Type      string `json:"type"`
}

func main() {
	if err := publicsuffix.Update(); err != nil {
		panic(err.Error())
	}

	file, err := os.Open(os.Args[1])

	if err != nil {
		fmt.Println(err.Error())
	}

	reader, err := gzip.NewReader(file)

	if err != nil {
		fmt.Println(err.Error())
	}

	scanner := bufio.NewScanner(reader)

	writer := bufio.NewWriterSize(os.Stdout, 4096)

	var ipv4_int uint32
	var record RDNS

	for scanner.Scan() {
		if err := json.Unmarshal([]byte(scanner.Text()), &record); err != nil {
			log.Fatal("Unable to parse: %w", err)
		} else {
			binary.Read(bytes.NewBuffer(record.Name.To4()), binary.BigEndian, &ipv4_int)

			var suffix, _ = publicsuffix.PublicSuffix(record.Value)

			no_tld := strings.TrimRight(record.Value, suffix)
			dots := strings.Split(no_tld, ".")
			fmt.Fprintln(writer, strconv.FormatUint(uint64(ipv4_int), 10)+","+dots[len(dots)-1])
		}
	}
}
