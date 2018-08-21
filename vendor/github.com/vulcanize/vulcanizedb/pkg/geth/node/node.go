package node

import (
	"context"

	"strconv"

	"regexp"

	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"strings"
)

type IPropertiesReader interface {
	NodeInfo() (id string, name string)
	NetworkId() float64
	GenesisBlock() string
}

type PropertiesReader struct {
	client core.RpcClient
}

type ParityClient struct {
	PropertiesReader
}

type GethClient struct {
	PropertiesReader
}

type InfuraClient struct {
	PropertiesReader
}

func MakeNode(rpcClient core.RpcClient) core.Node {
	pr := makePropertiesReader(rpcClient)
	id, name := pr.NodeInfo()
	return core.Node{
		GenesisBlock: pr.GenesisBlock(),
		NetworkID:    pr.NetworkId(),
		ID:           id,
		ClientName:   name,
	}
}

func makePropertiesReader(client core.RpcClient) IPropertiesReader {
	switch getNodeType(client) {
	case core.GETH:
		return GethClient{PropertiesReader: PropertiesReader{client: client}}
	case core.PARITY:
		return ParityClient{PropertiesReader: PropertiesReader{client: client}}
	case core.INFURA:
		return InfuraClient{PropertiesReader: PropertiesReader{client: client}}
	default:
		return PropertiesReader{client: client}
	}
}

func getNodeType(client core.RpcClient) core.NodeType {
	if strings.Contains(client.IpcPath(), "infura") {
		return core.INFURA
	}
	modules, _ := client.SupportedModules()
	if _, ok := modules["admin"]; ok {
		return core.GETH
	}
	return core.PARITY
}

func (reader PropertiesReader) NetworkId() float64 {
	var version string
	err := reader.client.CallContext(context.Background(), &version, "net_version")
	if err != nil {
		log.Println(err)
	}
	networkId, _ := strconv.ParseFloat(version, 64)
	return networkId
}

func (reader PropertiesReader) GenesisBlock() string {
	var header *types.Header
	blockZero := "0x0"
	includeTransactions := false
	reader.client.CallContext(context.Background(), &header, "eth_getBlockByNumber", blockZero, includeTransactions)
	return header.Hash().Hex()
}

func (reader PropertiesReader) NodeInfo() (string, string) {
	var info p2p.NodeInfo
	reader.client.CallContext(context.Background(), &info, "admin_nodeInfo")
	return info.ID, info.Name
}

func (client ParityClient) NodeInfo() (string, string) {
	nodeInfo := client.parityNodeInfo()
	id := client.parityID()
	return id, nodeInfo
}

func (client InfuraClient) NodeInfo() (string, string) {
	return "infura", "infura"
}

func (client ParityClient) parityNodeInfo() string {
	var nodeInfo core.ParityNodeInfo
	client.client.CallContext(context.Background(), &nodeInfo, "parity_versionInfo")
	return nodeInfo.String()
}

func (client ParityClient) parityID() string {
	var enodeId = regexp.MustCompile(`^enode://(.+)@.+$`)
	var enodeURL string
	client.client.CallContext(context.Background(), &enodeURL, "parity_enode")
	enode := enodeId.FindStringSubmatch(enodeURL)
	if len(enode) < 2 {
		return ""
	}
	return enode[1]
}
