package supervisor_test

import (
	hydra_config "github.com/innotech/hydra/config"
	. "github.com/innotech/hydra/supervisor"
	mock "github.com/innotech/hydra/supervisor/mock"
	etcd_config "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"

	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"time"
)

var _ = FDescribe("EtcdSupervisor", func() {
	var (
		mockCtrl             *gomock.Controller
		mockEtcdManager      *mock.MockEtcdController
		mockClusterInspector *mock.MockClusterAnalyzer
		mockStateManager     *mock.MockStateController
		hydraConfig          *hydra_config.Config
		etcdSupervisor       *EtcdSupervisor
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEtcdManager = mock.NewMockEtcdController(mockCtrl)
		mockClusterInspector = mock.NewMockClusterAnalyzer(mockCtrl)
		mockStateManager = mock.NewMockStateController(mockCtrl)
		hydraConfig = hydra_config.New()
		etcdSupervisor = NewEtcdSupervisor(hydraConfig)
		etcdSupervisor.EtcdManager = mockEtcdManager
		etcdSupervisor.ClusterInspector = mockClusterInspector
		etcdSupervisor.StateManager = mockStateManager
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Run", func() {
		Context("when state manager can not set state and emits a stop signal", func() {
			It("should restart the etcd service as master prime", func() {
				mockClusterInspector.EXPECT().Run(gomock.Any()).Times(1)
				var stateChannel chan StateControllerState
				c1 := mockStateManager.EXPECT().Run(gomock.Any()).Do(func(ch chan StateControllerState) {
					stateChannel = ch
				}).Times(1)
				var etcdConf *etcd_config.Config
				mockEtcdManager.EXPECT().Restart(gomock.Any()).Do(func(config *etcd_config.Config) {
					etcdConf = config
				}).Times(1).After(c1)
				go func() {
					etcdSupervisor.Run()
				}()
				time.Sleep(time.Duration(100) * time.Millisecond)
				stateChannel <- Stopped
				time.Sleep(time.Duration(50) * time.Millisecond)
				Expect(etcdConf.Peers).To(HaveLen(0))
			})
		})
		Context("when cluster inspector emits a new leader", func() {
			It("should restart the etcd service as slave", func() {
				mockStateManager.EXPECT().Run(gomock.Any()).Times(1)
				var clusterInspectorChannel chan string
				c1 := mockClusterInspector.EXPECT().Run(gomock.Any()).Do(func(ch chan string) {
					clusterInspectorChannel = ch
				}).Times(1)
				var etcdConf *etcd_config.Config
				mockEtcdManager.EXPECT().Restart(gomock.Any()).Do(func(config *etcd_config.Config) {
					etcdConf = config
				}).Times(1).After(c1)
				go func() {
					etcdSupervisor.Run()
				}()
				time.Sleep(time.Duration(100) * time.Millisecond)
				const peerAddr string = "98.245.153.113:7001"
				clusterInspectorChannel <- peerAddr
				time.Sleep(time.Duration(50) * time.Millisecond)
				Expect(etcdConf.Peers).To(HaveLen(1))
				Expect(etcdConf.Peers).To(ContainElement(peerAddr))
			})
		})
	})
})
