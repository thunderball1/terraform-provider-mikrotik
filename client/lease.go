package client

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-routeros/routeros/proto"
)

type DhcpLease struct {
	Id         string
	Address    string
	MacAddress string
	Comment    string
	Server     string
}

func (client Mikrotik) AddDhcpLease(address, macaddress, name string, server string) (*DhcpLease, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := strings.Split(fmt.Sprintf("/ip/dhcp-server/lease/add =address=%s =mac-address=%s =comment=%s =server=%s", address, macaddress, name, server), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	id := r.Done.Map["ret"]

	return client.FindDhcpLease(id)
}

func (client Mikrotik) FindDhcpLease(id string) (*DhcpLease, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := strings.Split(fmt.Sprintf("/ip/dhcp-server/lease/print ?.id=%s", id), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	found := false
	var sentence *proto.Sentence

	if err != nil {
		return nil, err
	}

	for _, reply := range r.Re {
		for _, item := range reply.List {
			if item.Key == ".id" && item.Value == id {
				found = true
				sentence = reply
				log.Printf("[DEBUG] Found dhcp lease we were looking for: %v", sentence)
			}
		}
	}

	if !found {
		return nil, nil
	}

	address := ""
	macaddress := ""
	comment := ""
	server := ""
	for _, pair := range sentence.List {
		if pair.Key == "address" {
			address = pair.Value
		}
		if pair.Key == "mac-address" {
			macaddress = pair.Value
		}
		if pair.Key == "comment" {
			comment = pair.Value
		}
		if pair.Key == "server" {
			server = pair.Value
		}
	}

	return &DhcpLease{
		Id:         id,
		MacAddress: macaddress,
		Address:    address,
		Comment:    comment,
		Server:     server,
	}, nil
}

func (client Mikrotik) UpdateDhcpLease(id, address, macaddress, comment string, server string) (*DhcpLease, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	cmd := strings.Split(fmt.Sprintf("/ip/dhcp-server/lease/set =.id=%s =address=%s =mac-address=%s =comment=%s =server=%s", id, address, macaddress, comment, server), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	return client.FindDhcpLease(id)
}

func (client Mikrotik) DeleteDhcpLease(id string) error {
	c, err := client.getMikrotikClient()

	if err != nil {
		return err
	}

	cmd := strings.Split(fmt.Sprintf("/ip/dhcp-server/lease/remove =.id=%s", id), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	return err
}
