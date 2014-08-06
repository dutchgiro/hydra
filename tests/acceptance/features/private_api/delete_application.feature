Feature: Removing an application

	Scenario: Retrieve a Not Found response 
		Given one hydra server that is running without an application configuration file (apps.json)
		When the client requests DELETE /apps/app3
    Then the response should be 404 Not Found

  Scenario: Remove a concrete application succesfully
    Given one hydra server that is running with an application configuration file (apps.json) containing only one application
    When the client requests DELETE /apps/app1
    Then the response should be 200 OK
    When the client requests GET /apps/app1
    Then the response should be 404 Not Found