package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	//"runtime"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

const (
	serverURL      = "http://localhost:8080/api/v1" // Change this to your server URL in production
	submitPath     = "/resource"
	dnsSubmitPath  = "/internet-usage"
	retryInterval  = 5 * time.Second
	updateInterval = 10 * time.Second
	snapLen        = 1600
	promiscuous    = false
	timeout        = pcap.BlockForever
)

type ResourceData struct {
	ComputerID    string  `json:"computer_id"`
	CPU         float64 `json:"cpu"`
	Memory      float64 `json:"memory"`
	NetworkIn   float64 `json:"network_in"`
	NetworkOut  float64 `json:"network_out"`
}

type DNSData struct {
	ComputerID  string    `json:"computer_id"`
	Domain    string    `json:"domain"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	// Parse command line arguments
	computerID := flag.String("systemID", "", "System ID for this computer")
	flag.Parse()

	if *computerID == "" {
		log.Fatal("System ID is required. Use --systemID flag")
	}

	// Display privacy notice
	fmt.Println("NOTICE: This system monitors resource and internet usage for lab management purposes.")
	fmt.Printf("System ID: %s\n", *computerID)
	fmt.Println("Press Ctrl+C to stop monitoring.")

	// Create log file
	logFile := setupLogging()
	defer logFile.Close()

	// Start DNS monitoring in a separate goroutine
	go monitorDNS(*computerID)

	// Previous network stats for calculating rate
	var prevNetStats []net.IOCountersStat

	// Main monitoring loop for system resources
	for {
		data := ResourceData{
			ComputerID: *computerID,
		}

		// Get CPU usage
		cpuPercent, err := cpu.Percent(time.Second, false)
		if err != nil {
			log.Printf("Error getting CPU usage: %v", err)
		} else if len(cpuPercent) > 0 {
			data.CPU = cpuPercent[0]
		}

		// Get memory usage
		memInfo, err := mem.VirtualMemory()
		if err != nil {
			log.Printf("Error getting memory usage: %v", err)
		} else {
			data.Memory = memInfo.UsedPercent
		}

		// Get network usage
		netStats, err := net.IOCounters(false)
		if err != nil {
			log.Printf("Error getting network stats: %v", err)
		} else if len(netStats) > 0 {
			if prevNetStats != nil {
				// Calculate network rate (bytes per second)
				timeDiff := updateInterval.Seconds()
				data.NetworkIn = float64(netStats[0].BytesRecv-prevNetStats[0].BytesRecv) / timeDiff
				data.NetworkOut = float64(netStats[0].BytesSent-prevNetStats[0].BytesSent) / timeDiff
			}
			prevNetStats = netStats
		}

		// Submit data to server
		if err := submitData(data); err != nil {
			log.Printf("Error submitting data: %v", err)
			time.Sleep(retryInterval)
			continue
		}

		time.Sleep(updateInterval)
	}
}

func monitorDNS(computerID string) {
	// Find all network devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Printf("Error finding network devices: %v", err)
		return
	}

	// Create a WaitGroup to manage multiple packet captures
	var wg sync.WaitGroup

	// Start packet capture on each device
	for _, device := range devices {
		wg.Add(1)
		go func(deviceName string) {
			defer wg.Done()

			handle, err := pcap.OpenLive(deviceName, snapLen, promiscuous, timeout)
			if err != nil {
				log.Printf("Error opening device %s: %v", deviceName, err)
				return
			}
			defer handle.Close()

			// Set BPF filter for DNS queries
			err = handle.SetBPFFilter("udp and port 53")
			if err != nil {
				log.Printf("Error setting BPF filter on device %s: %v", deviceName, err)
				return
			}

			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			for packet := range packetSource.Packets() {
				dnsLayer := packet.Layer(layers.LayerTypeDNS)
				if dnsLayer == nil {
					continue
				}

				dns, _ := dnsLayer.(*layers.DNS)
				if !dns.QR { // Only process DNS queries, not responses
					for _, question := range dns.Questions {
						domain := string(question.Name)
						if domain != "" {
							dnsData := DNSData{
								ComputerID:  computerID,
								Domain:    domain,
								Timestamp: time.Now(),
							}
							
							// Submit DNS data to server
							if err := submitDNSData(dnsData); err != nil {
								log.Printf("Error submitting DNS data: %v", err)
							}
						}
					}
				}
			}
		}(device.Name)
	}

	wg.Wait()
}

func submitData(data ResourceData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}

	resp, err := http.Post(serverURL+submitPath, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	return nil
}

func submitDNSData(data DNSData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling DNS data: %v", err)
	}

	resp, err := http.Post(serverURL+dnsSubmitPath, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending DNS data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	return nil
}

func setupLogging() *os.File {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	// Open log file
	logFile, err := os.OpenFile(
		fmt.Sprintf("%s/collector_%s.log", logsDir, time.Now().Format("2006-01-02")),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Set log output to file
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return logFile
}
