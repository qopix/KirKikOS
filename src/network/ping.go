package network

import (
	"kirkikos/arch/amd64"
)

const (
	ICMP_ECHO_REQUEST = 8
	ICMP_ECHO_REPLY   = 0
	ICMP_DEST_UNREACH = 3
	ICMP_TIME_EXCEEDED = 11
	ICMP_PARAM_PROB   = 12
	
	IPPROTO_ICMP = 1
	IPPROTO_TCP  = 6
	IPPROTO_UDP  = 17
	
	IP_DF = 0x4000
	IP_MF = 0x2000
	IP_OFFSET = 0x1FFF
)

type ICMPHeader struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	ID       uint16
	Seq      uint16
}

type IPHeader struct {
	VersionIHL uint8
	TOS        uint8
	TotalLen   uint16
	ID         uint16
	FlagsFrag  uint16
	TTL        uint8
	Protocol   uint8
	Checksum   uint16
	SrcIP      [4]byte
	DstIP      [4]byte
}

var (
	icmpID   uint16 = 0x1234
	icmpSeq  uint16 = 1
)

func Ping(host string, count int, timeout int) error {
	ip := resolveHost(host)
	if ip == [4]byte{0, 0, 0, 0} {
		return ErrHostUnreachable
	}
	
	printString("PING ")
	printString(host)
	printString(" (")
	printIP(ip)
	printString("): ")
	printUint(uint64(64), 0)
	printString(" bytes of data.\n")
	
	sent := 0
	received := 0
	minRTT := uint64(0)
	maxRTT := uint64(0)
	totalRTT := uint64(0)
	
	for i := 0; i < count; i++ {
		rtt, err := sendICMPEcho(ip)
		if err != nil {
			printString("Request timeout for icmp_seq ")
			printUint(uint64(i+1), 0)
			printString("\n")
		} else {
			received++
			printString("64 bytes from ")
			printIP(ip)
			printString(": icmp_seq=")
			printUint(uint64(i+1), 0)
			printString(" ttl=64 time=")
			printUint(rtt, 0)
			printString(" ms\n")
			
			if minRTT == 0 || rtt < minRTT {
				minRTT = rtt
			}
			if rtt > maxRTT {
				maxRTT = rtt
			}
			totalRTT += rtt
		}
		sent++
		
		if i < count-1 {
			amd64.Sleep(timeout)
		}
	}
	
	printString("\n--- ")
	printString(host)
	printString(" ping statistics ---\n")
	printUint(uint64(sent), 0)
	printString(" packets transmitted, ")
	printUint(uint64(received), 0)
	printString(" received, ")
	if sent > 0 {
		loss := ((sent - received) * 100) / sent
		printUint(uint64(loss), 0)
	} else {
		printString("100")
	}
	printString("% packet loss\n")
	
	if received > 0 {
		avgRTT := totalRTT / uint64(received)
		printString("rtt min/avg/max = ")
		printUint(minRTT, 0)
		printString("/")
		printUint(avgRTT, 0)
		printString("/")
		printUint(maxRTT, 0)
		printString(" ms\n")
	}
	
	return nil
}

func sendICMPEcho(dstIP [4]byte) (uint64, error) {
	icmpSeq++
	
	icmp := ICMPHeader{
		Type:     ICMP_ECHO_REQUEST,
		Code:     0,
		Checksum: 0,
		ID:       icmpID,
		Seq:      icmpSeq,
	}
	
	payload := make([]byte, 56)
	for i := 0; i < len(payload); i++ {
		payload[i] = byte(i & 0xFF)
	}
	
	packet := make([]byte, 8+56)
	packet[0] = icmp.Type
	packet[1] = icmp.Code
	packet[2] = 0
	packet[3] = 0
	packet[4] = byte(icmp.ID >> 8)
	packet[5] = byte(icmp.ID & 0xFF)
	packet[6] = byte(icmp.Seq >> 8)
	packet[7] = byte(icmp.Seq & 0xFF)
	copy(packet[8:], payload)
	
	checksum := calculateChecksum(packet)
	packet[2] = byte(checksum >> 8)
	packet[3] = byte(checksum & 0xFF)
	
	ipPacket := buildIPPacket(dstIP, packet)
	
	start := amd64.GetTimerTicks()
	err := sendRawIP(ipPacket)
	if err != nil {
		return 0, err
	}
	
	reply, err := receiveICMPEcho(dstIP, timeout)
	if err != nil {
		return 0, err
	}
	end := amd64.GetTimerTicks()
	
	if reply != nil {
		return uint64(end - start), nil
	}
	
	return 0, ErrTimeout
}

func buildIPPacket(dstIP [4]byte, payload []byte) []byte {
	packet := make([]byte, 20+len(payload))
	
	packet[0] = 0x45
	packet[1] = 0
	packet[2] = byte((20 + len(payload)) >> 8)
	packet[3] = byte((20 + len(payload)) & 0xFF)
	
	id := amd64.GetTimerTicks() & 0xFFFF
	packet[4] = byte(id >> 8)
	packet[5] = byte(id & 0xFF)
	
	packet[6] = byte(IP_DF >> 8)
	packet[7] = byte(IP_DF & 0xFF)
	
	packet[8] = 64
	packet[9] = IPPROTO_ICMP
	
	packet[10] = 0
	packet[11] = 0
	
	copy(packet[12:16], getLocalIP()[:])
	copy(packet[16:20], dstIP[:])
	
	ipChecksum := calculateChecksum(packet[:20])
	packet[10] = byte(ipChecksum >> 8)
	packet[11] = byte(ipChecksum & 0xFF)
	
	copy(packet[20:], payload)
	
	return packet
}

func calculateChecksum(data []byte) uint16 {
	sum := uint32(0)
	for i := 0; i < len(data); i += 2 {
		if i+1 < len(data) {
			sum += uint32(data[i])<<8 | uint32(data[i+1])
		} else {
			sum += uint32(data[i]) << 8
		}
	}
	for (sum >> 16) != 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return uint16(^sum)
}

func getLocalIP() [4]byte {
	return [4]byte{10, 0, 2, 15}
}

func resolveHost(host string) [4]byte {
	if host == "localhost" || host == "127.0.0.1" {
		return [4]byte{127, 0, 0, 1}
	}
	
	parts := splitIP(host)
	if len(parts) == 4 {
		var ip [4]byte
		for i, p := range parts {
			ip[i] = byte(p)
		}
		return ip
	}
	
	return [4]byte{0, 0, 0, 0}
}

func splitIP(s string) []int {
	var parts []int
	current := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			parts = append(parts, current)
			current = 0
		} else if s[i] >= '0' && s[i] <= '9' {
			current = current*10 + int(s[i]-'0')
		}
	}
	parts = append(parts, current)
	return parts
}

func printIP(ip [4]byte) {
	for i := 0; i < 4; i++ {
		if i > 0 {
			printString(".")
		}
		printUint(uint64(ip[i]), 0)
	}
}

func sendRawIP(packet []byte) error {
	return nil
}

func receiveICMPEcho(dstIP [4]byte, timeout int) ([]byte, error) {
	return nil, ErrTimeout
}

var (
	ErrHostUnreachable = &NetworkError{"host unreachable"}
	ErrTimeout         = &NetworkError{"timeout"}
	ErrNetworkUnreachable = &NetworkError{"network unreachable"}
)

type NetworkError struct {
	msg string
}

func (e *NetworkError) Error() string {
	return e.msg
}

func printString(s string) {
	for i := 0; i < len(s); i++ {
		amd64.PutChar(s[i])
	}
}

func printUint(n uint64, width int) {
	var buf [20]byte
	i := len(buf)
	if n == 0 {
		i--
		buf[i] = '0'
	} else {
		for n > 0 && i > 0 {
			i--
			buf[i] = byte('0' + n%10)
			n /= 10
		}
	}
	for j := i; j < len(buf); j++ {
		printString(string(buf[j:j+1]))
	}
}