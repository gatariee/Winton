package utils

import (
	"fmt"
	"net"
	"runtime"

)

func fallback() (string) {
	addrs, err := net.InterfaceAddrs()
    if err != nil {
        return "-"
    }

    for _, addr := range addrs {
        ipNet, ok := addr.(*net.IPNet)
        if ok && !ipNet.IP.IsLoopback() {
            if ipNet.IP.To4() != nil {
                return ipNet.IP.String()
            }
        }
    }

    return "-"
}

func GetInternalIP() (string, error) {
    interfaces, err := net.Interfaces()
    if err != nil {
        return "", err
    }

    for _, iface := range interfaces {
        if iface.Name == "Ethernet" {
            addrs, err := iface.Addrs()
            if err != nil {
                return "", err
            }

            for _, addr := range addrs {
                ipNet, ok := addr.(*net.IPNet)
                if ok && !ipNet.IP.IsLoopback() {
                    if ipNet.IP.To4() != nil {
                        return ipNet.IP.String(), nil
                    }
                }
            }
        }
    }

	return fallback(), nil
}

func GetSystemInfo() string {
    os := runtime.GOOS
    arch := runtime.GOARCH

    switch os {
    case "darwin":
        os = "macOS"
    case "windows":
        os = "Windows"
    case "linux":
        os = "Linux"
    }

    switch arch {
    case "amd64":
        arch = "x64"
    case "386":
        arch = "x86"
    }

    return fmt.Sprintf("%s (%s)", os, arch)
}