package supervisor_test

import (
	. "github.com/innotech/hydra/etcd"
	. "github.com/innotech/hydra/supervisor"

	mock "github.com/innotech/hydra/supervisor/mock"
	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"

	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	// . "github.com/innotech/hydra/vendors/github.com/onsi/gomega"
)

var _ = Describe("EtcdManager", func() {
	var (
		mockCtrl         *gomock.Controller
		etcdFactoryCache *EtcdFactory
		// mockEtcdFactory *mock.MockEtcdFactory
		mockEtcd    *mock.MockEtcdService
		etcdManager *EtcdManager
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		etcdFactoryCache = EtcdFactory
		EtcdFactory = mock.NewMockEtcdFactory(mockCtrl)
		mockEtcd = mock.NewMockEtcdService()
		etcdManager = NewEtcdManager(mockEtcd)
	})

	AfterEach(func() {
		mockCtrl.Finish()
		EtcdFactory = etcdFactoryCache
	})

	Describe("Start", func() {
		It("should configure and start etcd servers", func() {
			c1 := EtcdFactory.EXPECT().Config(gomock.Any()).
				Times(1)
			c2 := EtcdFactory.EXPECT().Build().
				Return(mockEtcd).Times(1).After(c1)
			mockEtcd.EXPECT().Start().
				Times(1).After(c2)

			etcdManager.Start()
		})
	})
	Describe("Stop", func() {
		It("should close etcd listeners", func() {
			const mockEtcdDataDir string = "/tmp/mockEtcdDataDir"
			etcdManager.etcdService.Config.DataDir = mockEtcdDataDir
			mockEtcd.EXPECT().Stop().Times(1)

			// TODO: need config and check remove files

			etcdManager.Stop()
		})
	})
	Describe("Restart", func() {
		// TODO: Call Stop and Start tests
	})
})
