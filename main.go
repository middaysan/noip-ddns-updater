package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Username       string
	Password       string
	Hostname       string
	Interval       int
	NoIPURL        string
	CheckIPURL     string
	LastUpdatedIP  string
}

func getPublicIP(checkIPURL string) (string, error) {
	resp, err := http.Get(checkIPURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(ip)), nil
}

func updateNoIP(config *Config, ip string) error {
	client := &http.Client{}
	requestURL := fmt.Sprintf("%s?hostname=%s&myip=%s", config.NoIPURL, url.QueryEscape(config.Hostname), ip)
	fmt.Printf("Updating IP: %s\n", requestURL)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(config.Username, config.Password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("No-IP Response:", string(body))
	return nil
}

func getEnvVariable(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func validateConfig(config *Config) {
	if config.Username == "" || config.Password == "" || config.Hostname == "" {
		fmt.Println("Missing required environment variables: NOIP_USER, NOIP_PASS, NOIP_HOST")
		fmt.Println("Username:", config.Username)
		fmt.Println("Password len*:", len(config.Password))
		fmt.Println("Hostname:", config.Hostname)
		fmt.Println("Update Interval:", config.Interval, "minutes")
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("No-IP Dynamic DNS Updater")
	fmt.Println("Usage: noip_ddns_updater [-h]")
	fmt.Println("Environment Variables:")
	fmt.Println("  NOIP_USER              - No-IP account username")
	fmt.Println("  NOIP_PASS              - No-IP account password")
	fmt.Println("  NOIP_HOST              - Hostname to update")
	fmt.Println("  NOIP_INTERVAL_MINUTES  - Interval between updates (default: 5 minutes)")
	fmt.Println("  NOIP_URL               - No-IP update URL (default: https://dynupdate.no-ip.com/nic/update)")
	fmt.Println("  CHECK_IP_URL           - URL to check public IP (default: https://checkip.amazonaws.com)")
	os.Exit(0)
}

func loadConfig() *Config {
	if len(os.Args) > 1 && os.Args[1] == "-h" {
		printUsage()
	}

	interval, err := strconv.Atoi(os.Getenv("NOIP_INTERVAL_MINUTES"))
	if err != nil {
		fmt.Println("Error reading interval, setting default to 5 minutes.")
		interval = 5
	}

	config := &Config{
		Username:   getEnvVariable("NOIP_USER", ""),
		Password:   getEnvVariable("NOIP_PASS", ""),
		Hostname:   getEnvVariable("NOIP_HOST", ""),
		Interval:   interval,
		NoIPURL:    getEnvVariable("NOIP_URL", "https://dynupdate.no-ip.com/nic/update"),
		CheckIPURL: getEnvVariable("CHECK_IP_URL", "https://checkip.amazonaws.com"),
	}

	validateConfig(config)

	fmt.Println("Configuration Loaded:")
	fmt.Println("Username:", config.Username)
	fmt.Println("Hostname:", config.Hostname)
	fmt.Println("Update Interval:", config.Interval, "minutes")
	fmt.Println("No-IP Update URL:", config.NoIPURL)
	fmt.Println("Check IP URL:", config.CheckIPURL)

	return config
}

func runUpdater(config *Config) {
	ticker := time.NewTicker(time.Duration(config.Interval) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		currentIP, err := getPublicIP(config.CheckIPURL)
		if err != nil {
			fmt.Println("Error retrieving public IP:", err)
			continue
		}

		if currentIP == config.LastUpdatedIP {
			fmt.Println("IP address has not changed. Skipping update.")
			continue
		}

		fmt.Println("New IP detected:", currentIP)
		if err := updateNoIP(config, currentIP); err == nil {
			config.LastUpdatedIP = currentIP
		}
	}
}

func main() {
	config := loadConfig()
	runUpdater(config)
}
