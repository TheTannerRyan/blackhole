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
	blackholeIP     = os.Getenv("BLACKHOLE_IP")                // Answer to all lookups
	blackholeTTL, _ = strconv.Atoi(os.Getenv("BLACKHOLE_TTL")) // TTL of response
	blackholePort   = os.Getenv("BLACKHOLE_PORT")              // Port of DNS server
	logging         = os.Getenv("LOGGING")                     // Enable logging to stdout
)

type handler struct {
	IP  net.IP // Answer to all lookups
	TTL uint32 // TTL of response
}

// ServeDNS is the DNS server handler.
func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		domain := msg.Question[0].Name
		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: h.TTL},
			A:   h.IP,
		})
	}
	w.WriteMsg(&msg)
	if logging == "true" {
		log.Println(msg.Question[0].Name)
	}
}

func main() {
	// check environment variables
	if blackholeIP == "" {
		panic("error: \"BLACKHOLE_IP\" not defined")
	}
	if blackholePort == "" {
		panic("error: \"BLACKHOLE_PORT\" not defined")
	}

	// listen on all interfaces
	srv := &dns.Server{Addr: ":" + blackholePort, Net: "udp"}
	srv.Handler = &handler{
		IP:  net.ParseIP(blackholeIP),
		TTL: uint32(blackholeTTL),
	}
	if err := srv.ListenAndServe(); err != nil {
		panic("failed to start DNS blackhole " + err.Error())
	}
}
