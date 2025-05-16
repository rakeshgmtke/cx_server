# cx_server

A high-performance DIAMETER server for the Cx interface in IMS networks.

[![Go](https://img.shields.io/badge/Go-lang-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)

## Overview

`cx_server` is a Go-based implementation of a DIAMETER server specifically designed to handle the Cx interface as defined in the 3GPP specifications ETSI TS 129 228 and ETSI TS 129 229. It efficiently processes key Cx interface messages such as User-Authorization-Request (UAR), Multimedia-Authentication-Request (MAR), Server-Assignment-Request (SAR), and Location-Information-Request (LIR). This server is built with performance in mind to meet the demands of modern IMS networks. IMS restoration procedures are also considered as outlined in ETSI TS 123 380.


## Key Features

* **Cx Interface Support:** Handles core Cx interface messages: UAR, MAR, SAR, and LIR.
* **Standards Compliance:** Adheres to ETSI TS 129 228 and ETSI TS 129 229 specifications for the Cx interface based on the DIAMETER protocol.
* **Performance Optimized:** Developed in Go for high concurrency and efficient resource utilization.
* **Statistics:** real-time server statistics.
* **pprof Support:** Includes standard Go profiling for performance analysis.
* **IMS Restoration Awareness:** Takes into account procedures defined in ETSI TS 123 380 for IMS restoration.
* **Configurable:** Uses command-line arguments for easy configuration of essential parameters.


## Getting Started

### Prerequisites

* [Go](https://go.dev/doc/install) (version 1.18 or higher)

### Running the Server

You can run the `cx_server` directly from the command line after building it (if you haven't already). Here's an example of how to start the server with specific configurations:

**./cx_server -diam_host hss1.ims.com -diam_realm ims.com -pprof_addr :9000 -addr 10.1.10.10:3868 -ifcxml ifc_default.xml**

Command-line Arguments
The cx_server accepts the following command-line arguments:

**-addr** 
	  Address in the form of ip:port to listen on (default: "127.0.0.1:3868").
	
**-cert_file** 
	  Path to the TLS certificate file (optional). Enables TLS if provided.
	
**-diam_host** 
	  Diameter identity host (default: "server").
	
**-diam_realm** 
	  Diameter identity realm (default: "go-diameter").
	
**-ifcxml** 
	  Path to the User-Data XML file. If not provided, it will use the default ifc_default.xml file present in the current directory.
	
**-key_file**
	  Path to the TLS key file (optional). Required if -cert_file is provided.
	
**-log** 
	  Enable logging to a file (default: "false").

**-logpath** 
	  Path to the log file (default: "/tmp/hss.log"). This is used only if -log is enabled.

**-network_type**
	  Protocol type to use for DIAMETER communication (tcp or sctp) (default: "tcp").

**-pprof_addr**
	  Address in the form of ip:port for the pprof server (default: ":9000").
