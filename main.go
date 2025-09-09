package main

import (
	"fmt"
	"net"
	"os"
	"text/template"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/beevik/ntp"
)

const (
	defaultTemplate = "{{if .Validate}}{{.Validate}}{{else}}OK{{end}} from {{.Address}}:" +
		"seq={{.Seq}} stratum={{.Stratum}} " +
		"offset={{.ClockOffset}} distance={{.RootDistance}} RTT={{.RTT}} ref={{.ReferenceString}}\n"
	version = "dev"
)

type PrettyRsp struct {
	ntp.Response
	Address string
	Seq     int
}

type args struct {
	Address string        `arg:"positional,required"`
	Timeout time.Duration `arg:"-t,--timeout" help:"timeout duration" default:"3s"`
	Port    int           `arg:"-p,--port" help:"NTP server port" default:"123"`
	Count   int           `arg:"-c,--count" help:"stop after COUNT replies, if only 1 count, it act like ping" default:"16"`
	IPv6    bool          `arg:"-6,--ipv6" help:"prefer IPv6" default:"false"`
}

func (args) Version() string {
	return version
}

func resolveToIP(s string, ipv6 bool) (net.IP, error) {
	ip := net.ParseIP(s)
	if ip != nil {
		return ip, nil
	}
	network := "ip4"
	if ipv6 {
		network = "ip6"
	}
	ipa, err := net.ResolveIPAddr(network, s)
	if err != nil {
		return nil, err
	}
	return ipa.IP, nil
}

func main() {

	var args args

	arg.MustParse(&args)

	raddr, err := resolveToIP(args.Address, args.IPv6)
	if err != nil {
		fmt.Printf("error from %s:%s\n", args.Address, err)
		os.Exit(1)
	}
	addr := raddr.String()

	if args.Port != 123 {
		addr = fmt.Sprintf("%s:%d", args.Address, args.Port)
	}

	tmpl := template.Must(template.New("").Parse(defaultTemplate))

	for i := 0; i < args.Count; i++ {
		rsp, err := ntp.QueryWithOptions(addr,
			ntp.QueryOptions{Timeout: args.Timeout})

		if err != nil {
			fmt.Printf("error from %s:%s\n", addr, err)
			break
		}

		err = tmpl.Execute(os.Stdout, &PrettyRsp{Response: *rsp,
			Address: addr, Seq: i + 1})
		if err != nil {
			fmt.Printf("error from %s:%s\n", addr, err)
			break
		}

		if rsp.Stratum == 1 {
			break
		}

		addr = rsp.ReferenceString()
		if net.ParseIP(addr) == nil {
			fmt.Printf("error from %s:invalid ip address\n", addr)
			break
		}
	}
}
