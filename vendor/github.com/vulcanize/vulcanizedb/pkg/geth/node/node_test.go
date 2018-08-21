package node_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
)

var EmpytHeaderHash = "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"

var _ = Describe("Parity Node Info", func() {

	It("verifies parity_versionInfo can be unmarshalled into ParityNodeInfo", func() {
		var parityNodeInfo core.ParityNodeInfo
		nodeInfoJSON := []byte(
			`{
        "hash": "0x2ae8b4ca278dd7b896090366615fef81cbbbc0e0",
        "track": "null",
        "version": {
          "major": 1,
          "minor": 6,
          "patch": 0
             }
        }`)
		json.Unmarshal(nodeInfoJSON, &parityNodeInfo)
		Expect(parityNodeInfo.Hash).To(Equal("0x2ae8b4ca278dd7b896090366615fef81cbbbc0e0"))
		Expect(parityNodeInfo.Track).To(Equal("null"))
		Expect(parityNodeInfo.Major).To(Equal(1))
		Expect(parityNodeInfo.Minor).To(Equal(6))
		Expect(parityNodeInfo.Patch).To(Equal(0))
	})

	It("Creates client string", func() {
		parityNodeInfo := core.ParityNodeInfo{
			Track: "null",
			ParityVersion: core.ParityVersion{
				Major: 1,
				Minor: 6,
				Patch: 0,
			},
			Hash: "0x1232144j",
		}
		Expect(parityNodeInfo.String()).To(Equal("Parity/v1.6.0/"))
	})

	It("returns the genesis block for any client", func() {
		client := fakes.NewMockRpcClient()
		n := node.MakeNode(client)
		Expect(n.GenesisBlock).To(Equal(EmpytHeaderHash))
	})

	It("returns the network id for any client", func() {
		client := fakes.NewMockRpcClient()
		n := node.MakeNode(client)
		Expect(n.NetworkID).To(Equal(float64(1234)))
	})

	It("returns parity ID and client name for parity node", func() {
		client := fakes.NewMockRpcClient()
		client.SetNodeType(core.PARITY)
		n := node.MakeNode(client)
		Expect(n.ID).To(Equal("ParityNode"))
		Expect(n.ClientName).To(Equal("Parity/v1.2.3/"))
	})

	It("returns geth ID and client name for geth node", func() {
		client := fakes.NewMockRpcClient()
		client.SetNodeType(core.GETH)
		n := node.MakeNode(client)
		Expect(n.ID).To(Equal("enode://GethNode@172.17.0.1:30303"))
		Expect(n.ClientName).To(Equal("Geth/v1.7"))
	})

	It("returns infura ID and client name for infura node", func() {
		client := fakes.NewMockRpcClient()
		client.SetNodeType(core.INFURA)
		client.SetIpcPath("infura/path")
		n := node.MakeNode(client)
		Expect(n.ID).To(Equal("infura"))
		Expect(n.ClientName).To(Equal("infura"))
	})
})
