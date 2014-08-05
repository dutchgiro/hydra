PRIVATE_API_ADDR = '127.0.0.1:7601'

Given(/^one hydra server that is running without an application configuration file \(apps\.json\)$/) do
	@private_api_addr = PRIVATE_API_ADDR
	args = {
		'apps-file' => FIXTURES_PATH + 'non-existent-apps-file.json',
		'private-api-addr' => @private_api_addr
	}
  start_standalone args
end

Given(/^one hydra server that is running with an application configuration file \(apps\.json\) containing an empty array of applications$/) do
  @private_api_addr = PRIVATE_API_ADDR
	args = {
		'apps-file' => FIXTURES_PATH + 'empty_array.json',
		'private-api-addr' => @private_api_addr
	}
  start_standalone args
end

Given(/^one hydra server that is running with an application configuration file \(apps\.json\) containing only one application$/) do
  @private_api_addr = PRIVATE_API_ADDR
	args = {
		'apps-file' => FIXTURES_PATH + 'only_one_app_v1.json',
		'private-api-addr' => @private_api_addr
	}
  start_standalone args
end

Given(/^one hydra server that is running with an application configuration file \(apps\.json\) containing two applications$/) do
  @private_api_addr = PRIVATE_API_ADDR
	args = {
		'apps-file' => FIXTURES_PATH + 'two_apps.json',
		'private-api-addr' => @private_api_addr
	}
  start_standalone args
end

Given(/^a hydra server cluster that is running with an application configuration file \(apps\.json\) containing only one application$/) do
  @private_api_addr = PRIVATE_API_ADDR
	args = {
		'apps-file' => FIXTURES_PATH + 'only_one_app_v2.json',
	}
  start_cluster(3, args)
end