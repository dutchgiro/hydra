package supervisor_test

// import (
// 	. "github.com/innotech/hydra/supervisor"

// 	mock "github.com/innotech/hydra/supervisor/mock"
// 	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"

// 	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
// 	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"
// )

// var _ = Describe("EtcdManager", func() {
// 	var (
// 		mockCtrl       *gomock.Controller
// 		mockEtcd       *mock.MockEtcd
// 		etcdSupervisor *EtcdSupervisor
// 		stateKeyPath   string
// 	)

// 	BeforeEach(func() {
// 		mockCtrl = gomock.NewController(GinkgoT())
// 		mockEtcd = mock.NewMockEtcd(mockCtrl)
// 		etcdManager = NewEtcdManager(mockEtcd)
// 	})

// 	AfterEach(func() {
// 		mockCtrl.Finish()
// 	})

// 	Describe("Stop", func() {
// 		c1 := mockEtcd.EXPECT().Set(gomock.Eq(stateKeyPath), gomock.Eq(StateAvailable), gomock.Eq(etcdSupervisor.GetStateTTL())).
// 			Return(successResponse, nil).Times(1)
// 	})
// 	Describe("Start", func() {

// 	})
// 	Describe("Restart", func() {

// 	})
// })
