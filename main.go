package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/JeffreyVdb/slack-working-from-home/slack"
	"github.com/JeffreyVdb/slack-working-from-home/util"
)

type publicIPResponse struct {
	IP string `json:"ip"`
}

type configuration struct {
	slackAPIToken       string
	slackStatusFilePath string
}

type statusEntry struct {
	StatusText  string   `json:"status_text"`
	StatusEmoji string   `json:"status_emoji,omitempty"`
	WifiNames   []string `json:"wifi_names"`
	PublicIPs   []string `json:"public_ips,omitempty"`
}

var programConfig = new(configuration)
var outputVersion bool

var publicIPServices = []string{
	"https://api.ipify.org/?format=json",
}

func getSlackStatusFromFile(filename, ssid, publicIP string) (*slack.SlackStatus, error) {
	fileHandle, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fileHandle.Close()

	config := []statusEntry{}
	err = json.NewDecoder(fileHandle).Decode(&config)
	if err != nil {
		return nil, err
	}

	for _, entry := range config {
		hasIPConstraint := len(entry.PublicIPs) > 0

		if util.ContainsString(entry.WifiNames, ssid) {
			if !hasIPConstraint || util.ContainsString(entry.PublicIPs, publicIP) {
				return &slack.SlackStatus{StatusText: entry.StatusText, StatusEmoji: entry.StatusEmoji}, nil
			}
		}
	}

	return nil, nil
}

func getPublicIPAddress(ipServices []string) (string, error) {
	ipService := ipServices[rand.Intn(len(ipServices))]
	resp, err := http.Get(ipService)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	response := &publicIPResponse{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return "", err
	}

	return response.IP, nil
}

const linuxCmd = "iwgetid"
const macosCmd = "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport"

func getCmdOutput(program string, args ...string) (string, error) {
	output, err := exec.Command(program, args...).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func getLinuxWifiSSID() (string, error) {
	return getCmdOutput(linuxCmd, "--raw")
}

func getMacSSID() (string, error) {
	output, err := getCmdOutput(macosCmd, "-I")
	if err != nil {
		return "", err
	}

	r := regexp.MustCompile(`s*SSID: (.+)s*`)
	name := r.FindAllStringSubmatch(output, -1)

	if len(name) <= 1 {
		return "", fmt.Errorf("Could not get SSID")
	} else {
		return name[1][1], nil
	}
}

func getWifiSSID() (string, error) {
	platform := runtime.GOOS
	var err error
	var ssid string
	if platform == "darwin" {
		ssid, err = getMacSSID()
	} else if platform == "win32" {
		err = fmt.Errorf("Not yet implemented for win32")
	} else {
		ssid, err = getLinuxWifiSSID()
	}

	return ssid, err
}

func readEnvironment(config *configuration) error {
	if value, ok := os.LookupEnv("SLACK_API_TOKEN"); ok {
		config.slackAPIToken = value
	} else {
		return fmt.Errorf("SLACK_API_TOKEN environment variable was not defined")
	}

	return nil
}

func init() {
	rand.Seed(time.Now().Unix())

	flag.StringVar(&programConfig.slackStatusFilePath, "status-file", "/etc/slack-status.json", "Slack status file path")
	flag.BoolVar(&outputVersion, "version", false, "Get version of this program")
}

func main() {
	flag.Parse()
	if outputVersion {
		fmt.Println("Version 1.0")
		fmt.Printf("Go version: %s\n", runtime.Version())
		os.Exit(0)
	}

	err := readEnvironment(programConfig)
	if err != nil {
	}

	ssid, err := getWifiSSID()
	if err != nil {
		util.PrintError(os.Stderr, "Error while trying to get SSID: %v\n", err)
		os.Exit(1)
	}

	publicIP, err := getPublicIPAddress(publicIPServices)
	if err != nil {
		util.PrintError(os.Stderr, "Error while trying to get public IP: %v\n", err)
		os.Exit(1)
	}

	status, err := getSlackStatusFromFile(programConfig.slackStatusFilePath, ssid, publicIP)
	if err != nil {
		util.PrintError(os.Stderr, "Error while trying to get status: %v\n", err)
		os.Exit(1)
	}

	if status == nil {
		fmt.Printf("No status for ssid: %s\n", ssid)
		os.Exit(0)
	}

	slackClient, err := slack.NewClient(programConfig.slackAPIToken, slack.Timeout(10))
	if err != nil {
		util.PrintError(os.Stderr, "Error while trying to create slack client: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Setting status...")
	err = slackClient.SetProfileStatus(status)
	if err != nil {
		util.PrintError(os.Stderr, "Error setting slack status: %v\n", err)
		os.Exit(1)
	}
}
