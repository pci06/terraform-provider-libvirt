package libvirt

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/dmacvicar/terraform-provider-libvirt/libvirt/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	//"libvirt.org/go/libvirtxml"
	"log"
	"strconv"
	//libvirt "github.com/digitalocean/go-libvirt"
)

type NetworkGeneric struct {
	Network xml.Name `xml:"network"`
	Name string `xml:"name"`
	UUID string `xml:"uuid"`
	Forward struct {
		Forward xml.Name `xml:"forward"`
		Mode string `xml:"mode,attr"`
	} `xml:"forward"`
	Bridge struct {
		Bridge xml.Name `xml:"bridge"`
		Name string `xml:"name,attr"`
		STP string `xml:"stp,attr"`
		Delay string `xml:"delay,attr"`
	} `xml:"bridge"`
	MAC struct {
		MAC xml.Name `xml:"mac"`
		Address string `xml:"address,attr"`
	} `xml:"mac"`
	IP struct {
		IP xml.Name `xml:"ip"`
		Address string `xml:"address,attr"`
		NetMask string `xml:"netmask,attr"`
		DHCP struct {
			DHCP xml.Name `xml:"dhcp"`
			Range struct {
				Range xml.Name `xml:"range"`
				Start string `xml:"start,attr"`
				End string `xml:"end,attr"`
			} `xml:"range"`
		} `xml:"dhcp"`
	} `xml:"ip"`
}

func datasourceLibvirtNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Info(ctx, "XXX datasourceLibvirtNetwoCreate XXX")
	//client := meta.(*Client)
	//virConn := client.libvirt
	return datasourceLibvirtNetworkRead(ctx, d, meta)
}

func datasourceLibvirtNetwork() *schema.Resource {
	log.Printf("[DEBUG] create schema")
	return &schema.Resource {
		ReadContext: datasourceLibvirtNetworkRead,
		Schema: map[string]*schema.Schema {
			"name": {
				Type: schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type: schema.TypeString,
				Optional: true,
			},
			"forward": {
				Type: schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func datasourceLibvirtNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Read data source libvirt_network")

	virConn := meta.(*Client).libvirt

	var networkName string

	if name, ok := d.GetOk("name"); ok {
		networkName = name.(string)
		tflog.Debug(ctx, "Got name: ", map[string]interface{}{ "networkName": networkName })
	}

	network, err := virConn.NetworkLookupByName(networkName)
	if err != nil {
		return diag.Errorf("failed to lookup network: %w", err)
	}

	xmlDesc, err := virConn.NetworkGetXMLDesc(network, 0)
	if err != nil {
		return diag.Errorf("failed to get XML for network: %w", err)
	}

	networkXML := NetworkGeneric{}

	err = xml.Unmarshal([]byte(xmlDesc), &networkXML)
	if err != nil {
		tflog.Error(ctx, "failed to unmarshal XML into networkXML:", map[string]interface{}{"error": err})
	}
	tflog.Debug(ctx, "XXX Parsed network into networkXML: ", map[string]interface{}{"xml": networkXML})

	d.Set("xml", xmlDesc)
	d.Set("uuid", networkXML.UUID)
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(networkXML.Forward)
	json.Unmarshal(inrec, &inInterface)
	tflog.Debug(ctx, "Setting interface to:", inInterface)
	d.Set("forward", &inInterface)
//	d.Set("path", networkXML.Path)
//	d.Set("parent", networkXML.Parent)
	//log.Printf("[DEBUG] d.Set capability type %s : %s", capability[0]["type"], capability)
	d.SetId(strconv.Itoa(hashcode.String(fmt.Sprintf("%v", xmlDesc))))

	return nil
}
