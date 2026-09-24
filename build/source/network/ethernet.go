package network

const (
	ETH_P_ALL     = 0x0003
	ETH_P_IP      = 0x0800
	ETH_P_ARP     = 0x0806
	ETH_P_IPV6    = 0x86DD
	ETH_ALEN      = 6
	ETH_HLEN      = 14
	ETH_ZLEN      = 60
	ETH_DATA_LEN  = 1500
	ETH_FRAME_LEN = 1514
	ETH_FCS_LEN   = 4
)

type EthernetHeader struct {
	Dest   [ETH_ALEN]byte
	Source [ETH_ALEN]byte
	Type   uint16
}

type EthernetDevice struct {
	Name        string
	MAC         [ETH_ALEN]byte
	MTU         int
	IOBase      uint16
	IRQ         uint8
	Driver      *EthernetDriver
	RxBuffer    []byte
	TxBuffer    []byte
	Stats       EthernetStats
	Up          bool
	Promiscuous bool
}

type EthernetStats struct {
	RxPackets   uint64
	TxPackets   uint64
	RxBytes     uint64
	TxBytes     uint64
	RxErrors    uint64
	TxErrors    uint64
	RxDropped   uint64
	TxDropped   uint64
	Collisions  uint64
}

type EthernetDriver struct {
	Name        string
	Probe       func(*EthernetDevice) bool
	Init        func(*EthernetDevice) error
	Start       func(*EthernetDevice) error
	Stop        func(*EthernetDevice) error
	Transmit    func(*EthernetDevice, []byte) error
	Receive     func(*EthernetDevice) ([]byte, error)
	SetMAC      func(*EthernetDevice, [ETH_ALEN]byte) error
	SetPromisc  func(*EthernetDevice, bool) error
	GetStats    func(*EthernetDevice) EthernetStats
}

var (
	ethernetDevices []*EthernetDevice
	ethernetDrivers []*EthernetDriver
)

func RegisterEthernetDriver(driver *EthernetDriver) {
	ethernetDrivers = append(ethernetDrivers, driver)
}

func ProbeEthernetDevices() {
	for _, driver := range ethernetDrivers {
		dev := &EthernetDevice{
			Driver: driver,
			MTU:    ETH_DATA_LEN,
		}
		if driver.Probe(dev) {
			ethernetDevices = append(ethernetDevices, dev)
		}
	}
}

func GetEthernetDevices() []*EthernetDevice {
	return ethernetDevices
}

func (dev *EthernetDevice) Init() error {
	if dev.Driver != nil && dev.Driver.Init != nil {
		return dev.Driver.Init(dev)
	}
	return nil
}

func (dev *EthernetDevice) Start() error {
	if dev.Driver != nil && dev.Driver.Start != nil {
		dev.Up = true
		return dev.Driver.Start(dev)
	}
	return nil
}

func (dev *EthernetDevice) Stop() error {
	if dev.Driver != nil && dev.Driver.Stop != nil {
		dev.Up = false
		return dev.Driver.Stop(dev)
	}
	return nil
}

func (dev *EthernetDevice) Transmit(data []byte) error {
	if dev.Driver != nil && dev.Driver.Transmit != nil {
		dev.Stats.TxPackets++
		dev.Stats.TxBytes += uint64(len(data))
		return dev.Driver.Transmit(dev, data)
	}
	return nil
}

func (dev *EthernetDevice) Receive() ([]byte, error) {
	if dev.Driver != nil && dev.Driver.Receive != nil {
		data, err := dev.Driver.Receive(dev)
		if err == nil && data != nil {
			dev.Stats.RxPackets++
			dev.Stats.RxBytes += uint64(len(data))
		}
		return data, err
	}
	return nil, nil
}

func (dev *EthernetDevice) SetMAC(mac [ETH_ALEN]byte) error {
	if dev.Driver != nil && dev.Driver.SetMAC != nil {
		dev.MAC = mac
		return dev.Driver.SetMAC(dev, mac)
	}
	return nil
}

func (dev *EthernetDevice) SetPromiscuous(promisc bool) error {
	if dev.Driver != nil && dev.Driver.SetPromisc != nil {
		dev.Promiscuous = promisc
		return dev.Driver.SetPromisc(dev, promisc)
	}
	return nil
}

func (dev *EthernetDevice) GetStats() EthernetStats {
	if dev.Driver != nil && dev.Driver.GetStats != nil {
		return dev.Driver.GetStats(dev)
	}
	return dev.Stats
}