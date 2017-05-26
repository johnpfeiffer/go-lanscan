# README #
Go LanScan is written in Go (aka <https://golang.org/>)

```
./go-lanscan -help
```

```
# an example using the nc utility as a "server" and using all of the parameters to detect it
nc -lvp 443
./go-lanscan -remote "4.4.4.4" -subnet "/30" -port 443 -verbose
```

This is a simple command line utility to quickly scan the network (default is the /24 subnet) for hosts listening on TCP on a specific port (i.e. port 22)

I solved a problem I had (nmap parameters are terrible to remember, not trivially available on every platform, and can get flagged by malware scans), I hope it helps someone else.  (Oh and it was fun writing the code ;)

Some things I might get to fixing:

- If you have blocked outbound port 80 to 8.8.8.8 then it will probably fail (as it attempts to detect the default outbound network interface), it probably does not need to do this at all or could degrade more gracefully (though with -remote you can at least pass it a different ip address, i.e. your gateway)
- configurable timeouts
- unit tests


Binaries for major platforms will be provided on a best effort basis.
