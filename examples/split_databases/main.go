package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/groove-x/rdb"
	"github.com/groove-x/rdb/crc64"
	"github.com/groove-x/rdb/nopdecoder"
)

type decoder struct {
	nopdecoder.NopDecoder
	db    []int
	start []int
	end   []int
}

func (d *decoder) StartDatabase(n int, offset int) {
	d.start = append(d.start, offset)
}

func (d *decoder) EndDatabase(n int, offset int) {
	d.end = append(d.end, offset)
}

func maybeFatal(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	f, err := os.Open(os.Args[1])
	maybeFatal(err)
	d := decoder{start: make([]int, 0), end: make([]int, 0)}
	err = rdb.Decode(f, &d)
	maybeFatal(err)
	fmt.Printf("start: %v\n", d.start)
	fmt.Printf("end: %v\n", d.end)

	f, err = os.Open(os.Args[1])
	maybeFatal(err)

	header := make([]byte, d.start[0])
	_, err = f.Read(header)
	maybeFatal(err)

	for i, s := range d.start {
		e := d.end[i]
		db := make([]byte, e-s)
		_, err = f.Read(db)
		maybeFatal(err)
		out, err := os.Create(os.Args[1] + "." + strconv.Itoa(i))
		maybeFatal(err)

		crc := crc64.New()
		content := append(header, db...)
		content = append(content, 0xff)
		crc.Write(content)
		out.Write(content)
		out.Write(crc.Sum(nil))
	}

}
