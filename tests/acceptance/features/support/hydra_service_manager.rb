module HydraServiceManager
	def services
		@services ||= []
	end

	def start_hydra_service(node_index, args = {})
		service = HydraService.new(node_index, args)
		service.start
		services << service
	end

	def start_standalone(args = {})
		start_hydra_service(1, args)
	end

	def start_cluster(num_of_nodes, leader_args = {})
		raise "Invalid number of nodes, the minimal number is 2" if num_of_nodes < 2
		start_hydra_service(1, leader_args)
		(2..num_of_nodes).each do |i|
			start_hydra_service(i)
		end
	end

	def kill_all_nodes
		services.each do |service|
			service.stop
		end if services
	end
end

World(HydraServiceManager)