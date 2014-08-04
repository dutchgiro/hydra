# module HydraServiceManager
# 	HYDRA_BIN_PATH = Dir.getwd + '/../../../bin/hydra'
# 	HYDRA_BASE_CONFIG = {
# 		'addr' => '127.0.0.1:740#[node_index]',
# 		# 'apps-file' = ""
# 		# 'bind-addr' = ""
# 		# 'ca-file' = ""
# 		# 'cert-file' = ""
# 		# 'data-dir' = "."
# 		# 'discovery' = ""
# 		# 'force' = false
# 		# 'instance-expiration-time' = 300
# 		# 'key-file' = ""
# 		'load-balancer-addr' => '*:750#[node_index]',
# 		'peers' => [],
# 		'private-addr' => '127.0.0.1:760#[node_index]',
# 		'public-addr' => '127.0.0.1:780#[node_index]',
# 		'name' => 'hydra#[node_index]',
# 		# 'snapshot' = true
# 		# 'snapshot-count' = 2000
# 		# 'verbose' = false
# 		# peer
# 		'peer-addr' => '127.0.0.1:770#[node_index]',
# 		# 'peer-bind-addr' = ""
# 		# 'peer-ca-file' = ""
# 		# 'peer-cert-file' = ""
# 		# 'peer-key-file' = ""
# 		# 'peer-heartbeat-timeout' = 100
# 		# 'peer-election-timeout' = 400
# 	}

# 	extend self

# 	def services
# 		return @services if @services
# 		@services = []
# 	end

# 	def start_hydra_service(node_index)
# 		service = HydraServiceManager.HydraService.new(HYDRA_BASE_CONFIG, node_index)
# 		service.start
# 		@services << service
# 	end

# 	def start_standalone
# 		start_hydra_service(1)
# 	end

# 	def start_cluster(num_of_nodes)
# 		num_of_nodes.times do |i|
# 			start_hydra_service(i)
# 		end
# 	end

# 	def kill_all_nodes
# 		@services.each do |service|
# 			service.stop
# 		end 
# 	end
# end

# require "hydra_service_manager/hydra_service"