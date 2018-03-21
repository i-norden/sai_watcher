package gov_test

import (
	"github.com/8thlight/sai_watcher/event_triggered/tub/gov"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gov conversion", func() {
	It("converts a gov entity to gov model", func() {
		var_ := "0x6361700000000000000000000000000000000000000000000000000000000000"
		ge := &gov.GovEntity{
			Block: 1,
			Tx:    "0x1b0272e44fb92b9bfc654946e6aead6ea3d2a29a3b2d5f4c95170dc6617411e7",
			Var:   var_,
			Arg:   "0x000000000000000000000000000000000000000000295be96e64066972000000",
			Guy:   "0x000000000000000000000000f2c5369cffb8ea6284452b0326e326dbfdcb867c",
			Cap:   "50000000000000000000000000",
			Mat:   "1500000000000000000000000000",
			Tax:   "1000000000000000000000000000",
			Fee:   "1000000000158153903837946257",
			Axe:   "1130000000000000000000000000",
			Gap:   "1000000000000000000",
		}
		gm := gov.ConvertToModel(ge)
		Expect(gm.Block).To(Equal(int64(1)))
		Expect(gm.Tx).To(Equal(ge.Tx))
		Expect(gm.Var).To(Equal("cap"))
		Expect(gm.Arg).To(Equal("50000000"))
		Expect(gm.Cap).To(Equal("50000000"))
		Expect(gm.Guy).To(Equal("0xf2c5369cffb8ea6284452b0326e326dbfdcb867c"))
		Expect(gm.Mat).To(Equal("1.5"))
		Expect(gm.Tax).To(Equal("1"))
		Expect(gm.Fee).To(Equal("1.000000000158154"))
		Expect(gm.Axe).To(Equal("1.13"))
		Expect(gm.Gap).To(Equal("1"))
	})

	It("removes zero a string", func() {
		var_ := "0x6361700000000000000000000000000000000000000000000000000000000000"
		p := gov.RemoveZeroPadding(var_)
		Expect(p).To(Equal("0x636170"))

	})

})
