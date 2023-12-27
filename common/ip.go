package common

import (
	"fmt"
	"net"
)

func GetIP() string {
	var ipv6 string
	// 指定要获取信息的网络接口名称
	targetInterfaceName := "eth0" // 替换为你的网络接口名称

	// 获取指定名称的网络接口
	iface, err := net.InterfaceByName(targetInterfaceName)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	// 获取指定网络接口的地址信息
	addrs, err := iface.Addrs()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	// 遍历每个地址
	for _, addr := range addrs {
		// 使用类型断言将 net.Addr 转换为具体类型（例如，net.IPNet）
		switch v := addr.(type) {
		case *net.IPNet:
			// 判断 IP 地址是 IPv6
			if v.IP.To16() != nil && v.IP.To4() == nil {
				// 排除 Link-Local 地址
				if !v.IP.IsLinkLocalUnicast() {
					ipv6 = string(v.IP)
				}
			}
		}
	}
	return ipv6
}
