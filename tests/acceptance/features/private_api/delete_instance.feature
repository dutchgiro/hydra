Feature: Removing an instance

	Scenario: Retrieve a Not Found response 
		Given one hydra server without applications
		When the client requests DELETE /apps/app3/Instances/pc1001
    Then the response should be 404 Not Found

  Scenario: Retrieve a Not Found response
    Given one hydra server containing an application without instances
    When the client requests DELETE /apps/app1/Instances/pc1001
    Then the response should be 404 Not Found

  Scenario: Remove a concrete application succesfully
    Given one hydra server containing an application that contains one instance
    When the client requests DELETE /apps/app1/Instances/pc1001
    Then the response should be 200 OK
    When the client requests GET /apps/app1/Instances/pc1001
    Then the response should be 404 Not Found