require 'childprocess'
require 'timeout'
require 'httparty'
require 'fileutils'

class HydraService
	HYDRA_BIN_PATH = Dir.getwd + '/../../bin/hydra'
	HYDRA_BASE_CONFIG = {
		'addr' => '127.0.0.1:740#[node_index]',
		# 'apps-file' = ""
		# 'bind-addr' = ""
		# 'ca-file' = ""
		# 'cert-file' = ""
		'data-dir' => "/tmp/hydra#[node_index]",
		# 'discovery' = ""
		# 'force' = false
		# 'instance-expiration-time' = 300
		# 'key-file' = ""
		'load-balancer-addr' => '*:750#[node_index]',
		'peers' => '',
		'private-api-addr' => '127.0.0.1:760#[node_index]',
		'public-api-addr' => '127.0.0.1:780#[node_index]',
		'name' => 'hydra#[node_index]',
		# 'snapshot' = true
		# 'snapshot-count' = 2000
		# 'verbose' = false
		# peer
		'peer-addr' => '127.0.0.1:770#[node_index]',
		# 'peer-bind-addr' = ""
		# 'peer-ca-file' = ""
		# 'peer-cert-file' = ""
		# 'peer-key-file' = ""
		# 'peer-heartbeat-timeout' = 100
		# 'peer-election-timeout' = 400
	}

	def initialize(node_index, args = {})
		@node_index = node_index
		@config = {}
		@node_proccess = nil
		customize_configuration(HYDRA_BASE_CONFIG.merge(args))
	end
	
	def start
		FileUtils.rm_rf(@config['data-dir'])
		sleep 0.5
		args = get_array_of_args
		@node_proccess = ChildProcess.build(HYDRA_BIN_PATH, *args)
		@node_proccess.start
		sleep 2
		Timeout.timeout(3) do
			loop do
				begin
					HTTParty.get('http://' + @config['private-api-addr'] + '/apps')
					break
				rescue Errno::ECONNREFUSED => try_again
					sleep 0.1
				end
			end
		end
	end

	def stop
		@node_proccess.stop
	end

private
	def customize_configuration(base_config)
		if @node_index > 1
			@config['peers'] = compose_peers_arg(base_config['peers'])
		end

		base_config.each_pair do |key, value|
			@config[key] = customize_port(value, @node_index)
		end
	end

	def customize_port(address, index)
		return address.gsub(/#\[node_index\]/, index.to_s)
	end

	def compose_peers_arg(base_peer_addr)
		peers_arg = ''
		(@node_index-1).times do |i|
			peers_arg += customize_port(base_peer_addr, i)
			peers_arg += ',' if i < @node_index-1
		end

		return peers_arg
	end

	def get_array_of_args
		args = []
		@config.each do |key, value|
			if value != ''
				args << make_arg_from_string(key)
				args << value
			end
		end
		return args
	end

	def make_arg_from_string(string)
		return '-' + string
	end
	
end