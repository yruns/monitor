package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	// 打开网络接口
	handle, err := pcap.OpenLive("eth0", 65535, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// 设置过滤器
	filter := "port 3000"
	if err := handle.SetBPFFilter(filter); err != nil {
		log.Fatal(err)
	}

	// 开始捕获数据包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			fmt.Printf("Source IP: %s\n", packet.NetworkLayer().NetworkFlow().Src().String())
			fmt.Printf("Destination IP: %s\n", packet.NetworkLayer().NetworkFlow().Dst().String())
			fmt.Printf("Source Port: %d\n", tcp.SrcPort)
			fmt.Printf("Destination Port: %d\n", tcp.DstPort)
		}
	}
}
