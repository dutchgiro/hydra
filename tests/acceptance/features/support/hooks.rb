# Before('@standalone') do
# 	HydraServiceManager.start_standalone
# end

# Before('@cluster') do
# 	HydraServiceManager.start_cluster 3
# end

# After('@standalone, @cluster') do
# 	HydraServiceManager.kill_all_nodes
# end

After do
	kill_all_nodes
end