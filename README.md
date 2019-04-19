# blackhole - lightweight lying DNS server
[![GitHub
license](https://img.shields.io/github/license/TheTannerRyan/blackhole.svg?style=flat-square)](https://github.com/TheTannerRyan/blackhole/blob/master/LICENSE)

blackhole is a lightweight DNS service. It will redirect all DNS lookups to a
provided IP. 

## Docker Usage
There's nothing funky about the process. Clone, build, and run to your heart's
desire.
```
docker build -t blackhole .
docker run -d --restart=unless-stopped \
  -p 9999:9999/udp \
  -e BLACKHOLE_IP=127.0.0.1 \
  -e BLACKHOLE_TTL=3600 \
  -e BLACKHOLE_PORT=9999 \
  -e LOGGING=true \
  --name blackhole \
  blackhole
```
You are able to configure three environment variables:
- `BLACKHOLE_IP`: the IP that all DNS lookups will return.
- `BLACKHOLE_TTL`: the TTL that all DNS lookups will return.
- `BLACKHOLE_PORT`: the UDP port to listen for DNS requests on.
- `LOGGING`: boolean to enable stdout logging of domains.

## Response
All __A record__ lookups will return the `BLACKHOLE_IP` with the
`BLACKHOLE_TTL`:
```
$ dig google.com @172.16.0.18 -p 9999

; <<>> DiG 9.10.6 <<>> google.com @172.16.0.18 -p 9999
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 30094
;; flags: qr rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;google.com.                    IN      A

;; ANSWER SECTION:
google.com.             3600    IN      A       127.0.0.1

;; Query time: 1 msec
;; SERVER: 172.16.0.18#9999(172.16.0.18)
;; WHEN: Fri Apr 19 18:07:16 EDT 2019
;; MSG SIZE  rcvd: 54
```

Any lookup other than __A record__ will return no recursion available.

## Purpose
I built this for another project I was working on. When dealing with millions of
RPZ zones, it was much faster and more efficient to redirect lookups to this
blackholing DNS server, rather than creating an A record for each blacklisted
record.

If you do find some other purpose for this, please let me know.

## License
Copyright (c) 2019 Tanner Ryan. All rights reserved. Use of this source code is
governed by a BSD-style license that can be found in the LICENSE file.
