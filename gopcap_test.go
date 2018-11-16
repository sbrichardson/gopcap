package gopcap

import (
	"fmt"
	"io"
	"os"
	"testing"

	gopcap "github.com/0intro/pcap"
	"github.com/google/gopacket/pcap"
	mpcap "github.com/miekg/pcap"
	npcap "go.universe.tf/netboot/pcap"
)

const file = "maccdc2012_00000.pcap"

func BenchmarkReadPcap(b *testing.B) {

	r, err := Open(file)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	for n := 0; n < b.N; n++ {
		_, _, err := r.ReadNextPacket()
		if err == io.EOF {
			fmt.Println("EOF")
			break
		} else if err != nil {
			panic(err)
		}
	}
}

func BenchmarkReadPcap0Intro(b *testing.B) {

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r, err := gopcap.NewReader(f)
	if err != nil {
		panic(err)
	}

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		h, err := r.Next()
		if err == io.EOF {
			fmt.Println("EOF")
			break
		} else if err != nil {
			panic(err)
		}

		data := make([]byte, h.CapLen)
		_, err = r.Read(data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkReadPcapGoPacket(b *testing.B) {

	h, err := pcap.OpenOffline(file)
	if err != nil {
		panic(err)
	}
	defer h.Close()

	for n := 0; n < b.N; n++ {
		_, _, err := h.ZeroCopyReadPacketData()
		if err == io.EOF {
			fmt.Println("EOF")
			break
		} else if err != nil {
			panic(err)
		}
	}
}

func BenchmarkReadPcapNetboot(b *testing.B) {

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r, err := npcap.NewReader(f)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		if !r.Next() {
			fmt.Println("EOF")
			break
		}
		_ = r.Packet()
	}
}

func BenchmarkReadPcapMiekg(b *testing.B) {

	h, err := mpcap.OpenOffline(file)
	if err != nil {
		panic(err)
	}
	defer h.Close()

	for n := 0; n < b.N; n++ {
		_, r := h.NextEx()
		if r != 1 {
			fmt.Println("got result != 1", r)
			break
		}
	}
}
