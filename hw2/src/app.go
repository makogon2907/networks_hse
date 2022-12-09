package mtu

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	minMTU     = 0
	maxMTU     = 65536
	headerSize = 28
)

func getFlags(packetSize int) (flags []string) {
	switch runtime.GOOS {
	case "darwin":
		flags = strings.Split(fmt.Sprintf("-D -s %d -c 1", packetSize), " ")
	case "linux":
		flags = strings.Split(fmt.Sprintf("-c 1 -M do -s %d", packetSize), " ")
	case "windows":
		flags = strings.Split(fmt.Sprintf("-n 1 -M do -s %d", packetSize), " ")
	}

	return flags
}

func mtuSuccessful(hostname string, packetSize int) bool {
	flags := getFlags(packetSize)
	flags = append(flags, hostname)

	cmd := exec.Command("ping", flags...)

	fmt.Printf("%s ", cmd.String())
	if err := cmd.Run(); err != nil {
		return false
	} else {
		return true
	}
}

func checkAvailable(hostname string) bool {
	cmd := exec.Command("ping", "-c", "1", hostname)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func checkIcmpEnabled() bool {
	if runtime.GOOS == "linux" {
		fContent, err := ioutil.ReadFile("/proc/sys/net/ipv4/icmp_echo_ignore_all")
		if err != nil {
			return true
		}
		icmp, err := strconv.Atoi(strings.TrimSpace(string(fContent)))
		if err != nil {
			return true
		}
		if icmp == 1 {
			fmt.Println("icmp is blocked in system settings: /proc/sys/net/ipv4/icmp_echo_ignore_all")
			return false
		}
	}
	return true
}

func FindMTUWith(hostname string) {
	if !checkIcmpEnabled() {
		log.Fatalf("ICMP is disabled :(")
	}

	addr, err := net.LookupIP(hostname)
	if err != nil {
		log.Fatalf("Unknown host: %q", hostname)
	} else {
		fmt.Printf("IP addresses: %s\n", addr)
	}

	if hostAvailable := checkAvailable(hostname); hostAvailable {
		fmt.Println("Host available, start pinging...")
	} else {
		log.Fatalf("Host unavailable :(")
	}

	left := minMTU
	right := maxMTU - headerSize

	for left < right {
		mtuToTest := (left + right + 1) / 2

		if mtuSuccessful(hostname, mtuToTest) {
			left = mtuToTest
			fmt.Printf("...OK\n")
		} else {
			right = mtuToTest - 1
			fmt.Printf("...FAILED\n")
		}
	}

	fmt.Printf("FOUND MTU TO HOST %q: %d bytes\n", hostname, left+headerSize)
}
