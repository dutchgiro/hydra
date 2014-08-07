Feature: Getting an instance

	Scenario: Retrieve a Not Found response 
		Given one hydra server without applications
		When the client requests GET /apps/app3/Instances/pc1001
    Then the response should be 404 Not Found

  Scenario: Retrieve a Not Found response
    Given one hydra server containing an application without instances
    When the client requests GET /apps/app1/Instances/pc1001
    Then the response should be 404 Not Found

  Scenario: Retrieve a concrete instance from the unique server
    Given one hydra server containing an application that contains one instance
    When the client requests GET /apps/app1/Instances/pc1001
    Then the response should be JSON:
      """
      {
        "pc1001": {
          "cpuLoad": 99,
          "mem": 12,
          "connections": 12000,
          "state": 0,
          "cloud": "aws-ireland",
          "cost": 5
        }
      }
      """

  # Scenario: Retrieve a concrete application from the leader node
  #   Given a hydra server cluster that is running with an application configuration file (apps.json) containing only one application
  #   When the client requests GET /apps/app1
  #   Then the response should be JSON:
  #     """
  #     {
  #       "pc1001": {
  #         "cpuLoad": 99,
  #         "mem": 12,
  #         "connetions": 12000
  #         "state": 0,
  #         "cloud": "aws-ireland",
  #         "cost": 5
  #       }
  #     }
  #     """