package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// LogEntry represents the structure of each log entry.
type LogEntry struct {
	Timestamp       string `json:"timestamp"`
	Action          string `json:"action"`
	Username        string `json:"username"`
	DN              string `json:"dn"`
	Status          string `json:"status"`
	ResultCode      int    `json:"result_code"`
	MessageResponse string `json:"message_response"`
}

var (
	openVPNLogEntriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "openvpn_ldap_user_logs_total",
			Help: "Total number of OpenVPN LDAP user log entries",
		},
		[]string{
			"Timestamp",
			"action",
			"status",
			"Username",
			"DN",
			"ResultCode",
			"MessageResponse"},
	)
)

func init() {
	prometheus.MustRegister(openVPNLogEntriesTotal)
}

func main() {

	var (
		listenAddress    = flag.String("web.listen-address", ":9177", "Address to listen on for web interface and telemetry.")
		metricsPath      = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		openvpnAuthPaths = flag.String("openvpn.auth_paths", "/go/src/app/openvpn-status-auth.log", "Paths at which OpenVPN places its auth files.")
	)
	flag.Parse()

	log.Printf("Starting OpenVPN Exporter - Auth logs\n")
	log.Printf("Listen address: %v\n", *listenAddress)
	log.Printf("Metrics path: %v\n", *metricsPath)
	log.Printf("openvpn.auth_log_path: %v\n", *openvpnAuthPaths)

	// Specify the path to your log file
	authFilePath := *openvpnAuthPaths

	// Start a goroutine to read and process the log file
	go readLogEntries(authFilePath)

	// Expose the metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*listenAddress, nil)
}

func readLogEntries(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		processLogEntry(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading log file:", err)
	}
}

func processLogEntry(logEntryJSON string) {
	var logEntry LogEntry
	err := json.Unmarshal([]byte(logEntryJSON), &logEntry)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	// Increment the appropriate metric based on log entry attributes
	openVPNLogEntriesTotal.WithLabelValues(
		logEntry.Timestamp,
		logEntry.Action,
		logEntry.Status,
		logEntry.Username,
		logEntry.DN,
		strconv.Itoa(logEntry.ResultCode),
		logEntry.MessageResponse).Inc()
}
