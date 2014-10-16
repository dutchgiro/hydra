package etcd_test

import (
	hydra_config "github.com/innotech/hydra/config"
	. "github.com/innotech/hydra/etcd"

	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"
)

var _ = Describe("EtcdFactory", func() {
	Describe("Build", func() {
		It("should build a etcd object", func() {
			var etcd *Etcd
			hydraConfig := hydra_config.New()
			// TODO: Add BindAddrs to default configuration
			hydraConfig.EtcdConf.BindAddr = "127.0.0.1:7401"
			hydraConfig.EtcdConf.Name = "hydra_0"
			hydraConfig.EtcdConf.Peer.BindAddr = "127.0.0.1:7701"
			EtcdFactory.Config(hydraConfig.EtcdConf)
			Expect(EtcdFactory.Build()).To(BeAssignableToTypeOf(etcd))
		})
	})
})
