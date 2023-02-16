package cmd

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"go-protobuf-tcp/client"
	"go-protobuf-tcp/server"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-protobuf-tcp",
	Short: "a simple TCP server and client using Protobuf",
	Long:  `go-protobuf-tcp is a simple TCP server and client using Protobuf`,
	Run:   run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var serverArg bool
var clientArg bool
var SrcPortArg int
var DstAddrArg string

func init() {
	rootCmd.Flags().BoolVarP(&serverArg, "server", "s", false, "run as server")
	rootCmd.Flags().BoolVarP(&clientArg, "client", "c", false, "run as client")
	rootCmd.Flags().IntVarP(&SrcPortArg, "src-port", "p", 0,
		"source port, can be used to specify the port of the client, or the port of the server to listen to")
	rootCmd.Flags().StringVarP(&DstAddrArg, "dst-ip", "d", "127.0.0.1:8080",
		"destination address, only for client, can be used to specify the address of the server to connect to")
}

func run(cmd *cobra.Command, args []string) {
	if serverArg {
		// start server
		s := server.NewServer(net.ParseIP("127.0.0.1"), SrcPortArg)
		s.Listen()
	} else {
		// start client
		addr := strings.Split(DstAddrArg, ":")[:2]
		if len(addr) != 2 {
			log.Fatal("invalid dst address")
		}
		dstPort, _ := strconv.Atoi(addr[1])

		c := client.NewClient(net.ParseIP(addr[0]), dstPort, SrcPortArg)
		err := c.Connect()
		if err != nil {
			log.Fatal(err)
		}
		c.Start()
	}
}
