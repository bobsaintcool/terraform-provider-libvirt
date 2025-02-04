package libvirt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	UUIDBuflen = 16
)

func TestAccLibvirtNetworkDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "libvirt_network" "default" {
name = "default"
}`,
				Check: resource.ComposeTestCheckFunc(
					checkTemplate("data.libvirt_network.default", "name", "default"),
					func(state *terraform.State) error {
						rs, err := getResourceValue("data.libvirt_network.default", "id", state)
						if err != nil {
							return err
						}
						value := parseUUID(rs)
						if len(value) != UUIDBuflen {
							return fmt.Errorf(
								"Network UUID format not expected %v", uuidString(value))
						}
						return nil
					},
				),
			},
		},
	})
}

func TestAccLibvirtNetworkDataSourceNonExistingNetwork(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "libvirt_network" "NonExistingNetwork" {
name = "NonExistingNetworkName"
}`,
				Check: resource.ComposeTestCheckFunc(
					func(state *terraform.State) error {
						_, err := getResourceValue("data.libvirt_network.NonExistingNetwork", "id", state)
						if err != nil {
							return nil
						}

						return fmt.Errorf(
							"No libvirt network is expected with this name")
					},
				),
			},
		},
	})
}

func TestAccLibvirtNetworkDataSource_DNSHostTemplate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{

			{
				Config: `data "libvirt_network_dns_host_template" "bootstrap" {
  count = 2
  ip = "1.1.1.${count.index}"
  hostname = "myhost${count.index}"
}`,
				Check: resource.ComposeTestCheckFunc(
					checkTemplate("data.libvirt_network_dns_host_template.bootstrap.0", "ip", "1.1.1.0"),
					checkTemplate("data.libvirt_network_dns_host_template.bootstrap.0", "hostname", "myhost0"),
					checkTemplate("data.libvirt_network_dns_host_template.bootstrap.1", "ip", "1.1.1.1"),
					checkTemplate("data.libvirt_network_dns_host_template.bootstrap.1", "hostname", "myhost1"),
				),
			},
		},
	})
}

func getResourceValue(id string, name string, state *terraform.State) (string, error) {
	rs, err := getResourceFromTerraformState(id, state)
	if err != nil {
		return "", err
	}

	return rs.Primary.Attributes[name], nil
}

func checkTemplate(id, name, value string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		v, err := getResourceValue(id, name, state)
		if err != nil {
			return err
		}

		if v != value {
			return fmt.Errorf(
				"Value for %s is %s, not %s", name, v, value)
		}

		return nil
	}
}

func TestAccLibvirtNetworkDataSource_DNSSRVTemplate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{

			{
				Config: `data "libvirt_network_dns_srv_template" "etcd_cluster" {
  count = 2
  service = "etcd-server-ssl"
  protocol = "tcp"
  target = "my-etcd-${count.index}.tt.testing"
}`,
				Check: resource.ComposeTestCheckFunc(
					checkTemplate("data.libvirt_network_dns_srv_template.etcd_cluster.0", "target", "my-etcd-0.tt.testing"),
					checkTemplate("data.libvirt_network_dns_srv_template.etcd_cluster.0", "service", "etcd-server-ssl"),
					checkTemplate("data.libvirt_network_dns_srv_template.etcd_cluster.0", "protocol", "tcp"),
					checkTemplate("data.libvirt_network_dns_srv_template.etcd_cluster.1", "target", "my-etcd-1.tt.testing"),
					checkTemplate("data.libvirt_network_dns_srv_template.etcd_cluster.1", "service", "etcd-server-ssl"),
					checkTemplate("data.libvirt_network_dns_srv_template.etcd_cluster.1", "protocol", "tcp"),
				),
			},
		},
	})
}

func TestAccLibvirtNetworkDataSource_DNSDnsmasqTemplate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{

			{
				Config: `data "libvirt_network_dnsmasq_options_template" "options" {
  count = 2
  option_name = "address"
  option_value = "/.apps.tt${count.index}.testing/1.1.1.${count.index+1}"
}`,
				Check: resource.ComposeTestCheckFunc(
					checkTemplate("data.libvirt_network_dnsmasq_options_template.options.0", "option_name", "address"),
					checkTemplate("data.libvirt_network_dnsmasq_options_template.options.0", "option_value", "/.apps.tt0.testing/1.1.1.1"),
					checkTemplate("data.libvirt_network_dnsmasq_options_template.options.1", "option_name", "address"),
					checkTemplate("data.libvirt_network_dnsmasq_options_template.options.1", "option_value", "/.apps.tt1.testing/1.1.1.2"),
				),
			},
		},
	})
}
