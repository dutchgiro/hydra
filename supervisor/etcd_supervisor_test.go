package supervisor_test

import (
	. "github.com/innotech/hydra/supervisor"
	mock "github.com/innotech/hydra/supervisor/mock"
	etcd_config "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"

	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"time"
)

var _ = Describe("EtcdSupervisor", func() {
	var (
		mockCtrl             *gomock.Controller
		mockEtcdManager      *mock.MockEtcdController
		mockClusterInspector *mock.MockClusterAnalyzer
		mockStateManager     *mock.MockStateController
		etcdConfig           *etcd_config.Config
		etcdSupervisor       *EtcdSupervisor
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEtcdManager = mock.NewMockEtcdController(mockCtrl)
		mockClusterInspector = mock.NewMockClusterAnalyzer(mockCtrl)
		mockStateManager = mock.NewMockStateController(mockCtrl)
		etcdConfig = etcd_config.New()
		etcdConfig.Peers = []string{"98.245.153.111:4001", "98.245.153.112:4001"}
		etcdSupervisor = NewEtcdSupervisor(etcdConfig)
		etcdSupervisor.EtcdManager = mockEtcdManager
		etcdSupervisor.ClusterInspector = mockClusterInspector
		etcdSupervisor.StateManager = mockStateManager
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Run", func() {
		It("should start the etcd service as master prime", func() {
			var etcdConf *etcd_config.Config
			c0 := mockEtcdManager.EXPECT().Start(gomock.Any()).
				Do(func(config *etcd_config.Config) {
				etcdConf = config
			}).Times(1)
			mockStateManager.EXPECT().Run(gomock.Any()).AnyTimes().After(c0)
			mockClusterInspector.EXPECT().Run(gomock.Any()).AnyTimes().After(c0)
			go func() {
				etcdSupervisor.Run()
			}()
			time.Sleep(time.Duration(100) * time.Millisecond)
			Expect(etcdConf).To(Equal(etcdConfig))
			Expect(etcdConf.Peers).To(HaveLen(0))
		})
		It("should initialize cluster inspector with configured peers", func() {
			initPeers := []string{"98.245.153.111:4001", "98.245.153.112:4001"}
			etcdConfig = etcd_config.New()
			etcdConfig.Peers = initPeers
			etcdSupervisor = NewEtcdSupervisor(etcdConfig)
			clusterInspector := etcdSupervisor.ClusterInspector.(*ClusterInspector)
			Expect(clusterInspector.PeerCluster.Peers).To(HaveLen(len(initPeers)))
			for i := 0; i < len(etcdConfig.Peers); i++ {
				Expect(clusterInspector.PeerCluster.Peers[0].Addr).To(Equal(initPeers[0]))
				Expect(clusterInspector.PeerCluster.Peers[0].State).To(Equal(PeerStateEnabled))
			}
		})
		Context("when state manager can not set state and emits a stop signal", func() {
			It("should restart the etcd service as master prime", func() {
				c0 := mockEtcdManager.EXPECT().Start(gomock.Any()).Times(1)
				mockClusterInspector.EXPECT().Run(gomock.Any()).Times(1).After(c0)
				var stateChannel chan StateControllerState
				c1 := mockStateManager.EXPECT().Run(gomock.Any()).Do(func(ch chan StateControllerState) {
					stateChannel = ch
				}).Times(1).After(c0)
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
				c0 := mockEtcdManager.EXPECT().Start(gomock.Any()).Times(1)
				mockStateManager.EXPECT().Run(gomock.Any()).Times(1).After(c0)
				var clusterInspectorChannel chan string
				c1 := mockClusterInspector.EXPECT().Run(gomock.Any()).Do(func(ch chan string) {
					clusterInspectorChannel = ch
				}).Times(1).After(c0)
				var etcdConf *etcd_config.Config
				mockEtcdManager.EXPECT().Restart(gomock.Any()).Do(func(config *etcd_config.Config) {
					etcdConf = config
				}).Times(1).After(c1)
				go func() {
					etcdSupervisor.Run()
				}()
				time.Sleep(time.Duration(100) * time.Millisecond)
				const newLeader string = "98.245.153.113:7001"
				clusterInspectorChannel <- newLeader
				time.Sleep(time.Duration(50) * time.Millisecond)
				Expect(etcdConf.Peers).To(HaveLen(1))
				Expect(etcdConf.Peers).To(ContainElement(newLeader))
			})
		})
	})
})
