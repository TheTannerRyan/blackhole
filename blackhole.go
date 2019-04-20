// Copyright (c) 2019 Tanner Ryan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/miekg/dns"
)

var (
	blackholeAnswer = os.Getenv("BLACKHOLE_ANSWER")            // Answer to all lookups (IP or "NXDOMAIN")
	blackholeTTL, _ = strconv.Atoi(os.Getenv("BLACKHOLE_TTL")) // TTL of response
	blackholePort   = os.Getenv("BLACKHOLE_PORT")              // Port of DNS server
	logging         = os.Getenv("LOGGING")                     // Enable logging to stdout
)

type handler struct {
	Answer string // Answer to all lookups (IP or "NXDOMAIN")
	TTL    uint32 // TTL of response
}

// ServeDNS is the DNS server handler.
func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	if h.Answer == "NXDOMAIN" {
		msg.SetRcode(r, dns.RcodeNameError)
	} else {
		switch r.Question[0].Qtype {
		case dns.TypeA:
			domain := msg.Question[0].Name
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: h.TTL},
				A:   net.ParseIP(h.Answer),
			})
		}
	}
	w.WriteMsg(&msg)
	if logging == "true" {
		log.Println(msg.Question[0].Name)
	}
}

func main() {
	// check environment variables
	if blackholeAnswer == "" {
		panic("error: \"BLACKHOLE_ANSWER\" not defined")
	}
	if blackholePort == "" {
		panic("error: \"BLACKHOLE_PORT\" not defined")
	}

	// listen on all interfaces
	srv := &dns.Server{Addr: ":" + blackholePort, Net: "udp"}
	srv.Handler = &handler{
		Answer: blackholeAnswer,
		TTL:    uint32(blackholeTTL),
	}
	if err := srv.ListenAndServe(); err != nil {
		panic("failed to start DNS blackhole " + err.Error())
	}
}
