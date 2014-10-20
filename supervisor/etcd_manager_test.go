package supervisor_test

import (
	hydra_config "github.com/innotech/hydra/config"
	. "github.com/innotech/hydra/etcd"
	. "github.com/innotech/hydra/supervisor"

	mock "github.com/innotech/hydra/supervisor/mock"
	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"

	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"

	"os"
)

func existsPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

var _ = FDescribe("EtcdManager", func() {
	var (
		mockCtrl        *gomock.Controller
		mockEtcdFactory *mock.MockEtcdBuilder
		mockEtcd        *mock.MockEtcdService
		etcdManager     *EtcdManager
	)

	etcdFactoryCache := EtcdFactory

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEtcdFactory = mock.NewMockEtcdBuilder(mockCtrl)
		mockEtcd = mock.NewMockEtcdService(mockCtrl)
		etcdManager = NewEtcdManager()
	})

	AfterEach(func() {
		mockCtrl.Finish()
		EtcdFactory = etcdFactoryCache
	})

	shouldStartSuccessfully := func(call *gomock.Call) {
		It("should configure and start etcd servers", func() {
			EtcdFactory = mockEtcdFactory
			var c1 *gomock.Call
			if call == nil {
				c1 = mockEtcdFactory.EXPECT().Config(gomock.Any()).
					Times(1)
			} else {
				c1 = mockEtcdFactory.EXPECT().Config(gomock.Any()).
					Times(1).After(call)
			}
			c2 := mockEtcdFactory.EXPECT().Build().
				Return(mockEtcd).
				Times(1).After(c1)
			mockEtcd.EXPECT().Start().
				Times(1).After(c2)

			hydraConfig := hydra_config.New()
			hydraConfig.EtcdConf.BindAddr = "127.0.0.1:7401"
			hydraConfig.EtcdConf.Name = "hydra_0"
			hydraConfig.EtcdConf.Peer.BindAddr = "127.0.0.1:7701"
			etcdManager.Start(hydraConfig.EtcdConf)
		})
	}

	shouldStopSuccessfully := func() (call *gomock.Call) {
		It("should close etcd listeners", func() {
			etcdManager.EtcdService = mockEtcd
			call = mockEtcd.EXPECT().Stop().Times(1)
			etcdManager.Stop()
		})
		return
	}

	Describe("Start", func() {
		shouldStartSuccessfully(nil)
	})
	Describe("Stop", func() {
		shouldStopSuccessfully()
	})
	Describe("Restart", func() {
		callToStop := shouldStopSuccessfully()
		shouldStartSuccessfully(callToStop)
	})
})
