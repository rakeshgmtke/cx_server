# cx_server
Diameter Server to handle Cx interface


Usage of ./cx_server:
  -addr string
        address in the form of ip:port to listen on (default "127.0.0.1:3868")
  -cert_file string
        tls certificate file (optional)
  -diam_host string
        diameter identity host (default "server")
  -diam_realm string
        diameter identity realm (default "go-diameter")
  -ifcxml string
        Path to the User-Data XML file or it will use default ifc_default.xml file present in current directory
  -key_file string
        tls key file (optional)
  -log
        Enable logging to a file
  -logpath string
        Path to the log file (default "/tmp/hss.log")
  -network_type string
        protocol type tcp/sctp (default "tcp")
  -pprof_addr string
        address in form of ip:port for the pprof server (default ":9000")
