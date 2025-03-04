//go:build windows

package runtime

import (
	"golang.org/x/sys/windows"
	"syscall"
)

var (
	iphlpapi     = syscall.NewLazyDLL("iphlpapi.dll")
	hGetIfTable2 = iphlpapi.NewProc("GetIfTable2")

	wlanapi             = syscall.NewLazyDLL("wlanapi.dll")
	hWlanOpenHandle     = wlanapi.NewProc("WlanOpenHandle")
	hWlanCloseHandle    = wlanapi.NewProc("WlanCloseHandle")
	hWlanQueryInterface = wlanapi.NewProc("WlanQueryInterface")
	hWlanEnumInterfaces = wlanapi.NewProc("WlanEnumInterfaces")
)

type MIN_IF_TABLE2 struct {
	NumEntries uint64
	Table      [256]MIB_IF_ROW2
}

const (
	IF_MAX_STRING_SIZE         uint64 = 256
	IF_MAX_PHYS_ADDRESS_LENGTH uint64 = 32
)

type MIB_IF_ROW2 struct {
	InterfaceLuid            uint64
	InterfaceIndex           uint32
	InterfaceGUID            windows.GUID
	Alias                    [IF_MAX_STRING_SIZE + 1]uint16
	Description              [IF_MAX_STRING_SIZE + 1]uint16
	PhysicalAddressLength    uint32
	PhysicalAddress          [IF_MAX_PHYS_ADDRESS_LENGTH]uint8
	PermanentPhysicalAddress [IF_MAX_PHYS_ADDRESS_LENGTH]uint8

	Mtu                uint32
	Type               uint32
	TunnelType         uint32
	MediaType          uint32
	PhysicalMediumType uint32
	AccessType         uint32
	DirectionType      uint32

	InterfaceAndOperStatusFlags struct {
		HardwareInterface bool
		FilterInterface   bool
		ConnectorPresent  bool
		NotAuthenticated  bool
		NotMediaConnected bool
		Paused            bool
		LowPower          bool
		EndPointInterface bool
	}

	OperStatus        uint32
	AdminStatus       uint32
	MediaConnectState uint32
	NetworkGUID       windows.GUID
	ConnectionType    uint32

	TransmitLinkSpeed uint64
	ReceiveLinkSpeed  uint64

	InOctets           uint64
	InUcastPkts        uint64
	InNUcastPkts       uint64
	InDiscards         uint64
	InErrors           uint64
	InUnknownProtos    uint64
	InUcastOctets      uint64
	InMulticastOctets  uint64
	InBroadcastOctets  uint64
	OutOctets          uint64
	OutUcastPkts       uint64
	OutNUcastPkts      uint64
	OutDiscards        uint64
	OutErrors          uint64
	OutUcastOctets     uint64
	OutMulticastOctets uint64
	OutBroadcastOctets uint64
	OutQLen            uint64
}

type WLAN_INTERFACE_INFO_LIST struct {
	dwNumberOfItems uint32
	dwIndex         uint32
	InterfaceInfo   [1]WLAN_INTERFACE_INFO
}

type WLAN_INTERFACE_INFO struct {
	InterfaceGuid           syscall.GUID
	strInterfaceDescription [256]uint16
	isState                 uint32
}

const (
	WLAN_MAX_NAME_LENGTH  int64 = 256
	DOT11_SSID_MAX_LENGTH int64 = 32
)

type WLAN_CONNECTION_ATTRIBUTES struct {
	isState                   uint32
	wlanConnectionMode        uint32
	strProfileName            [WLAN_MAX_NAME_LENGTH]uint16
	wlanAssociationAttributes WLAN_ASSOCIATION_ATTRIBUTES
	wlanSecurityAttributes    WLAN_SECURITY_ATTRIBUTES
}

type WLAN_ASSOCIATION_ATTRIBUTES struct {
	dot11Ssid         DOT11_SSID
	dot11BssType      uint32
	dot11Bssid        [6]uint8
	dot11PhyType      uint32
	uDot11PhyIndex    uint32
	wlanSignalQuality uint32
	ulRxRate          uint32
	ulTxRate          uint32
}

type WLAN_SECURITY_ATTRIBUTES struct {
	bSecurityEnabled     uint32
	bOneXEnabled         uint32
	dot11AuthAlgorithm   uint32
	dot11CipherAlgorithm uint32
}

type DOT11_SSID struct {
	uSSIDLength uint32
	ucSSID      [DOT11_SSID_MAX_LENGTH]uint8
}

type IFTYPE string
type NDIS_MEDIUM string
type NDIS_PHYSICAL_MEDIUM string

type WifiType string

type NetworkInfo struct {
	Alias                 string
	Interface             string
	InterfaceType         IFTYPE
	NDISMediaType         NDIS_MEDIUM
	NDISPhysicalMeidaType NDIS_PHYSICAL_MEDIUM
	TransmitLinkSpeed     uint64
	ReceiveLinkSpeed      uint64
	SSID                  string // Wi-Fi only
}

type WifiInfo struct {
	SSID           string
	Interface      string
	RadioType      WifiType
	PhysType       WifiType
	Authentication WifiType
	Cipher         WifiType
	Channel        int
	ReceiveRate    int
	TransmitRate   int
	Signal         int
	Error          string
}

const (
	// see https://docs.microsoft.com/en-us/windows/win32/api/netioapi/ns-netioapi-mib_if_row2

	// InterfaceType
	IF_TYPE_OTHER              IFTYPE = "Other"            // 1
	IF_TYPE_ETHERNET_CSMACD    IFTYPE = "Ethernet/802.3"   // 6
	IF_TYPE_ISO88025_TOKENRING IFTYPE = "Token Ring/802.5" // 9
	IF_TYPE_FDDI               IFTYPE = "FDDI"             // 15
	IF_TYPE_PPP                IFTYPE = "PPP"              // 23
	IF_TYPE_SOFTWARE_LOOPBACK  IFTYPE = "Loopback"         // 24
	IF_TYPE_ATM                IFTYPE = "ATM"              // 37
	IF_TYPE_IEEE80211          IFTYPE = "Wi-Fi/802.11"     // 71
	IF_TYPE_TUNNEL             IFTYPE = "Tunnel"           // 131
	IF_TYPE_IEEE1394           IFTYPE = "FireWire/1394"    // 144
	IF_TYPE_IEEE80216_WMAN     IFTYPE = "WMAN/802.16"      // 237 WiMax
	IF_TYPE_WWANPP             IFTYPE = "WWANPP/GSM"       // 243 GSM
	IF_TYPE_WWANPP2            IFTYPE = "WWANPP/CDMA"      // 244 CDMA
	IF_TYPE_UNKNOWN            IFTYPE = "Unknown"

	// NDISMediaType
	NdisMedium802_3        NDIS_MEDIUM = "802.3"         // 0
	NdisMedium802_5        NDIS_MEDIUM = "802.5"         // 1
	NdisMediumFddi         NDIS_MEDIUM = "FDDI"          // 2
	NdisMediumWan          NDIS_MEDIUM = "WAN"           // 3
	NdisMediumLocalTalk    NDIS_MEDIUM = "LocalTalk"     // 4
	NdisMediumDix          NDIS_MEDIUM = "DIX"           // 5
	NdisMediumArcnetRaw    NDIS_MEDIUM = "ARCNET"        // 6
	NdisMediumArcnet878_2  NDIS_MEDIUM = "ARCNET(878.2)" // 7
	NdisMediumAtm          NDIS_MEDIUM = "ATM"           // 8
	NdisMediumWirelessWan  NDIS_MEDIUM = "WWAN"          // 9
	NdisMediumIrda         NDIS_MEDIUM = "IrDA"          // 10
	NdisMediumBpc          NDIS_MEDIUM = "Broadcast"     // 11
	NdisMediumCoWan        NDIS_MEDIUM = "CO WAN"        // 12
	NdisMedium1394         NDIS_MEDIUM = "1394"          // 13
	NdisMediumInfiniBand   NDIS_MEDIUM = "InfiniBand"    // 14
	NdisMediumTunnel       NDIS_MEDIUM = "Tunnel"        // 15
	NdisMediumNative802_11 NDIS_MEDIUM = "Native 802.11" // 16
	NdisMediumLoopback     NDIS_MEDIUM = "Loopback"      // 17
	NdisMediumWiMax        NDIS_MEDIUM = "WiMax"         // 18
	NdisMediumUnknown      NDIS_MEDIUM = "Unknown"

	// NDISPhysicalMeidaType
	NdisPhysicalMediumUnspecified  NDIS_PHYSICAL_MEDIUM = "Unspecified"   // 0
	NdisPhysicalMediumWirelessLan  NDIS_PHYSICAL_MEDIUM = "Wireless LAN"  // 1
	NdisPhysicalMediumCableModem   NDIS_PHYSICAL_MEDIUM = "Cable Modem"   // 2
	NdisPhysicalMediumPhoneLine    NDIS_PHYSICAL_MEDIUM = "Phone Line"    // 3
	NdisPhysicalMediumPowerLine    NDIS_PHYSICAL_MEDIUM = "Power Line"    // 4
	NdisPhysicalMediumDSL          NDIS_PHYSICAL_MEDIUM = "DSL"           // 5
	NdisPhysicalMediumFibreChannel NDIS_PHYSICAL_MEDIUM = "Fibre Channel" // 6
	NdisPhysicalMedium1394         NDIS_PHYSICAL_MEDIUM = "1394"          // 7
	NdisPhysicalMediumWirelessWan  NDIS_PHYSICAL_MEDIUM = "Wireless WAN"  // 8
	NdisPhysicalMediumNative802_11 NDIS_PHYSICAL_MEDIUM = "Native 802.11" // 9
	NdisPhysicalMediumBluetooth    NDIS_PHYSICAL_MEDIUM = "Bluetooth"     // 10
	NdisPhysicalMediumInfiniband   NDIS_PHYSICAL_MEDIUM = "Infini Band"   // 11
	NdisPhysicalMediumWiMax        NDIS_PHYSICAL_MEDIUM = "WiMax"         // 12
	NdisPhysicalMediumUWB          NDIS_PHYSICAL_MEDIUM = "UWB"           // 13
	NdisPhysicalMedium802_3        NDIS_PHYSICAL_MEDIUM = "802.3"         // 14
	NdisPhysicalMedium802_5        NDIS_PHYSICAL_MEDIUM = "802.5"         // 15
	NdisPhysicalMediumIrda         NDIS_PHYSICAL_MEDIUM = "IrDA"          // 16
	NdisPhysicalMediumWiredWAN     NDIS_PHYSICAL_MEDIUM = "Wired WAN"     // 17
	NdisPhysicalMediumWiredCoWan   NDIS_PHYSICAL_MEDIUM = "Wired CO WAN"  // 18
	NdisPhysicalMediumOther        NDIS_PHYSICAL_MEDIUM = "Other"         // 19
	NdisPhysicalMediumUnknown      NDIS_PHYSICAL_MEDIUM = "Unknown"
)

const (
	FHSS   WifiType = "FHSS"
	DSSS   WifiType = "DSSS"
	IR     WifiType = "IR"
	A      WifiType = "802.11a"
	HRDSSS WifiType = "HRDSSS"
	G      WifiType = "802.11g"
	N      WifiType = "802.11n"
	AC     WifiType = "802.11ac"

	Infrastructure WifiType = "Infrastructure"
	Independent    WifiType = "Independent"
	Any            WifiType = "Any"

	OpenSystem WifiType = "802.11 Open System"
	SharedKey  WifiType = "802.11 Shared Key"
	WPA        WifiType = "WPA"
	WPAPSK     WifiType = "WPA PSK"
	WPANone    WifiType = "WPA NONE"
	WPA2       WifiType = "WPA2"
	WPA2PSK    WifiType = "WPA2 PSK"
	Disabled   WifiType = "disabled"

	None   WifiType = "None"
	WEP40  WifiType = "WEP40"
	TKIP   WifiType = "TKIP"
	CCMP   WifiType = "CCMP"
	WEP104 WifiType = "WEP104"
	WEP    WifiType = "WEP"
)
