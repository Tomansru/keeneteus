package keenetic_api

import (
	"bytes"
	"encoding/json"
	"io"
)

type InterfaceStat struct {
	Devices       []Eth  `json:"-"`
	Interfaces    []Eth  `json:"-"`
	InterfacesStr string `json:"-"`

	Show struct {
		Interface struct {
			Stat []struct {
				InterfaceName      string `json:"-"`
				Rxpackets          int    `json:"rxpackets"`
				RxMulticastPackets int    `json:"rx-multicast-packets"`
				RxBroadcastPackets int    `json:"rx-broadcast-packets"`
				Rxbytes            int64  `json:"rxbytes"`
				Rxerrors           int    `json:"rxerrors"`
				Rxdropped          int    `json:"rxdropped"`
				Txpackets          int    `json:"txpackets"`
				TxMulticastPackets int    `json:"tx-multicast-packets"`
				TxBroadcastPackets int    `json:"tx-broadcast-packets"`
				Txbytes            int64  `json:"txbytes"`
				Txerrors           int    `json:"txerrors"`
				Txdropped          int    `json:"txdropped"`
				Timestamp          string `json:"timestamp"`
				LastOverflow       string `json:"last-overflow"`
				Rxspeed            int    `json:"rxspeed"`
				Txspeed            int    `json:"txspeed"`
			} `json:"stat"`
		} `json:"interface"`
		Ip struct {
			Hotspot struct {
				Chart struct {
					Bar []struct {
						Mac  string `json:"mac,omitempty"`
						Bars []struct {
							Attribute string `json:"attribute"`
							Data      []struct {
								T int `json:"t"`
								V int `json:"v"`
							} `json:"data"`
						} `json:"bars"`
						Multicast bool `json:"multicast,omitempty"`
						Others    bool `json:"others,omitempty"`
					} `json:"bar"`
				} `json:"chart"`
			} `json:"hotspot"`
		} `json:"ip"`
	} `json:"show"`
}

func (i *InterfaceStat) GetRqBody() io.Reader {
	if i.InterfacesStr == "" {
		i.InterfacesStr = `{"show":{"interface":{"stat":[`
		for k := range i.Interfaces {
			i.InterfacesStr += `{"name":"` + i.Interfaces[k].Code + `"},`
		}
		i.InterfacesStr = i.InterfacesStr[:len(i.InterfacesStr)-1]
		i.InterfacesStr += `]},"ip":{"hotspot":{"chart":{"items":"`
		for k := range i.Devices {
			i.InterfacesStr += i.Devices[k].Code + `,`
		}
		i.InterfacesStr = i.InterfacesStr[:len(i.InterfacesStr)-1]
		i.InterfacesStr += `","detail":0,"attributes":"rxbytes,txbytes"}}}}}`
	}
	return bytes.NewBufferString(i.InterfacesStr)
}

func (i *InterfaceStat) Unmarshal(b io.Reader) error {
	if err := json.NewDecoder(b).Decode(i); err != nil {
		return err
	}
	for k := range i.Show.Interface.Stat {
		i.Show.Interface.Stat[k].InterfaceName = i.GetInterfaceName(k)
	}
	return nil
}

type Eth struct {
	Name string
	Code string
}

// SetInterfaces Set interfaces for monitoring
func (i *InterfaceStat) SetInterfaces(interfaces []Eth) {
	i.Interfaces = interfaces
}

func (i *InterfaceStat) GetInterfaceName(k int) string {
	if len(i.Interfaces) > k {
		return i.Interfaces[k].Name
	}
	return ""
}

// SetDevices Set interfaces for monitoring
func (i *InterfaceStat) SetDevices(interfaces []Eth) {
	i.Devices = interfaces
}

func (i *InterfaceStat) GetDeviceName(k int) string {
	if len(i.Devices) > k {
		return i.Devices[k].Name
	}
	return ""
}

type Metrics struct {
	Whoami struct {
		User  string `json:"user"`
		Agent string `json:"agent"`
		Host  string `json:"host"`
		Mac   string `json:"mac"`
		Where string `json:"where"`
	} `json:"whoami"`
	Show struct {
		Version struct {
			Release string `json:"release"`
			Sandbox string `json:"sandbox"`
			Title   string `json:"title"`
			Arch    string `json:"arch"`
			Ndm     struct {
				Exact string `json:"exact"`
				Cdate string `json:"cdate"`
			} `json:"ndm"`
			Bsp struct {
				Exact string `json:"exact"`
				Cdate string `json:"cdate"`
			} `json:"bsp"`
			Ndw struct {
				Version    string `json:"version"`
				Features   string `json:"features"`
				Components string `json:"components"`
			} `json:"ndw"`
			Manufacturer string `json:"manufacturer"`
			Vendor       string `json:"vendor"`
			Series       string `json:"series"`
			Model        string `json:"model"`
			HwVersion    string `json:"hw_version"`
			HwId         string `json:"hw_id"`
			Device       string `json:"device"`
			Region       string `json:"region"`
			Description  string `json:"description"`
		} `json:"version"`
		System struct {
			Hostname   string `json:"hostname"`
			Domainname string `json:"domainname"`
			Cpuload    int    `json:"cpuload"`
			Memory     string `json:"memory"`
			Swap       string `json:"swap"`
			Memtotal   int    `json:"memtotal"`
			Memfree    int    `json:"memfree"`
			Membuffers int    `json:"membuffers"`
			Memcache   int    `json:"memcache"`
			Swaptotal  int    `json:"swaptotal"`
			Swapfree   int    `json:"swapfree"`
			Uptime     string `json:"uptime"`
		} `json:"system"`
		Media struct {
			Media0 struct {
				Usb struct {
					Port    int    `json:"port"`
					Version string `json:"version"`
				} `json:"usb"`
				State        string `json:"state"`
				Manufacturer string `json:"manufacturer"`
				Product      string `json:"product"`
				Serial       string `json:"serial"`
				Size         string `json:"size"`
				Partition    []struct {
					Uuid   string `json:"uuid"`
					Label  string `json:"label"`
					Fstype string `json:"fstype"`
					State  string `json:"state"`
					Total  string `json:"total"`
					Free   string `json:"free"`
				} `json:"partition"`
			} `json:"Media0"`
		} `json:"media"`
		Interface struct {
			GigabitEthernet0 struct {
				Id            string `json:"id"`
				Index         int    `json:"index"`
				Type          string `json:"type"`
				Description   string `json:"description"`
				InterfaceName string `json:"interface-name"`
				Link          string `json:"link"`
				Connected     string `json:"connected"`
				State         string `json:"state"`
				Mtu           int    `json:"mtu"`
				TxQueueLength int    `json:"tx-queue-length"`
				Port          struct {
					Field1 struct {
						Id            string `json:"id"`
						Index         int    `json:"index"`
						InterfaceName string `json:"interface-name"`
						Type          string `json:"type"`
						Link          string `json:"link"`
						Role          []struct {
							For  string `json:"for"`
							Role string `json:"role"`
						} `json:"role"`
						Speed           string `json:"speed"`
						Duplex          string `json:"duplex"`
						AutoNegotiation string `json:"auto-negotiation"`
						FlowControl     string `json:"flow-control"`
						Eee             string `json:"eee"`
						LastChange      string `json:"last-change"`
						LastOverflow    string `json:"last-overflow"`
						Public          bool   `json:"public"`
						LinkGroup       struct {
							Supported bool `json:"supported"`
						} `json:"link-group"`
					} `json:"1"`
					Field2 struct {
						Id              string `json:"id"`
						Index           int    `json:"index"`
						InterfaceName   string `json:"interface-name"`
						Type            string `json:"type"`
						Link            string `json:"link"`
						Speed           string `json:"speed"`
						Duplex          string `json:"duplex"`
						AutoNegotiation string `json:"auto-negotiation"`
						FlowControl     string `json:"flow-control"`
						Eee             string `json:"eee"`
						LastChange      string `json:"last-change"`
						LastOverflow    string `json:"last-overflow"`
						Public          bool   `json:"public"`
						LinkGroup       struct {
							Supported bool `json:"supported"`
						} `json:"link-group"`
					} `json:"2"`
					Field3 struct {
						Id            string `json:"id"`
						Index         int    `json:"index"`
						InterfaceName string `json:"interface-name"`
						Type          string `json:"type"`
						Link          string `json:"link"`
						LastChange    string `json:"last-change"`
						LastOverflow  string `json:"last-overflow"`
						Public        bool   `json:"public"`
						LinkGroup     struct {
							Supported bool `json:"supported"`
						} `json:"link-group"`
					} `json:"3"`
					Field4 struct {
						Id              string `json:"id"`
						Index           int    `json:"index"`
						InterfaceName   string `json:"interface-name"`
						Type            string `json:"type"`
						Link            string `json:"link"`
						Speed           string `json:"speed"`
						Duplex          string `json:"duplex"`
						AutoNegotiation string `json:"auto-negotiation"`
						FlowControl     string `json:"flow-control"`
						Eee             string `json:"eee"`
						LastChange      string `json:"last-change"`
						LastOverflow    string `json:"last-overflow"`
						Public          bool   `json:"public"`
						LinkGroup       struct {
							Supported bool `json:"supported"`
						} `json:"link-group"`
					} `json:"4"`
				} `json:"port"`
			} `json:"GigabitEthernet0"`
			Field2 struct {
				Id            string `json:"id"`
				Index         int    `json:"index"`
				InterfaceName string `json:"interface-name"`
				Type          string `json:"type"`
				Link          string `json:"link"`
				Role          []struct {
					For  string `json:"for"`
					Role string `json:"role"`
				} `json:"role"`
				Speed           string `json:"speed"`
				Duplex          string `json:"duplex"`
				AutoNegotiation string `json:"auto-negotiation"`
				FlowControl     string `json:"flow-control"`
				Eee             string `json:"eee"`
				LastChange      string `json:"last-change"`
				LastOverflow    string `json:"last-overflow"`
				Public          bool   `json:"public"`
				LinkGroup       struct {
					Supported bool `json:"supported"`
				} `json:"link-group"`
			} `json:"1"`
			Field3 struct {
				Id              string `json:"id"`
				Index           int    `json:"index"`
				InterfaceName   string `json:"interface-name"`
				Type            string `json:"type"`
				Link            string `json:"link"`
				Speed           string `json:"speed"`
				Duplex          string `json:"duplex"`
				AutoNegotiation string `json:"auto-negotiation"`
				FlowControl     string `json:"flow-control"`
				Eee             string `json:"eee"`
				LastChange      string `json:"last-change"`
				LastOverflow    string `json:"last-overflow"`
				Public          bool   `json:"public"`
				LinkGroup       struct {
					Supported bool `json:"supported"`
				} `json:"link-group"`
			} `json:"2"`
			Field4 struct {
				Id            string `json:"id"`
				Index         int    `json:"index"`
				InterfaceName string `json:"interface-name"`
				Type          string `json:"type"`
				Link          string `json:"link"`
				LastChange    string `json:"last-change"`
				LastOverflow  string `json:"last-overflow"`
				Public        bool   `json:"public"`
				LinkGroup     struct {
					Supported bool `json:"supported"`
				} `json:"link-group"`
			} `json:"3"`
			Field5 struct {
				Id              string `json:"id"`
				Index           int    `json:"index"`
				InterfaceName   string `json:"interface-name"`
				Type            string `json:"type"`
				Link            string `json:"link"`
				Speed           string `json:"speed"`
				Duplex          string `json:"duplex"`
				AutoNegotiation string `json:"auto-negotiation"`
				FlowControl     string `json:"flow-control"`
				Eee             string `json:"eee"`
				LastChange      string `json:"last-change"`
				LastOverflow    string `json:"last-overflow"`
				Public          bool   `json:"public"`
				LinkGroup       struct {
					Supported bool `json:"supported"`
				} `json:"link-group"`
			} `json:"4"`
			GigabitEthernet0Vlan1 struct {
				Id            string   `json:"id"`
				Index         int      `json:"index"`
				Type          string   `json:"type"`
				Description   string   `json:"description"`
				InterfaceName string   `json:"interface-name"`
				Link          string   `json:"link"`
				Connected     string   `json:"connected"`
				State         string   `json:"state"`
				Mtu           int      `json:"mtu"`
				TxQueueLength int      `json:"tx-queue-length"`
				Group         string   `json:"group"`
				Usedby        []string `json:"usedby"`
				Mac           string   `json:"mac"`
				AuthType      string   `json:"auth-type"`
			} `json:"GigabitEthernet0/Vlan1"`
			GigabitEthernet0Vlan2 struct {
				Id            string   `json:"id"`
				Index         int      `json:"index"`
				Type          string   `json:"type"`
				Description   string   `json:"description"`
				InterfaceName string   `json:"interface-name"`
				Link          string   `json:"link"`
				Connected     string   `json:"connected"`
				State         string   `json:"state"`
				Mtu           int      `json:"mtu"`
				TxQueueLength int      `json:"tx-queue-length"`
				Group         string   `json:"group"`
				Usedby        []string `json:"usedby"`
				Mac           string   `json:"mac"`
				AuthType      string   `json:"auth-type"`
			} `json:"GigabitEthernet0/Vlan2"`
			GigabitEthernet0Vlan4 struct {
				Id            string `json:"id"`
				Index         int    `json:"index"`
				Type          string `json:"type"`
				Description   string `json:"description"`
				InterfaceName string `json:"interface-name"`
				Link          string `json:"link"`
				Connected     string `json:"connected"`
				State         string `json:"state"`
				Mtu           int    `json:"mtu"`
				TxQueueLength int    `json:"tx-queue-length"`
				Address       string `json:"address"`
				Mask          string `json:"mask"`
				Uptime        int    `json:"uptime"`
				Global        bool   `json:"global"`
				Defaultgw     bool   `json:"defaultgw"`
				Priority      int    `json:"priority"`
				SecurityLevel string `json:"security-level"`
				Mac           string `json:"mac"`
				AuthType      string `json:"auth-type"`
			} `json:"GigabitEthernet0/Vlan4"`
			ISP struct {
				Id            string `json:"id"`
				Index         int    `json:"index"`
				Type          string `json:"type"`
				Description   string `json:"description"`
				InterfaceName string `json:"interface-name"`
				Link          string `json:"link"`
				Connected     string `json:"connected"`
				State         string `json:"state"`
				Mtu           int    `json:"mtu"`
				TxQueueLength int    `json:"tx-queue-length"`
				Address       string `json:"address"`
				Mask          string `json:"mask"`
				Uptime        int    `json:"uptime"`
				Global        bool   `json:"global"`
				Defaultgw     bool   `json:"defaultgw"`
				Priority      int    `json:"priority"`
				SecurityLevel string `json:"security-level"`
				Mac           string `json:"mac"`
				AuthType      string `json:"auth-type"`
				Port          struct {
					Id              string `json:"id"`
					Index           int    `json:"index"`
					InterfaceName   string `json:"interface-name"`
					Type            string `json:"type"`
					Link            string `json:"link"`
					AutoNegotiation string `json:"auto-negotiation"`
					Speed           string `json:"speed"`
					Duplex          string `json:"duplex"`
					FlowControl     string `json:"flow-control"`
					Transceiver     string `json:"transceiver"`
					SfpCombo        bool   `json:"sfp-combo"`
				} `json:"port"`
			} `json:"ISP"`
			Field10 struct {
				Id              string `json:"id"`
				Index           int    `json:"index"`
				InterfaceName   string `json:"interface-name"`
				Type            string `json:"type"`
				Link            string `json:"link"`
				AutoNegotiation string `json:"auto-negotiation"`
				Speed           string `json:"speed"`
				Duplex          string `json:"duplex"`
				FlowControl     string `json:"flow-control"`
				Transceiver     string `json:"transceiver"`
				SfpCombo        bool   `json:"sfp-combo"`
			} `json:"0"`
			WifiMaster0 struct {
				Id            string `json:"id"`
				Index         int    `json:"index"`
				Type          string `json:"type"`
				Description   string `json:"description"`
				InterfaceName string `json:"interface-name"`
				Link          string `json:"link"`
				Connected     string `json:"connected"`
				State         string `json:"state"`
				Mtu           int    `json:"mtu"`
				TxQueueLength int    `json:"tx-queue-length"`
				Hwstate       string `json:"hwstate"`
				Bitrate       int    `json:"bitrate"`
				Channel       int    `json:"channel"`
				Temperature   int    `json:"temperature"`
			} `json:"WifiMaster0"`
			AccessPoint struct {
				Id            string   `json:"id"`
				Index         int      `json:"index"`
				Type          string   `json:"type"`
				Description   string   `json:"description"`
				InterfaceName string   `json:"interface-name"`
				Link          string   `json:"link"`
				Connected     string   `json:"connected"`
				State         string   `json:"state"`
				Mtu           int      `json:"mtu"`
				TxQueueLength int      `json:"tx-queue-length"`
				Group         string   `json:"group"`
				Usedby        []string `json:"usedby"`
				Mac           string   `json:"mac"`
				AuthType      string   `json:"auth-type"`
				Ssid          string   `json:"ssid"`
				Encryption    string   `json:"encryption"`
			} `json:"AccessPoint"`
			GuestWiFi struct {
				Id            string   `json:"id"`
				Index         int      `json:"index"`
				Type          string   `json:"type"`
				Description   string   `json:"description"`
				InterfaceName string   `json:"interface-name"`
				Link          string   `json:"link"`
				Connected     string   `json:"connected"`
				State         string   `json:"state"`
				Mtu           int      `json:"mtu"`
				TxQueueLength int      `json:"tx-queue-length"`
				Group         string   `json:"group"`
				Usedby        []string `json:"usedby"`
				Mac           string   `json:"mac"`
				AuthType      string   `json:"auth-type"`
				Encryption    string   `json:"encryption"`
			} `json:"GuestWiFi"`
			WifiMaster0AccessPoint2 struct {
				Id            string `json:"id"`
				Index         int    `json:"index"`
				Type          string `json:"type"`
				Description   string `json:"description"`
				InterfaceName string `json:"interface-name"`
				Link          string `json:"link"`
				Connected     string `json:"connected"`
				State         string `json:"state"`
				Mtu           int    `json:"mtu"`
				TxQueueLength int    `json:"tx-queue-length"`
				Global        bool   `json:"global"`
				SecurityLevel string `json:"security-level"`
				Mac           string `json:"mac"`
				AuthType      string `json:"auth-type"`
				Encryption    string `json:"encryption"`
			} `json:"WifiMaster0/AccessPoint2"`
			WifiMaster0AccessPoint3 struct {
				Id            string `json:"id"`
				Index         int    `json:"index"`
				Type          string `json:"type"`
				Description   string `json:"description"`
				InterfaceName string `json:"interface-name"`
				Link          string `json:"link"`
				Connected     string `json:"connected"`
				State         string `json:"state"`
				Mtu           int    `json:"mtu"`
				TxQueueLength int    `json:"tx-queue-length"`
				Global        bool   `json:"global"`
				SecurityLevel string `json:"security-level"`
				Mac           string `json:"mac"`
				AuthType      string `json:"auth-type"`
				Encryption    string `json:"encryption"`
			} `json:"WifiMaster0/AccessPoint3"`
			WifiMaster0WifiStation0 struct {
				Ap            string `json:"ap"`
				Id            string `json:"id"`
				Index         int    `json:"index"`
				Type          string `json:"type"`
				Description   string `json:"description"`
				InterfaceName string `json:"interface-name"`
				Link          string `json:"link"`
				Connected     string `json:"connected"`
				State         string `json:"state"`
				Mtu           int    `json:"mtu"`
				TxQueueLength int    `json:"tx-queue-length"`
				Global        bool   `json:"global"`
				SecurityLevel string `json:"security-level"`
				Mac           string `json:"mac"`
				AuthType      string `json:"auth-type"`
				Encryption    string `json:"encryption"`
			} `json:"WifiMaster0/WifiStation0"`
			WifiMaster1 struct {
				Id            string `json:"id"`
				Index         int    `json:"index"`
				Type          string `json:"type"`
				Description   string `json:"description"`
				InterfaceName string `json:"interface-name"`
				Link          string `json:"link"`
				Connected     string `json:"connected"`
				State         string `json:"state"`
				Mtu           int    `json:"mtu"`
				TxQueueLength int    `json:"tx-queue-length"`
				Hwstate       string `json:"hwstate"`
				Bitrate       int    `json:"bitrate"`
				Channel       int    `json:"channel"`
				Temperature   int    `json:"temperature"`
			} `json:"WifiMaster1"`
			AccessPoint5G struct {
				Id            string   `json:"id"`
				Index         int      `json:"index"`
				Type          string   `json:"type"`
				Description   string   `json:"description"`
				InterfaceName string   `json:"interface-name"`
				Link          string   `json:"link"`
				Connected     string   `json:"connected"`
				State         string   `json:"state"`
				Mtu           int      `json:"mtu"`
				TxQueueLength int      `json:"tx-queue-length"`
				Group         string   `json:"group"`
				Usedby        []string `json:"usedby"`
				Mac           string   `json:"mac"`
				AuthType      string   `json:"auth-type"`
				Ssid          string   `json:"ssid"`
				Encryption    string   `json:"encryption"`
			} `json:"AccessPoint_5G"`
			GuestWiFi5G struct {
				Id            string   `json:"id"`
				Index         int      `json:"index"`
				Type          string   `json:"type"`
				Description   string   `json:"description"`
				InterfaceName string   `json:"interface-name"`
				Link          string   `json:"link"`
				Connected     string   `json:"connected"`
				State         string   `json:"state"`
				Mtu           int      `json:"mtu"`
				TxQueueLength int      `json:"tx-queue-length"`
				Group         string   `json:"group"`
				Usedby        []string `json:"usedby"`
				Mac           string   `json:"mac"`
				AuthType      string   `json:"auth-type"`
				Encryption    string   `json:"encryption"`
			} `json:"GuestWiFi_5G"`
			WifiMaster1AccessPoint2 struct {
				Id            string `json:"id"`
				Index         int    `json:"index"`
				Type          string `json:"type"`
				Description   string `json:"description"`
				InterfaceName string `json:"interface-name"`
				Link          string `json:"link"`
				Connected     string `json:"connected"`
				State         string `json:"state"`
				Mtu           int    `json:"mtu"`
				TxQueueLength int    `json:"tx-queue-length"`
				Global        bool   `json:"global"`
				SecurityLevel string `json:"security-level"`
				Mac           string `json:"mac"`
				AuthType      string `json:"auth-type"`
				Encryption    string `json:"encryption"`
			} `json:"WifiMaster1/AccessPoint2"`
			WifiMaster1AccessPoint3 struct {
				Id            string `json:"id"`
				Index         int    `json:"index"`
				Type          string `json:"type"`
				Description   string `json:"description"`
				InterfaceName string `json:"interface-name"`
				Link          string `json:"link"`
				Connected     string `json:"connected"`
				State         string `json:"state"`
				Mtu           int    `json:"mtu"`
				TxQueueLength int    `json:"tx-queue-length"`
				Global        bool   `json:"global"`
				SecurityLevel string `json:"security-level"`
				Mac           string `json:"mac"`
				AuthType      string `json:"auth-type"`
				Encryption    string `json:"encryption"`
			} `json:"WifiMaster1/AccessPoint3"`
			WifiMaster1WifiStation0 struct {
				Ap            string `json:"ap"`
				Id            string `json:"id"`
				Index         int    `json:"index"`
				Type          string `json:"type"`
				Description   string `json:"description"`
				InterfaceName string `json:"interface-name"`
				Link          string `json:"link"`
				Connected     string `json:"connected"`
				State         string `json:"state"`
				Mtu           int    `json:"mtu"`
				TxQueueLength int    `json:"tx-queue-length"`
				Global        bool   `json:"global"`
				SecurityLevel string `json:"security-level"`
				Mac           string `json:"mac"`
				AuthType      string `json:"auth-type"`
				Encryption    string `json:"encryption"`
			} `json:"WifiMaster1/WifiStation0"`
			Home struct {
				Id            string `json:"id"`
				Index         int    `json:"index"`
				Type          string `json:"type"`
				Description   string `json:"description"`
				InterfaceName string `json:"interface-name"`
				Link          string `json:"link"`
				Connected     string `json:"connected"`
				State         string `json:"state"`
				Mtu           int    `json:"mtu"`
				TxQueueLength int    `json:"tx-queue-length"`
				Address       string `json:"address"`
				Mask          string `json:"mask"`
				Uptime        int    `json:"uptime"`
				Global        bool   `json:"global"`
				SecurityLevel string `json:"security-level"`
				Mac           string `json:"mac"`
				AuthType      string `json:"auth-type"`
				Bridge        struct {
					Interface []struct {
						Link      bool   `json:"link"`
						Inherited string `json:"inherited,omitempty"`
						Interface string `json:"interface"`
					} `json:"interface"`
				} `json:"bridge"`
			} `json:"Home"`
			Guest struct {
				Id            string `json:"id"`
				Index         int    `json:"index"`
				Type          string `json:"type"`
				Description   string `json:"description"`
				InterfaceName string `json:"interface-name"`
				Link          string `json:"link"`
				Connected     string `json:"connected"`
				State         string `json:"state"`
				Mtu           int    `json:"mtu"`
				TxQueueLength int    `json:"tx-queue-length"`
				Address       string `json:"address"`
				Mask          string `json:"mask"`
				Uptime        int    `json:"uptime"`
				Global        bool   `json:"global"`
				SecurityLevel string `json:"security-level"`
				Mac           string `json:"mac"`
				AuthType      string `json:"auth-type"`
				Bridge        struct {
					Interface []struct {
						Link      bool   `json:"link"`
						Interface string `json:"interface"`
					} `json:"interface"`
				} `json:"bridge"`
			} `json:"Guest"`
			OpenVPN0 struct {
				Id            string   `json:"id"`
				Index         int      `json:"index"`
				Type          string   `json:"type"`
				Description   string   `json:"description"`
				InterfaceName string   `json:"interface-name"`
				Link          string   `json:"link"`
				Connected     string   `json:"connected"`
				State         string   `json:"state"`
				Role          []string `json:"role"`
				Mtu           int      `json:"mtu"`
				TxQueueLength int      `json:"tx-queue-length"`
				Global        bool     `json:"global"`
				Defaultgw     bool     `json:"defaultgw"`
				Priority      int      `json:"priority"`
				SecurityLevel string   `json:"security-level"`
				Mac           string   `json:"mac"`
				AuthType      string   `json:"auth-type"`
				Via           string   `json:"via"`
			} `json:"OpenVPN0"`
			Wireguard0 struct {
				Id            string `json:"id"`
				Index         int    `json:"index"`
				Type          string `json:"type"`
				Description   string `json:"description"`
				InterfaceName string `json:"interface-name"`
				Link          string `json:"link"`
				Connected     string `json:"connected"`
				State         string `json:"state"`
				Mtu           int    `json:"mtu"`
				TxQueueLength int    `json:"tx-queue-length"`
				Address       string `json:"address"`
				Mask          string `json:"mask"`
				Uptime        int    `json:"uptime"`
				Global        bool   `json:"global"`
				Defaultgw     bool   `json:"defaultgw"`
				Priority      int    `json:"priority"`
				SecurityLevel string `json:"security-level"`
				Wireguard     struct {
					PublicKey  string `json:"public-key"`
					ListenPort int    `json:"listen-port"`
					Status     string `json:"status"`
					Peer       []struct {
						PublicKey     string `json:"public-key"`
						Local         string `json:"local"`
						LocalPort     int    `json:"local-port"`
						Via           string `json:"via"`
						Remote        string `json:"remote"`
						RemotePort    int    `json:"remote-port"`
						Rxbytes       int64  `json:"rxbytes"`
						Txbytes       int    `json:"txbytes"`
						LastHandshake int    `json:"last-handshake"`
						Online        bool   `json:"online"`
					} `json:"peer"`
				} `json:"wireguard"`
			} `json:"Wireguard0"`
		} `json:"interface"`
		Ip struct {
			NameServer struct {
				Server []struct {
					Address   string `json:"address"`
					Port      string `json:"port"`
					Domain    string `json:"domain"`
					Global    int    `json:"global"`
					Service   string `json:"service"`
					Interface string `json:"interface"`
				} `json:"server"`
			} `json:"name-server"`
			Hotspot struct {
				Host []struct {
					Mac       string `json:"mac"`
					Via       string `json:"via"`
					Ip        string `json:"ip"`
					Hostname  string `json:"hostname"`
					Name      string `json:"name"`
					Interface struct {
						Id          string `json:"id"`
						Name        string `json:"name"`
						Description string `json:"description"`
					} `json:"interface,omitempty"`
					Registered    bool     `json:"registered"`
					Access        string   `json:"access"`
					Schedule      string   `json:"schedule"`
					Active        bool     `json:"active"`
					Rxbytes       int      `json:"rxbytes"`
					Txbytes       int      `json:"txbytes"`
					FirstSeen     int      `json:"first-seen,omitempty"`
					LastSeen      int      `json:"last-seen,omitempty"`
					Link          string   `json:"link,omitempty"`
					Ssid          string   `json:"ssid,omitempty"`
					Ap            string   `json:"ap,omitempty"`
					Authenticated bool     `json:"authenticated,omitempty"`
					Txrate        int      `json:"txrate,omitempty"`
					Uptime        int      `json:"uptime"`
					Ht            int      `json:"ht,omitempty"`
					Mode          string   `json:"mode,omitempty"`
					Gi            int      `json:"gi,omitempty"`
					Rssi          int      `json:"rssi,omitempty"`
					Mcs           int      `json:"mcs,omitempty"`
					Txss          int      `json:"txss,omitempty"`
					Ebf           bool     `json:"ebf,omitempty"`
					DlMu          bool     `json:"dl-mu,omitempty"`
					Field29       []string `json:"_11,omitempty"`
					Security      string   `json:"security,omitempty"`
					TrafficShape  struct {
						Rx       int    `json:"rx"`
						Tx       int    `json:"tx"`
						Mode     string `json:"mode"`
						Schedule string `json:"schedule"`
					} `json:"traffic-shape"`
					Roam string `json:"roam,omitempty"`
					Dhcp struct {
						Expires int `json:"expires"`
					} `json:"dhcp,omitempty"`
				} `json:"host"`
			} `json:"hotspot"`
		} `json:"ip"`
		Acme struct {
			ServerEnabled    bool   `json:"server-enabled"`
			RealTime         bool   `json:"real-time"`
			NdnsDomain       string `json:"ndns-domain"`
			NdnsDomainAcme   bool   `json:"ndns-domain-acme"`
			NdnsDomainError  bool   `json:"ndns-domain-error"`
			DefaultDomain    string `json:"default-domain"`
			AccountPending   bool   `json:"account-pending"`
			AccountRunning   bool   `json:"account-running"`
			GetPending       bool   `json:"get-pending"`
			GetRunning       bool   `json:"get-running"`
			RevokePending    bool   `json:"revoke-pending"`
			RevokeRunning    bool   `json:"revoke-running"`
			ReissueQueueSize int    `json:"reissue-queue-size"`
			RevokeQueueSize  int    `json:"revoke-queue-size"`
			Retries          int    `json:"retries"`
			CheckerTimer     int    `json:"checker-timer"`
			ApplyTimer       int    `json:"apply-timer"`
			AcmeAccount      string `json:"acme-account"`
			NextTryTa        int    `json:"next-try-ta"`
			Jitter           int    `json:"jitter"`
		} `json:"acme"`
		Cifs struct {
			Enabled    bool `json:"enabled"`
			Automount  bool `json:"automount"`
			Permissive bool `json:"permissive"`
			MapHidden  bool `json:"map-hidden"`
			Share      []struct {
				Mount       string `json:"mount"`
				Label       string `json:"label"`
				Timemachine bool   `json:"timemachine"`
				Description string `json:"description"`
				Active      bool   `json:"active"`
			} `json:"share"`
		} `json:"cifs"`
		Dlna struct {
			Running   bool `json:"running"`
			Directory struct {
				Fbd65Edf34004871970C9D010D087Download struct {
					MediaType string `json:"media-type"`
					Mounted   bool   `json:"mounted"`
					Found     bool   `json:"found"`
				} `json:"443fbd65-edf3-4004-8719-70c9d010d087:/download"`
			} `json:"directory"`
			Db struct {
				Name      string `json:"name"`
				MediaType string `json:"media-type"`
				Mounted   bool   `json:"mounted"`
				Found     bool   `json:"found"`
			} `json:"db"`
		} `json:"dlna"`
		Torrent struct {
			Status struct {
				State   string `json:"state"`
				RpcPort int    `json:"rpc-port"`
			} `json:"status"`
		} `json:"torrent"`
		Ndns struct {
			Name     string `json:"name"`
			Booked   string `json:"booked"`
			Domain   string `json:"domain"`
			Address  string `json:"address"`
			Address6 string `json:"address6"`
			Updated  bool   `json:"updated"`
			Access   string `json:"access"`
			Access6  string `json:"access6"`
			Xns      string `json:"xns"`
			Ttp      struct {
				Direct    bool   `json:"direct"`
				Interface string `json:"interface"`
				Address   string `json:"address"`
			} `json:"ttp"`
		} `json:"ndns"`
		Internet struct {
			Status struct {
				Checked           string `json:"checked"`
				Enabled           bool   `json:"enabled"`
				Reliable          bool   `json:"reliable"`
				GatewayAccessible bool   `json:"gateway-accessible"`
				DnsAccessible     bool   `json:"dns-accessible"`
				HostAccessible    bool   `json:"host-accessible"`
				CaptiveAccessible bool   `json:"captive-accessible"`
				Internet          bool   `json:"internet"`
				Gateway           struct {
					Interface  string `json:"interface"`
					Address    string `json:"address"`
					Failures   int    `json:"failures"`
					Accessible bool   `json:"accessible"`
					Excluded   bool   `json:"excluded"`
				} `json:"gateway"`
				Captive struct {
					Response string `json:"response"`
					Location string `json:"location"`
					Failures int    `json:"failures"`
					Resolved bool   `json:"resolved"`
				} `json:"captive"`
				Hosts struct {
					YaRu struct {
						Failures   int    `json:"failures"`
						Resolved   bool   `json:"resolved"`
						Accessible bool   `json:"accessible"`
						Response   string `json:"response"`
					} `json:"ya.ru"`
					NicRu struct {
						Failures   int    `json:"failures"`
						Resolved   bool   `json:"resolved"`
						Accessible bool   `json:"accessible"`
						Response   string `json:"response"`
					} `json:"nic.ru"`
					GoogleCom struct {
						Failures   int    `json:"failures"`
						Resolved   bool   `json:"resolved"`
						Accessible bool   `json:"accessible"`
						Response   string `json:"response"`
					} `json:"google.com"`
				} `json:"hosts"`
			} `json:"status"`
		} `json:"internet"`
		PingCheck struct {
			Pingcheck []struct {
				Profile   string `json:"profile"`
				Interface struct {
					GigabitEthernet0Vlan4 struct {
						Successcount int    `json:"successcount"`
						Failcount    int    `json:"failcount"`
						Status       string `json:"status"`
						Ipcache      []struct {
							Host      string   `json:"host"`
							Addresses []string `json:"addresses"`
						} `json:"ipcache"`
					} `json:"GigabitEthernet0/Vlan4"`
					ISP struct {
						Successcount int    `json:"successcount"`
						Failcount    int    `json:"failcount"`
						Status       string `json:"status"`
						Ipcache      []struct {
							Host      string   `json:"host"`
							Addresses []string `json:"addresses"`
						} `json:"ipcache"`
					} `json:"ISP"`
				} `json:"interface"`
			} `json:"pingcheck"`
		} `json:"ping-check"`
		Clock struct {
			Date struct {
				Weekday int    `json:"weekday"`
				Day     int    `json:"day"`
				Month   int    `json:"month"`
				Year    int    `json:"year"`
				Hour    int    `json:"hour"`
				Min     int    `json:"min"`
				Sec     int    `json:"sec"`
				Msec    int    `json:"msec"`
				Dst     string `json:"dst"`
				Tz      []struct {
					Locality  string `json:"locality"`
					Stdoffset int    `json:"stdoffset"`
					Dstoffset int    `json:"dstoffset"`
					Usesdst   bool   `json:"usesdst"`
					Rule      string `json:"rule"`
					Custom    bool   `json:"custom"`
				} `json:"tz"`
			} `json:"date"`
		} `json:"clock"`
		Usb struct {
			Device struct {
				Media0 struct {
					DEVICE       string `json:"DEVICE"`
					DEVPATH      string `json:"DEVPATH"`
					Manufacturer string `json:"manufacturer"`
					Product      string `json:"product"`
					Serial       string `json:"serial"`
					Subsystem    string `json:"subsystem"`
					Port         string `json:"port"`
					PowerControl string `json:"power-control"`
					UsbVersion   string `json:"usb-version"`
				} `json:"Media0"`
			} `json:"device"`
		} `json:"usb"`
	} `json:"show"`
}

func (i *Metrics) GetRqBody() io.Reader {
	return bytes.NewBufferString(`{"show":{"clock":{"date":{}},"internet":{"status":{}},"version":{},"system":{},"interface":{},"ip":{"name-server":{},"hotspot":{"details":"wireless"}},"ndns":{},"acme":{},"ping-check":{},"cifs":{},"dlna":{},"torrent":{"status":{}},"usb":{},"media":{}},"whoami":{}}`)
}

func (i *Metrics) Unmarshal(b io.Reader) error {
	return json.NewDecoder(b).Decode(i)
}
