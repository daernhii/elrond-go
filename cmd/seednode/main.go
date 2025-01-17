package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/display"
	"github.com/ElrondNetwork/elrond-go/hashing/sha256"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/p2p/libp2p"
	"github.com/ElrondNetwork/elrond-go/p2p/libp2p/discovery"
	factoryP2P "github.com/ElrondNetwork/elrond-go/p2p/libp2p/factory"
	"github.com/ElrondNetwork/elrond-go/p2p/loadBalancer"
	"github.com/btcsuite/btcd/btcec"
	libp2pCrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/urfave/cli"
)

var (
	seedNodeHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
VERSION:
   {{.Version}}
   {{end}}
`
	// port defines a flag for setting the port on which the node will listen for connections
	port = cli.IntFlag{
		Name:  "port",
		Usage: "Port number on which the application will start",
		Value: 10000,
	}
	// p2pSeed defines a flag to be used as a seed when generating P2P credentials. Useful for seed nodes.
	p2pSeed = cli.StringFlag{
		Name:  "p2p-seed",
		Usage: "P2P seed will be used when generating credentials for p2p component. Can be any string.",
		Value: "seed",
	}

	p2pConfigurationFile = "./config/p2p.toml"

	errNilSeed                     = errors.New("nil seed")
	errEmotySeed                   = errors.New("empty seed")
	errNilBuffer                   = errors.New("nil buffer")
	errEmptyBuffer                 = errors.New("empty buffer")
	errInvalidPort                 = errors.New("cannot start node on port < 0")
	errPeerDiscoveryShouldBeKadDht = errors.New("kad-dht peer discovery should have been enabled")
)

type seedRandReader struct {
	index int
	seed  []byte
}

// NewSeedRandReader will return a new instance of a seed-based reader
func NewSeedRandReader(seed []byte) *seedRandReader {
	return &seedRandReader{seed: seed, index: 0}
}

// Read to provided buffer pseudo-random generated bytes
func (srr *seedRandReader) Read(p []byte) (n int, err error) {
	if srr.seed == nil {
		return 0, errNilSeed
	}
	if len(srr.seed) == 0 {
		return 0, errEmotySeed
	}
	if p == nil {
		return 0, errNilBuffer
	}
	if len(p) == 0 {
		return 0, errEmptyBuffer
	}
	for i := 0; i < len(p); i++ {
		p[i] = srr.seed[srr.index]
		srr.index = (srr.index + 1) % len(srr.seed)
	}

	return len(p), nil
}

func main() {
	app := cli.NewApp()
	cli.AppHelpTemplate = seedNodeHelpTemplate
	app.Name = "SeedNode CLI App"
	app.Usage = "This is the entry point for starting a new seed node - the app will help bootnodes connect to the network"
	app.Flags = []cli.Flag{port, p2pSeed}
	app.Version = "v0.0.1"
	app.Authors = []cli.Author{
		{
			Name:  "The Elrond Team",
			Email: "contact@elrond.com",
		},
	}

	app.Action = func(c *cli.Context) error {
		return startNode(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func startNode(ctx *cli.Context) error {
	fmt.Println("Starting node...")

	stop := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	p2pConfig, err := core.LoadP2PConfig(p2pConfigurationFile)
	if err != nil {
		return err
	}
	fmt.Printf("Initialized with p2p config from: %s\n", p2pConfigurationFile)
	if ctx.IsSet(port.Name) {
		p2pConfig.Node.Port = ctx.GlobalInt(port.Name)
	}
	if ctx.IsSet(p2pSeed.Name) {
		p2pConfig.Node.Seed = ctx.GlobalString(p2pSeed.Name)
	}

	fmt.Println("Seed node....")
	messenger, err := createNode(p2pConfig)
	if err != nil {
		return err
	}
	err = messenger.Bootstrap()
	if err != nil {
		return err
	}

	go func() {
		<-sigs
		fmt.Println("terminating at user's signal...")
		stop <- true
	}()

	fmt.Println("Application is now running...")
	displayMessengerInfo(messenger)
	for {
		select {
		case <-stop:
			return nil
		case <-time.After(time.Second * 5):
			displayMessengerInfo(messenger)
		}
	}
}

func createNode(p2pConfig *config.P2PConfig) (p2p.Messenger, error) {
	var randReader io.Reader
	if p2pConfig.Node.Seed != "" {
		hasher := sha256.Sha256{}
		randReader = NewSeedRandReader(hasher.Compute(p2pConfig.Node.Seed))
	} else {
		randReader = rand.Reader
	}

	netMessenger, err := createNetMessenger(p2pConfig, randReader)
	if err != nil {
		return nil, err
	}

	return netMessenger, nil
}

func createNetMessenger(
	p2pConfig *config.P2PConfig,
	randReader io.Reader,
) (p2p.Messenger, error) {

	if p2pConfig.Node.Port < 0 {
		return nil, errInvalidPort
	}

	pDiscoveryFactory := factoryP2P.NewPeerDiscovererCreator(*p2pConfig)
	pDiscoverer, err := pDiscoveryFactory.CreatePeerDiscoverer()
	if err != nil {
		return nil, err
	}
	_, ok := pDiscoverer.(*discovery.KadDhtDiscoverer)
	if !ok {
		return nil, errPeerDiscoveryShouldBeKadDht
	}

	fmt.Printf("Starting with peer discovery: %s\n", pDiscoverer.Name())

	prvKey, _ := ecdsa.GenerateKey(btcec.S256(), randReader)
	sk := (*libp2pCrypto.Secp256k1PrivateKey)(prvKey)

	nm, err := libp2p.NewNetworkMessenger(
		context.Background(),
		p2pConfig.Node.Port,
		sk,
		nil,
		loadBalancer.NewOutgoingChannelLoadBalancer(),
		pDiscoverer,
		libp2p.ListenAddrWithIp4AndTcp,
		p2pConfig.Node.TargetPeerCount,
	)

	if err != nil {
		return nil, err
	}

	return nm, nil
}

func displayMessengerInfo(messenger p2p.Messenger) {
	headerSeedAddresses := []string{"Seednode addresses:"}
	addresses := make([]*display.LineData, 0)

	for _, address := range messenger.Addresses() {
		addresses = append(addresses, display.NewLineData(false, []string{address}))
	}

	tbl, _ := display.CreateTableString(headerSeedAddresses, addresses)
	fmt.Println(tbl)

	mesConnectedAddrs := messenger.ConnectedAddresses()
	sort.Slice(mesConnectedAddrs, func(i, j int) bool {
		return strings.Compare(mesConnectedAddrs[i], mesConnectedAddrs[j]) < 0
	})

	headerConnectedAddresses := []string{
		fmt.Sprintf("Seednode is connected to %d peers:", len(mesConnectedAddrs))}
	connAddresses := make([]*display.LineData, len(mesConnectedAddrs))

	for idx, address := range mesConnectedAddrs {
		connAddresses[idx] = display.NewLineData(false, []string{address})
	}

	tbl2, _ := display.CreateTableString(headerConnectedAddresses, connAddresses)
	fmt.Println(tbl2)
}
