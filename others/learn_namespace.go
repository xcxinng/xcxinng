package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"

	"github.com/vishvananda/netns"
)

var (
	ns    = flag.String("ns", "netns1", "specify ns to connect")
	iface = flag.String("iface", "eth0", "specify container interface")
)

func learnNamespace() {
	flag.Parse()
	runtime.LockOSThread()
	net.Interfaces()

	currentNs, err := netns.Get()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = netns.Set(currentNs)
		if err != nil {
			log.Fatal(err)
		}
	}()

	nsHandle, err := netns.GetFromName(*ns)
	if err != nil {
		log.Fatal("get ns error:", err)
	}
	if err = netns.Set(nsHandle); err != nil {
		log.Fatal("connect ns error:", err)
	} else {
		fmt.Printf("[debug]connect to netns:%s success\n", *ns)
	}

	c := exec.Command("/usr/sbin/ifconfig", *iface)
	output, err := c.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("netns:[%s],interface:[%s]\n\rifconfig_result:\n\t%s\n", *ns, *iface, output)
}
