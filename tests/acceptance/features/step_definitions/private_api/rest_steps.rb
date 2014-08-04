require 'httparty'

When(/^the client requests DELETE (.*)$/) do |path|
  @last_response = HTTParty.delete('http://' + @private_api_addr + path)
end

When(/^the client requests GET (.*)$/) do |path|
  @last_response = HTTParty.get('http://' + @private_api_addr + path)
end

When(/^the client requests POST (.*)$/) do |path, json|
  @last_response = HTTParty.post('http://' + @private_api_addr + path,
	  { 
	    :body => json,
	    :headers => { 'Content-Type' => 'application/json', 'Accept' => 'application/json'}
	  })
end

Then(/^the response should be (\d+) .*$/) do |status_code|
  @last_response.code == status_code
end

Then(/^the response should be JSON:$/) do |json|
  JSON.parse(@last_response.body).should == JSON.parse(json)
end