package supervisor_test

import (
	. "github.com/innotech/hydra/supervisor"
	mock "github.com/innotech/hydra/supervisor/mock"

	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	// . "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	// "errors"
	"time"
)

var _ = Describe("PeersMonitor", func() {
	const (
		peerAddrItself                   string        = "127.0.0.1:7701"
		stateKeyTTL                      uint64        = 5
		durationBetweenPublicationsState time.Duration = time.Duration(500) * time.Millisecond
	)

	var (
		// ch             chan StateControllerState
		mockCtrl       *gomock.Controller
		mockEtcdClient *mock.MockEtcdRequester
		peersMonitor   *PeersMonitor
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEtcdClient = mock.NewMockEtcdRequester(mockCtrl)
		peersMonitor = NewPeersMonitor(mockEtcdClient)
		// ch = make(chan StateControllerState)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Run", func() {
		It("should request the cluster directory periodically", func() {
			c1 := mockEtcdClient.EXPECT().Get(gomock.Eq("cluster"), gomock.Eq(true), gomock.Eq(true)).
				Times(1)
			mockEtcdClient.EXPECT().Get(gomock.Eq("cluster"), gomock.Eq(true), gomock.Eq(true)).
				AnyTimes().After(c1)

			go func() {
				// stateManager.Run(ch)
				peersMonitor.Run()
			}()
			time.Sleep(DefaultRequestClusterInterval + time.Duration(1)*time.Second)
		})
	})
})
