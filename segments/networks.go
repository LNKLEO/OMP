package segments

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/LNKLEO/OMP/properties"
	"github.com/LNKLEO/OMP/runtime"
)

type Networks struct {
	base

	Error string

	Networks     string
	Status       string
}

type Unit string

const (
	Auto Unit = "Auto"
	A    Unit = "A"
	Hide Unit = "Hide"
	bps  Unit = "bps"
	b    Unit = ""
	Kbps Unit = "Kbps"
	K    Unit = "K"
	Mbps Unit = "Mbps"
	M    Unit = "M"
	Gbps Unit = "Gbps"
	G    Unit = "G"
	Tbps Unit = "Tbps"
	T    Unit = "T"
)

func (n *Networks) Template() string {
	return "{{ if eq n.Status \"Connected\" }} {{ .Networks }} | {{ n.IconConnected }} {{ else }} {{ n.IconDisconnected }}"
}

func (n *Networks) Enabled() bool {
	// This segment only supports Windows/WSL for now
	if n.env.Platform() != runtime.WINDOWS && !n.env.IsWsl() {
		return false
	}
	Spliter := n.props.GetString("Spliter", "|")
	// IconConnected := n.props.GetString("IconConnected", "")
	// IconDisconnected := n.props.GetString("IconDisconnected", "")
	connections, err := n.env.Connection()
	displayError := n.props.GetBool(properties.DisplayError, false)
	if err != nil && displayError {
		n.Error = err.Error()
		return true
	}
	if err != nil || connections == nil {
		return false
	}
	if len(connections) == 0 {
		n.Status = "Disconnected"
	} else {
		n.Status = "Connected"
		connectionstrs := make([]string, 0)
		for _, connection := range connections {
			connectionstrs = append(connectionstrs, n.ConstructConnectionInfo(connection))
		}
		n.Networks = strings.Join(connectionstrs, Spliter)
	}
	return true
}

func (n *Networks) Init(props properties.Properties, env runtime.Environment) {
	n.props = props
	n.env = env
}

func (n *Networks) ConstructConnectionInfo(connection *runtime.Connection) string {
	str := ""
	IconMap := make(map[string]string)
	IconMap["Ethernet"] = "󰛳"
	IconMap["Wi-Fi"] = "󰖩"
	IconMap["Bluetooth"] = "󰂴"
	IconMap["Cellular"] = "󱄙"
	IconMap["Other"] = "󰇩"

	IconAsAT := n.props.GetBool("IconAsAT", false)
	ShowType := n.props.GetBool("ShowType", true)
	ShowSSID := n.props.GetBool("ShowSSID", true)
	SSIDAbbr := n.props.GetInt("SSIDAbbr", 0)
	LinkSpeedFull := n.props.GetBool("LinkSpeedFull", false)
	LinkSpeedUnit := Unit(n.props.GetString("LinkSpeedUnit", "Auto"))

	icon := IconMap[connection.Type]
	AT := "@"
	if IconAsAT {
		AT = icon
	} else {
		str += icon
		if !ShowType && !(ShowSSID && connection.Type == "Wi-Fi") {
			AT = ""
		}
	}

	if ShowSSID && connection.Type == "Wi-Fi" {
		if SSIDAbbr > 0 {
			abbr := connection.SSID
			for len(abbr) > SSIDAbbr {
				idx := strings.LastIndexAny(abbr, " #_-")
				if idx == -1 {
					break
				}
				abbr = abbr[:idx]
			}
			str += abbr
		} else {
			str += connection.SSID
		}
	}

	if ShowType && !(ShowSSID && connection.Type == "Wi-Fi") {
		str += connection.Type
	}

	if LinkSpeedUnit != Hide {
		var TransmitLinkSpeed string
		var TransmitLinkSpeedUnit Unit
		var ReceiveLinkSpeed string
		var ReceiveLinkSpeedUnit Unit

		switch LinkSpeedUnit {
		case b, bps:
			TransmitLinkSpeed = fmt.Sprintf("%d", connection.TransmitRate)
			TransmitLinkSpeedUnit = b
			ReceiveLinkSpeed = fmt.Sprintf("%d", connection.ReceiveRate)
			ReceiveLinkSpeedUnit = b
		case K, Kbps:
			TransmitLinkSpeed = strconv.FormatFloat(float64(connection.TransmitRate)/math.Pow10(3), 'f', -1, 64)
			TransmitLinkSpeedUnit = K
			ReceiveLinkSpeed = strconv.FormatFloat(float64(connection.ReceiveRate)/math.Pow10(3), 'f', -1, 64)
			ReceiveLinkSpeedUnit = K
		case M, Mbps:
			TransmitLinkSpeed = strconv.FormatFloat(float64(connection.TransmitRate)/math.Pow10(6), 'f', -1, 64)
			TransmitLinkSpeedUnit = M
			ReceiveLinkSpeed = strconv.FormatFloat(float64(connection.ReceiveRate)/math.Pow10(6), 'f', -1, 64)
			ReceiveLinkSpeedUnit = M
		case G, Gbps:
			TransmitLinkSpeed = strconv.FormatFloat(float64(connection.TransmitRate)/math.Pow10(9), 'f', -1, 64)
			TransmitLinkSpeedUnit = G
			ReceiveLinkSpeed = strconv.FormatFloat(float64(connection.ReceiveRate)/math.Pow10(9), 'f', -1, 64)
			ReceiveLinkSpeedUnit = G
		case T, Tbps:
			TransmitLinkSpeed = strconv.FormatFloat(float64(connection.TransmitRate)/math.Pow10(12), 'f', -1, 64)
			TransmitLinkSpeedUnit = T
			ReceiveLinkSpeed = strconv.FormatFloat(float64(connection.ReceiveRate)/math.Pow10(12), 'f', -1, 64)
			ReceiveLinkSpeedUnit = T
		case A, Auto:
			TransmitSpeedUnitIndex := (len(fmt.Sprintf("%d", connection.TransmitRate)) - 1) / 3
			if TransmitSpeedUnitIndex > 4 {
				TransmitSpeedUnitIndex = 4
			}
			switch TransmitSpeedUnitIndex {
			case 0:
				TransmitLinkSpeedUnit = b
			case 1:
				TransmitLinkSpeedUnit = K
			case 2:
				TransmitLinkSpeedUnit = M
			case 3:
				TransmitLinkSpeedUnit = G
			case 4:
				TransmitLinkSpeedUnit = T
			}
			ReceiveSpeedUnitIndex := (len(fmt.Sprintf("%d", connection.ReceiveRate)) - 1) / 3
			if ReceiveSpeedUnitIndex > 4 {
				ReceiveSpeedUnitIndex = 4
			}
			switch ReceiveSpeedUnitIndex {
			case 0:
				ReceiveLinkSpeedUnit = b
			case 1:
				ReceiveLinkSpeedUnit = K
			case 2:
				ReceiveLinkSpeedUnit = M
			case 3:
				ReceiveLinkSpeedUnit = G
			case 4:
				ReceiveLinkSpeedUnit = T
			}
			if TransmitSpeedUnitIndex == 0 {
				TransmitLinkSpeed = fmt.Sprintf("%d", connection.TransmitRate)
			} else {
				TransmitLinkSpeed = fmt.Sprintf("%.3g", float64(connection.TransmitRate)/math.Pow10(3*TransmitSpeedUnitIndex))
			}
			if ReceiveSpeedUnitIndex == 0 {
				ReceiveLinkSpeed = fmt.Sprintf("%d", connection.ReceiveRate)
			} else {
				ReceiveLinkSpeed = fmt.Sprintf("%.3g", float64(connection.ReceiveRate)/math.Pow10(3*ReceiveSpeedUnitIndex))
			}
		}
		switch LinkSpeedUnit {
		case bps, Kbps, Mbps, Gbps, Tbps, Auto:
			if TransmitLinkSpeedUnit != "" {
				TransmitLinkSpeedUnit += "bps"
			}
			if ReceiveLinkSpeedUnit != "" {
				ReceiveLinkSpeedUnit += "bps"
			}
		}

		if LinkSpeedFull || TransmitLinkSpeedUnit != ReceiveLinkSpeedUnit {
			str += fmt.Sprintf("%s%s%s/%s%s", AT, ReceiveLinkSpeed, ReceiveLinkSpeedUnit, TransmitLinkSpeed, TransmitLinkSpeedUnit)
		} else if TransmitLinkSpeed == ReceiveLinkSpeed {
			str += fmt.Sprintf("%s%s%s", AT, TransmitLinkSpeed, TransmitLinkSpeedUnit)
		} else {
			str += fmt.Sprintf("%s%s/%s%s", AT, ReceiveLinkSpeed, TransmitLinkSpeed, TransmitLinkSpeedUnit)
		}
	}

	if len(str) == 0 {
		return icon
	}

	return str
}
