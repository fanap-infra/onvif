package wsdiscovery

/*******************************************************
 * Copyright (C) 2018 Palanjyan Zhorzhik
 *
 * This file is part of ws-discovery project.
 *
 * ws-discovery can be copied and/or distributed without the express
 * permission of Palanjyan Zhorzhik
 *******************************************************/

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/gofrs/uuid"
	"golang.org/x/net/ipv4"
)

const bufSize = 8192

//SendProbeToSpecificDevice ...
func SendProbeToSpecificDevice(ip string, scopes, types []string, namespaces map[string]string) string {
	uuidV4 := uuid.Must(uuid.NewV4())
	probeSOAP := buildProbeMessage(uuidV4.String(), scopes, types, namespaces)
	response := writeUDP(ip, 3702, probeSOAP.String())
	return response
}

func writeUDP(ip string, port int, data string) string {

	address := net.UDPAddr{Port: port, IP: net.ParseIP(ip)}

	fmt.Printf("Client to contact server at %v\n", address)

	conn, err := net.DialUDP("udp", nil, &address)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Connected: %T, %v\n", conn, conn)

	fmt.Printf("Local address: %v\n", conn.LocalAddr())
	fmt.Printf("Remote address: %v\n", conn.RemoteAddr())

	b := []byte(data)
	readBytes := make([]byte, 4000)

	_, wrerr := conn.Write(b)

	if wrerr != nil {
		fmt.Printf("conn.Write() error: %s\n", wrerr)
		return ""
	}

	conn.SetDeadline(time.Now().Add(5 * time.Second))
	cc, rderr := conn.Read(readBytes)
	if rderr != nil {
		fmt.Printf("conn.Read() error: %s\n", rderr)
		return ""
	}

	if err = conn.Close(); err != nil {
		return ""
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromString(string(readBytes[0:cc])); err != nil {
		log.Printf("error %s", err)
		return ""
	}

	uuid := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/EndpointReference/Address")
	uuidStr := strings.Split(uuid[0].Text(), ":")[2]

	return uuidStr
}

//SendProbe to device
func SendProbe(interfaceName string, scopes, types []string, namespaces map[string]string) []string {
	// Creating UUID Version 4
	uuidV4 := uuid.Must(uuid.NewV4())
	// fmt.Printf("UUIDv4: %s\n", uuidV4)

	probeSOAP := buildProbeMessage(uuidV4.String(), scopes, types, namespaces)
	//probeSOAP = `<?xml version="1.0" encoding="UTF-8"?>
	//<Envelope xmlns="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing">
	//<Header>
	//<a:Action mustUnderstand="1">http://schemas.xmlsoap.org/ws/2005/04/discovery/Probe</a:Action>
	//<a:MessageID>uuid:78a2ed98-bc1f-4b08-9668-094fcba81e35</a:MessageID><a:ReplyTo>
	//<a:Address>http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:Address>
	//</a:ReplyTo><a:To mustUnderstand="1">urn:schemas-xmlsoap-org:ws:2005:04:discovery</a:To>
	//</Header>
	//<Body><Probe xmlns="http://schemas.xmlsoap.org/ws/2005/04/discovery">
	//<d:Types xmlns:d="http://schemas.xmlsoap.org/ws/2005/04/discovery" xmlns:dp0="http://www.onvif.org/ver10/network/wsdl">dp0:NetworkVideoTransmitter</d:Types>
	//</Probe>
	//</Body>
	//</Envelope>`

	return sendUDPMulticast(probeSOAP.String(), interfaceName)

}

func sendUDPMulticast(msg string, interfaceName string) []string {
	var result []string
	data := []byte(msg)
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		fmt.Println(err)
	}
	group := net.IPv4(239, 255, 255, 250)

	c, err := net.ListenPacket("udp4", "0.0.0.0:1024")
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()

	p := ipv4.NewPacketConn(c)
	if err := p.JoinGroup(iface, &net.UDPAddr{IP: group}); err != nil {
		fmt.Println(err)
	}

	dst := &net.UDPAddr{IP: group, Port: 3702}
	for _, ifi := range []*net.Interface{iface} {
		if err := p.SetMulticastInterface(ifi); err != nil {
			fmt.Println(err)
		}
		p.SetMulticastTTL(2)
		if _, err := p.WriteTo(data, nil, dst); err != nil {
			fmt.Println(err)
		}
	}

	if err := p.SetReadDeadline(time.Now().Add(time.Second * 1)); err != nil {
		log.Fatal(err)
	}

	for {
		b := make([]byte, bufSize)
		n, _, _, err := p.ReadFrom(b)
		if err != nil {
			fmt.Println(err)
			break
		}
		result = append(result, string(b[0:n]))
	}
	return result
}
