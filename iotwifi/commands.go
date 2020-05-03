package iotwifi

import (
	"os/exec"

	"github.com/bhoriuchi/go-bunyan/bunyan"
)

// Command for device network commands.
type Command struct {
	Log      bunyan.Logger
	Runner   CmdRunner
	SetupCfg *SetupCfg
}

// RemoveApInterface removes the AP interface.
func (c *Command) RemoveApInterface() {
	c.Log.Info("##NG Removing AP interface")
	cmd := exec.Command("iw", "dev", "uap0", "del")
	cmd.Start()
	cmd.Wait()
	c.Log.Info("##NG Removed AP interface")
}

// ConfigureApInterface configured the AP interface.
func (c *Command) ConfigureApInterface() {
	cmd := exec.Command("ifconfig", "uap0", c.SetupCfg.HostApdCfg.Ip)
	cmd.Start()
	cmd.Wait()
}

// UpApInterface ups the AP Interface.
func (c *Command) UpApInterface() {
	cmd := exec.Command("ifconfig", "uap0", "up")
	cmd.Start()
	cmd.Wait()
}

// AddApInterface adds the AP interface.
func (c *Command) AddApInterface() {
	cmd := exec.Command("iw", "phy", "phy0", "interface", "add", "uap0", "type", "__ap")
	cmd.Start()
	cmd.Wait()
}

// CheckInterface checks the AP interface.
func (c *Command) CheckApInterface() {
	cmd := exec.Command("ifconfig", "uap0")
	go c.Runner.ProcessCmd("ifconfig_uap0", cmd)
}

// StartWpaSupplicant starts wpa_supplicant.
func (c *Command) StartWpaSupplicant() {

	args := []string{
		"-Dnl80211",
		"-iwlan0",
		"-B",
		"-c",
		c.SetupCfg.WpaSupplicantCfg.CfgFile,
	}

	cmd := exec.Command("wpa_supplicant", args...)
	go c.Runner.ProcessCmd("wpa_supplicant", cmd)
}

// StartAPDnsmasq starts dnsmasq on AP mode.
func (c *Command) StartAPDnsmasq() {
	// hostapd is enabled, fire up dnsmasq
	args := []string{
		"--no-hosts", // Don't read the hostnames in /etc/hosts.
		"--keep-in-foreground",
		"--log-queries",
		"--address=" + c.SetupCfg.DnsmasqCfg.Address,
		"--dhcp-range=" + c.SetupCfg.DnsmasqCfg.DhcpRange,
		"--dhcp-vendorclass=" + c.SetupCfg.DnsmasqCfg.VendorClass,
		"--dhcp-authoritative",
		"--log-facility=-",
		"--interface=uap0",
		"--port=0",
	}

	cmd := exec.Command("dnsmasq", args...)
	go c.Runner.ProcessCmd("dnsmasq", cmd)
}

// StartCLDnsmasq starts dnsmasq in CL mode.
func (c *Command) StartCLDnsmasq() {
	args := []string{
		"--no-hosts", // Don't read the hostnames in /etc/hosts.
		"--keep-in-foreground",
		"--log-queries",
		"--server=" + c.SetupCfg.DnsmasqCfg.Server,
		"--log-facility=-",
		"--interface=wlan0",
	}

	cmd := exec.Command("dnsmasq", args...)
	go c.Runner.ProcessCmd("dnsmasq", cmd)
}

func (c *Command) StartHostAPD() {
	args := []string{
		"/etc/hostapd/hostapd.conf",
	}

	cmd := exec.Command("hostapd", args...)
	go c.Runner.ProcessCmd("hostapd", cmd)
}

func (c *Command) killIt(it string) {
	c.Log.Info("##NG Killing [" + it + "]")
	args := []string{
		it,
	}

	cmd := exec.Command("killall", args...)
	cmdId := "killall " + it
	c.Runner.ProcessCmd(cmdId, cmd)
	c.Log.Info("##NG Killed [" + it + "]")
}
