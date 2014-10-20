package main

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/innotech/hydra/config"
	"github.com/innotech/hydra/database/connector"
	"github.com/innotech/hydra/etcd"
	"github.com/innotech/hydra/load_balancer"
	"github.com/innotech/hydra/log"
	"github.com/innotech/hydra/server"
	"github.com/innotech/hydra/supervisor"
)

func main() {
	// Load configuration.
	var conf = config.New()
	if err := conf.Load(os.Args[1:]); err != nil {
		log.Fatal(err.Error() + "\n")
	}

	// Enable verbose option.
	if conf.Verbose {
		log.Verbose = true
	}

	if conf.DataDir == "" {
		log.Fatal("Data dir does't exist")
	}

	// Create data directory if it doesn't already exist.
	if err := os.MkdirAll(conf.DataDir, 0744); err != nil {
		log.Fatalf("Unable to create path: %s", err)
	}

	// Load etcd configuration.
	if err := conf.LoadEtcdConfig(); err != nil {
		log.Fatalf("Unable to load etcd conf: %s", err)
	}

	// Launch services
	// TODO: refactor visibility of supervisor component, export only NewEtcdSupervisor
	// TODO: Maybe pass conf.EtcdConf by reference
	etcdSupervisor := supervisor.NewEtcdSupervisor(conf.EtcdConf)
	go func() {
		etcdSupervisor.Run()
	}()

	// Private Server API
	privateHydraListener, err := net.Listen("tcp", conf.PrivateAPIAddr)
	if err != nil {
		log.Fatalf("Failed to create hydra private listener: ", err)
	}
	var privateServer = server.NewPrivateServer(privateHydraListener, conf.InstanceExpirationTime)
	privateServer.RegisterHandlers()
	go func() {
		log.Infof("hydra private server [name %s, listen on %s, advertised url %s]", conf.Name, conf.PrivateAPIAddr, "http://"+conf.PrivateAPIAddr)
		log.Fatal(http.Serve(privateServer.Listener, privateServer.Router))
	}()

	// Public Server API
	var loadBalancerFrontendEndpoint string = "ipc://" + conf.Name + "-frontend.ipc"
	publicHydraListener, err := net.Listen("tcp", conf.PublicAPIAddr)
	if err != nil {
		log.Fatalf("Failed to create hydra public listener: ", err)
	}
	var publicServer = server.NewPublicServer(publicHydraListener, loadBalancerFrontendEndpoint, conf.BalanceTimeout)
	publicServer.RegisterHandlers()
	go func() {
		log.Infof("hydra public server [name %s, listen on %s, advertised url %s]", conf.Name, conf.PublicAPIAddr, "http://"+conf.PublicAPIAddr)
		log.Fatal(http.Serve(publicServer.Listener, publicServer.Router))
	}()

	// Load applications.
	var appsConfig = config.NewApplicationsConfig()
	if _, err := os.Stat(conf.AppsFilePath); os.IsNotExist(err) {
		log.Warnf("Unable to find apps file: %s", err)
	} else {
		if err := appsConfig.Load(conf.AppsFilePath); err != nil {
			log.Fatalf("Unable to load applications: %s", err)
		}
	}

	time.Sleep(1 * time.Second)
	// Persist Configured applications
	if err := appsConfig.Persists(); err != nil {
		log.Fatalf("Failed to save configured applications: ", err)
	}

	// Load Balancer
	var inprocBackendEndpoint string = "ipc://" + conf.Name + "-backend.ipc"
	loadBalancer := load_balancer.NewLoadBalancer(loadBalancerFrontendEndpoint, "tcp://"+conf.LoadBalancerAddr, inprocBackendEndpoint)
	defer loadBalancer.Close()
	loadBalancer.Run()
}
