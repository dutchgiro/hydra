package supervisor_test

import (
	hydra_config "github.com/innotech/hydra/config"
	. "github.com/innotech/hydra/etcd"
	. "github.com/innotech/hydra/supervisor"

	mock "github.com/innotech/hydra/supervisor/mock"
	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"

	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"io/ioutil"
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
		mockCtrl *gomock.Controller
		// etcdFactoryCache *EtcdFactory
		mockEtcdFactory *mock.MockEtcdBuilder
		mockEtcd        *mock.MockEtcdService
		etcdManager     *EtcdManager
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		// etcdFactoryCache = EtcdFactory
		// EtcdFactory = mock.NewMockEtcdBuilder(mockCtrl)
		mockEtcdFactory = mock.NewMockEtcdBuilder(mockCtrl)
		mockEtcd = mock.NewMockEtcdService(mockCtrl)
		etcdManager = NewEtcdManager()
	})

	AfterEach(func() {
		mockCtrl.Finish()
		// EtcdFactory = etcdFactoryCache
	})

	shouldStartSuccessfully := func(call *gomock.Call) {
		It("should configure and start etcd servers", func() {
			EtcdFactory = mockEtcdFactory
			c1 := EtcdFactory.EXPECT().Config(gomock.Any()).
				Times(1).After(call)
			c2 := EtcdFactory.EXPECT().Build().
				Return(mockEtcd).Times(1).After(c1)
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
			dir, _ := ioutil.TempDir("", "")
			file, _ := ioutil.TempFile(dir, "")
			file.WriteString("Test content")
			file.Close()
			etcdManager.EtcdService.Config.DataDir = dir

			call = mockEtcd.EXPECT().Stop().Times(1)
			etcdManager.Stop()
			exists, err := existsPath(dir)
			Expect(err).ToNot(HaveOccurred())
			Expect(exists).To(BeFalse())
		})
		return
	}

	Describe("Start", func() {
		shouldStartSuccessfully(nil)
		// It("should configure and start etcd servers", func() {
		// 	c1 := EtcdFactory.EXPECT().Config(gomock.Any()).
		// 		Times(1)
		// 	c2 := EtcdFactory.EXPECT().Build().
		// 		Return(mockEtcd).Times(1).After(c1)
		// 	mockEtcd.EXPECT().Start().
		// 		Times(1).After(c2)

		// 	etcdManager.Start()
		// })
	})
	Describe("Stop", func() {
		shouldStopSuccessfully()
		// It("should close etcd listeners", func() {
		// 	dir, _ := ioutil.TempFile("", "")
		// 	file, _ := ioutil.TempFile(dir, "")
		// 	file.WriteString("Test content")
		// 	file.Close()
		// 	etcdManager.etcdService.Config.DataDir = dir

		// 	mockEtcd.EXPECT().Stop().Times(1)
		// 	etcdManager.Stop()
		// 	exists, err := existsPath(dir)
		// 	Expect(err).ToNot(HaveOcurred())
		// 	Expect(exists).To(BeFalse())
		// })
	})
	Describe("Restart", func() {
		callToStop := shouldStopSuccessfully()
		shouldStartSuccessfully(callToStop)
		// TODO: Call Stop and Start tests
	})
})
