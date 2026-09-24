package network

const (
	IEEE80211_ADDR_LEN = 6
	IEEE80211_MAX_SSID_LEN = 32
	IEEE80211_MAX_BSSID_LEN = 6
	IEEE80211_MAX_RATES = 16
	IEEE80211_MAX_CHANNELS = 14
)

type WiFiMode uint8

const (
	WiFiModeNone WiFiMode = iota
	WiFiModeStation
	WiFiModeAP
	WiFiModeMonitor
	WiFiModeMesh
)

type WiFiSecurity uint8

const (
	WiFiSecurityOpen WiFiSecurity = iota
	WiFiSecurityWEP
	WiFiSecurityWPA
	WiFiSecurityWPA2
	WiFiSecurityWPA3
)

type WiFiDevice struct {
	Name         string
	MAC          [IEEE80211_ADDR_LEN]byte
	Mode         WiFiMode
	SSID         [IEEE80211_MAX_SSID_LEN]byte
	BSSID        [IEEE80211_MAX_BSSID_LEN]byte
	Channel      uint8
	Frequency    uint32
	Security     WiFiSecurity
	Password     [64]byte
	Driver       *WiFiDriver
	RxBuffer     []byte
	TxBuffer     []byte
	Stats        WiFiStats
	Up           bool
	Connected    bool
	SignalStrength int8
	Bitrate      uint32
}

type WiFiStats struct {
	RxPackets   uint64
	TxPackets   uint64
	RxBytes     uint64
	TxBytes     uint64
	RxErrors    uint64
	TxErrors    uint64
	RxDropped   uint64
	TxDropped   uint64
	Retries     uint64
}

type WiFiDriver struct {
	Name        string
	Probe       func(*WiFiDevice) bool
	Init        func(*WiFiDevice) error
	Start       func(*WiFiDevice) error
	Stop        func(*WiFiDevice) error
	Transmit    func(*WiFiDevice, []byte) error
	Receive     func(*WiFiDevice) ([]byte, error)
	SetMode     func(*WiFiDevice, WiFiMode) error
	SetChannel  func(*WiFiDevice, uint8) error
	Scan        func(*WiFiDevice) ([]*AccessPoint, error)
	Connect     func(*WiFiDevice, string, string, WiFiSecurity) error
	Disconnect  func(*WiFiDevice) error
	GetStats    func(*WiFiDevice) WiFiStats
}

type AccessPoint struct {
	SSID         string
	BSSID        [IEEE80211_MAX_BSSID_LEN]byte
	Channel      uint8
	Frequency    uint32
	Security     WiFiSecurity
	SignalStrength int8
	Bitrates     []uint32
	Capabilities uint16
}

var (
	wifiDevices []*WiFiDevice
	wifiDrivers []*WiFiDriver
)

func RegisterWiFiDriver(driver *WiFiDriver) {
	wifiDrivers = append(wifiDrivers, driver)
}

func ProbeWiFiDevices() {
	for _, driver := range wifiDrivers {
		dev := &WiFiDevice{
			Driver: driver,
		}
		if driver.Probe(dev) {
			wifiDevices = append(wifiDevices, dev)
		}
	}
}

func GetWiFiDevices() []*WiFiDevice {
	return wifiDevices
}

func (dev *WiFiDevice) Init() error {
	if dev.Driver != nil && dev.Driver.Init != nil {
		return dev.Driver.Init(dev)
	}
	return nil
}

func (dev *WiFiDevice) Start() error {
	if dev.Driver != nil && dev.Driver.Start != nil {
		dev.Up = true
		return dev.Driver.Start(dev)
	}
	return nil
}

func (dev *WiFiDevice) Stop() error {
	if dev.Driver != nil && dev.Driver.Stop != nil {
		dev.Up = false
		return dev.Driver.Stop(dev)
	}
	return nil
}

func (dev *WiFiDevice) SetMode(mode WiFiMode) error {
	if dev.Driver != nil && dev.Driver.SetMode != nil {
		dev.Mode = mode
		return dev.Driver.SetMode(dev, mode)
	}
	return nil
}

func (dev *WiFiDevice) SetChannel(channel uint8) error {
	if dev.Driver != nil && dev.Driver.SetChannel != nil {
		dev.Channel = channel
		return dev.Driver.SetChannel(dev, channel)
	}
	return nil
}

func (dev *WiFiDevice) Scan() ([]*AccessPoint, error) {
	if dev.Driver != nil && dev.Driver.Scan != nil {
		return dev.Driver.Scan(dev)
	}
	return nil, nil
}

func (dev *WiFiDevice) Connect(ssid, password string, security WiFiSecurity) error {
	if dev.Driver != nil && dev.Driver.Connect != nil {
		copy(dev.SSID[:], ssid)
		copy(dev.Password[:], password)
		dev.Security = security
		err := dev.Driver.Connect(dev, ssid, password, security)
		if err == nil {
			dev.Connected = true
		}
		return err
	}
	return nil
}

func (dev *WiFiDevice) Disconnect() error {
	if dev.Driver != nil && dev.Driver.Disconnect != nil {
		dev.Connected = false
		return dev.Driver.Disconnect(dev)
	}
	return nil
}

func (dev *WiFiDevice) GetStats() WiFiStats {
	if dev.Driver != nil && dev.Driver.GetStats != nil {
		return dev.Driver.GetStats(dev)
	}
	return dev.Stats
}

func (dev *WiFiDevice) Transmit(data []byte) error {
	if dev.Driver != nil && dev.Driver.Transmit != nil {
		dev.Stats.TxPackets++
		dev.Stats.TxBytes += uint64(len(data))
		return dev.Driver.Transmit(dev, data)
	}
	return nil
}

func (dev *WiFiDevice) Receive() ([]byte, error) {
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