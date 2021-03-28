package main

import (
	"embed"
	"github.com/Richard87/wg-vpn-server/api"
	"github.com/Richard87/wg-vpn-server/config"
	"github.com/Richard87/wg-vpn-server/database"
	"github.com/Richard87/wg-vpn-server/wireguard"
	"log"

	_ "github.com/mdlayher/genetlink"
	_ "github.com/mdlayher/netlink"
	_ "github.com/mdlayher/netlink/nlenc"
	_ "golang.zx2c4.com/wireguard/device"
	"os"
	"os/signal"
	"syscall"
)

//go:embed ui/build
var embededFiles embed.FS

func main() {
	config.InitFlags()
	database.InitDatabase()
	database.InitUsers()

	wireguard.Init()
	wireguard.Run()

	api.Run(embededFiles)

	termSignal := make(chan os.Signal)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(termSignal, os.Interrupt)
	<-termSignal // Block until we receive our signal.
	log.Println("Closing...")
	api.Close()
	wireguard.Close()

	os.Exit(0)
}
