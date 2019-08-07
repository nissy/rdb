// This is a very basic example of a program that implements rdb.decoder and
// outputs a human readable diffable dump of the rdb file.
package main

import (
	"fmt"
	"os"

	"github.com/tao12345666333/rdb"
	"github.com/tao12345666333/rdb/nopdecoder"
)

type decoder struct {
	db      int
	dbcount int
	count   int
	nopdecoder.NopDecoder
}

func (p *decoder) StartDatabase(n int, offset int) {
	p.db = n
	p.dbcount = 0
}

func (p *decoder) EndDatabase(n int, offset int) {
	fmt.Printf("db=%d count=%d\n", p.db, p.dbcount)
	p.count += p.dbcount
}

func (p *decoder) Set(key, value []byte, expiry int64) {
	p.dbcount += 1
}

func (p *decoder) StartHash(key []byte, length, expiry int64) {
	p.dbcount += 1
}

func (p *decoder) StartSet(key []byte, cardinality, expiry int64) {
	p.dbcount += 1
}

func (p *decoder) StartList(key []byte, length, expiry int64) {
	p.dbcount += 1
}

func (p *decoder) StartZSet(key []byte, cardinality, expiry int64) {
	p.dbcount += 1
}

func (p *decoder) Aux(auxkey, auxval []byte) {
	fmt.Printf("db=%d %q -> %q\n", p.db, auxkey, auxval)
}

func (p *decoder) EndRDB(offset int) {
	fmt.Printf("\nthis rdb keys count=%d\n", p.count)
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
	err = rdb.Decode(f, &decoder{})
	maybeFatal(err)
}
