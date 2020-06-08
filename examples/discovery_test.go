package example

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/beevik/etree"
	"github.com/fanap-infra/onvif"
	"github.com/fanap-infra/onvif/device"
	discover "github.com/fanap-infra/onvif/ws-discovery"
)

func TestGetAvailableDevicesAtSpecificEthernetInterface(t *testing.T) {

	// client()
	// runDiscovery("en0")
	s := onvif.GetAvailableDevicesAtSpecificEthernetInterface("en0")

	log.Printf("%s", s)
}

func client() {
	dev, err := onvif.NewDevice("192.168.3.10")
	if err != nil {
		panic(err)
	}
	dev.Authenticate("admin", "zsyy12345")

	log.Printf("output %+v", dev.GetServices())

	res, err := dev.CallMethod(device.GetUsers{})
	bs, _ := ioutil.ReadAll(res.Body)
	log.Printf("output %+v %s", res.StatusCode, bs)
}

// Host host
type Host struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

func runDiscovery(interfaceName string) {
	var hosts []*Host
	devices := discover.SendProbe(interfaceName, nil, []string{"dn:NetworkVideoTransmitter"}, map[string]string{"dn": "http://www.onvif.org/ver10/network/wsdl"})
	for _, j := range devices {
		doc := etree.NewDocument()
		if err := doc.ReadFromString(j); err != nil {
			log.Printf("error %s", err)
		} else {

			endpoints := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/XAddrs")
			scopes := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/Scopes")

			flag := false

			host := &Host{}

			for _, xaddr := range endpoints {
				xaddr := strings.Split(strings.Split(xaddr.Text(), " ")[0], "/")[2]
				host.URL = xaddr
			}
			if flag {
				break
			}
			for _, scope := range scopes {
				re := regexp.MustCompile(`onvif:\/\/www\.onvif\.org\/name\/[A-Za-z0-9-]+`)
				match := re.FindStringSubmatch(scope.Text())
				host.Name = path.Base(match[0])
			}

			hosts = append(hosts, host)

		}

	}

	bys, _ := json.Marshal(hosts)
	log.Printf("done %s", bys)
}
