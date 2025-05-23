// Copyright 2013-2015 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter server example. This is by no means a complete server.
//
// If you'd like to test diameter over SSL, generate SSL certificates:
//   go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
//
// And start the server with `-cert_file cert.pem -key_file key.pem`.
//
// By default this server runs in a single OS thread. If you want to
// make it run on more, set the GOMAXPROCS=n environment variable.
// See Go's FAQ for details: http://golang.org/doc/faq#Why_no_multi_CPU

package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	//"time"
	//"runtime"
	//"io"
	//"sync"
	//	"strings"

	_ "net/http/pprof"

	"github.com/rakeshgmtke/go-diameter/diam"
	"github.com/rakeshgmtke/go-diameter/diam/avp"
	"github.com/rakeshgmtke/go-diameter/diam/datatype"
	_ "github.com/rakeshgmtke/go-diameter/diam/dict"
	"github.com/rakeshgmtke/go-diameter/diam/sm"
)

var (
	buildTime string
	version   string
)

const (
	VENDOR_3GPP         = 10415
	MandatoryCapability = 3
)

func main() {
	addr := flag.String("addr", "127.0.0.1:3868", "address in the form of ip:port to listen on")
	ppaddr := flag.String("pprof_addr", ":9000", "address in form of ip:port for the pprof server")
	host := flag.String("diam_host", "server", "diameter identity host")
	realm := flag.String("diam_realm", "go-diameter", "diameter identity realm")
	certFile := flag.String("cert_file", "", "tls certificate file (optional)")
	keyFile := flag.String("key_file", "", "tls key file (optional)")
	networkType := flag.String("network_type", "tcp", "protocol type tcp/sctp")
	enableLogging := flag.Bool("log", false, "Enable logging to a file")
	logFilePath := flag.String("logpath", "/tmp/hss.log", "Path to the log file")
	ifcXmlFile := flag.String("ifcxml", "", "Path to the User-Data XML file or it will use default ifc_default.xml file present in current directory")
	versionFlag := flag.Bool("version", false, "Print application version and build time")

	flag.Parse()

	// Check if the -version flag was provided
	if *versionFlag {
		fmt.Printf("Version: %s\n", defaultString(version, "N/A"))
		fmt.Printf("Build Time: %s\n", defaultString(buildTime, "N/A"))
		return // Exit after printing version information
	}	

	if *enableLogging {
		log.Printf("Version: %s\n", defaultString(version, "N/A"))
		log.Printf("Build Time: %s\n", defaultString(buildTime, "N/A"))
	}

	log.Printf("Build Version: %s Build Time: %s\n", defaultString(version, "N/A"), defaultString(buildTime, "N/A"))

	/*
		// Check is ifcxml file is passed
		if *ifcXmlFile == "" {
			fmt.Println("Error: XML file path is required. Please provide a value for the -ifcxml flag.")
			flag.Usage() // Show usage instructions
			os.Exit(1)
		}

	*/
	// Check if the ifcxml flag was provided and assign default if absent
	if *ifcXmlFile == "" {
		fmt.Println("Warning: No XML file path provided. Using the default file: ifc_default.xml")
		*ifcXmlFile = "ifc_default.xml" // Assign the default file path
	}

	// Check IFC XML syntax
	err := checkXMLSyntax(*ifcXmlFile)
	if err != nil {
		log.Fatalf("Error: XML syntax check failed for file '%s': %v", *ifcXmlFile, err)
	}
	log.Printf("XML file '%s' syntax is valid.\n", *ifcXmlFile)

	// Initialize the map.
	HSSDataMap = make(map[string]HSSData)
	//HSSDataMap = make(map[string]HSSData, 10000)

	log.Printf("HSSData initiation : %s and Length of HSSData: %d", HSSDataMap, len(HSSDataMap))

	populateHSS(10000) // Moved population to a function.

	log.Printf("HSSData initiated Length of HSSData: %d", len(HSSDataMap))

	//Function to start Diameter Stats
	stats := NewDiameterStats()

	// Generate the report
	ReportHeading := fmt.Sprintf("HostName: %s Realm: %s\nListening on : %s", *host, *realm, *addr)
	go printMetrics(stats, ReportHeading)

	settings := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity(*host),
		OriginRealm:      datatype.DiameterIdentity(*realm),
		VendorID:         VENDOR_3GPP,
		ProductName:      "go-diameter-cx",
		FirmwareRevision: 1,
	}

	// Open log file if logging is enabled
	if *enableLogging {
		logFile, err := os.OpenFile(*logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	// Create the state machine (mux) and set .CollectGarbage(context.Background(), &protos.Void{})its message handlers.
	mux := sm.New(settings)

	mux.Handle("UAR", handleUAR(*settings, stats, *enableLogging))
	mux.Handle("MAR", handleMAR(*settings, stats, *enableLogging))
	mux.Handle("SAR", handleSAR(*settings, stats, *enableLogging, *ifcXmlFile))
	mux.Handle("LIR", handleLIR(*settings, stats, *enableLogging))

	mux.HandleFunc("ALL", handleALL) // Catch all.

	// Print error reports.
	go printErrors(mux.ErrorReports())

	if len(*ppaddr) > 0 {
		go func() { log.Fatal(http.ListenAndServe(*ppaddr, nil)) }()
	}

	err = listen(*networkType, *addr, *certFile, *keyFile, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func checkXMLSyntax(xmlFilePath string) error {
	xmlData, err := os.ReadFile(xmlFilePath)
	if err != nil {
		return fmt.Errorf("error reading XML file '%s': %v", xmlFilePath, err)
	}

	// Attempt to unmarshal the XML data into a generic map or struct
	// We don't need the specific structure for syntax checking.
	var v interface{}
	err = xml.Unmarshal(xmlData, &v)
	if err != nil {
		return fmt.Errorf("error parsing XML syntax in file '%s': %v", xmlFilePath, err)
	}

	return nil // XML syntax is valid
}

func getTelNumber(s string) string {
	result := "tel:"
	for _, c := range s {
		if (c >= '0' && c <= '9') || c == '+' {
			result += string(c)
		}
	}
	return result
}

func onlyNumbers(s string) string {
	result := ""
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result += string(c)
		}
	}
	return result
}

func printErrors(ec <-chan *diam.ErrorReport) {
	for err := range ec {
		log.Println(err)
	}
}

func listen(networkType, addr, cert, key string, handler diam.Handler) error {
	// Start listening for connections.
	if len(cert) > 0 && len(key) > 0 {
		log.Println("Starting secure diameter server on", addr)
		return diam.ListenAndServeNetworkTLS(networkType, addr, cert, key, handler, nil)
	}
	log.Println("Starting diameter server on", addr)
	return diam.ListenAndServeNetwork(networkType, addr, handler, nil)
}

func handleALL(c diam.Conn, m *diam.Message) {
	//stats.IncrementReceived("HandleAll", "", "HandleAll")
	//if enableLogging {
		log.Printf("Received unexpected message from %s:\n%s", c.RemoteAddr(), m)
//	}	
}

// Assume msg is a *diam.Message
func parseGroupedAVP(msg *diam.Message) {
	// Find the grouped AVP — replace with actual AVP code you’re parsing
	groupedAVP, err := msg.FindAVP(avp.VendorSpecificApplicationID, 0)
	if err != nil {
		//log.Printf("Grouped AVP not found: %v", err)
		return
	}

	// Check that the AVP is of grouped type
	groupedData, ok := groupedAVP.Data.(*diam.GroupedAVP)
	if !ok {
		log.Println("Not a grouped AVP (expected *diam.GroupedAVP)")
		return
	}

	// Iterate the inner AVPs
	for _, innerAVP := range groupedData.AVP {
		fmt.Printf("AVP Code parseGroupedAVP --------->: %d, VendorID: %d, Data: %v\n",
			innerAVP.Code, innerAVP.VendorID, innerAVP.Data)
	}
}

// GetGroupedAVPs extracts a grouped AVP (by code and vendor ID) from a *diam.Message.
// It returns a slice of inner *diam.AVPs (inside the grouped AVP).
func GetGroupedAVPs(msg *diam.Message, code, vendor uint32) ([]*diam.AVP, error) {
	avpItem, err := msg.FindAVP(code, vendor)
	if err != nil {
		return nil, fmt.Errorf("AVP %d not found: %w", code, err)
	}

	grouped, ok := avpItem.Data.(*diam.GroupedAVP)
	if !ok {
		return nil, fmt.Errorf("AVP %d is not grouped", code)
	}
	return grouped.AVP, nil
}

func getIPFromAddress(address string) (string, error) {
	host, _, err := net.SplitHostPort(address)
	if err == nil {
		ip := net.ParseIP(host)
		if ip != nil {
			return ip.String(), nil
		}
		return "", fmt.Errorf("invalid IP address in host part: %s", host)
	}
	// If SplitHostPort fails, try parsing the whole address as IP
	ip := net.ParseIP(address)
	if ip != nil {
		return ip.String(), nil
	}
	return "", fmt.Errorf("invalid address format: %s", address)
}

// defaultString returns the input string if it's not empty, otherwise returns the defaultVal.
// This is a helper to display "N/A" if build flags aren't set.
func defaultString(s, defaultVal string) string {
	if s == "" {
		return defaultVal
	}
	return s
}
