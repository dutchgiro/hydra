package supervisor_test

import (
	. "github.com/innotech/hydra/supervisor"

	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"
)

var _ = FDescribe("PeerCluster", func() {
	var (
		initPeers   []Peer
		peerCluster *PeerCluster
	)

	BeforeEach(func() {
		initPeers = []Peer{
			Peer{
				Addr:     "98.245.153.111:4001",
				PeerAddr: "98.245.153.111:7001",
				State:    PeerStateEnabled,
			},
			Peer{
				Addr:     "98.245.153.112:4001",
				PeerAddr: "98.245.153.112:7001",
				State:    PeerStateEnabled,
			},
		}
		peerCluster = NewPeerCluster(initPeers)
	})

	// AfterEach(func() {
	// 	mockCtrl.Finish()
	// })

	Describe("Next", func() {
		It("should return the next value until end the collection", func() {
			for i := 0; i < len(initPeers); i++ {
				peer, err := peerCluster.Next()
				Expect(peer).To(Equal(initPeers[i]))
				Expect(err).NotTo(HaveOccurred())
			}
			peer, err := peerCluster.Next()
			Expect(peer).To(Equal(Peer{}))
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("HasNext", func() {
		It("should return true until all peers have been requested", func() {
			for i := 0; i < len(initPeers); i++ {
				Expect(peerCluster.HasNext()).To(BeTrue())
			}
			Expect(peerCluster.HasNext()).To(BeFalse())
		})
	})

	// Describe("Merge", func() {

	// })
	// Describe("Next", func() {
	// 	It("should return true as many times as peers have", func() {
	// 		for i := 0; i < len(initPeers); i++ {
	// 			Expect(peerCluster.Next()).To(BeTrue())
	// 		}
	// 		Expect(peerCluster.Next()).To(BeFalse())
	// 	})
	// })
	// Describe("Reset", func() {

	// })
	// Describe("Value", func() {
	// 	It("should return the current value", func() {
	// 		for i := 0; i < len(initPeers); i++ {
	// 			peerCluster.Next()
	// 			Expect(peerCluster.Value()).To(Equal(initPeers[i]))
	// 		}
	// 	})
	// })
})
