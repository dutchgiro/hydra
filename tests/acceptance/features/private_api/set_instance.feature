Feature: Setting an instance

  Scenario: Retrieve a Bad Request response
    Given one hydra server that is running without an application configuration file (apps.json)
    When the client requests POST /apps/app1/Instances
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
    Then the response should be 400 Bad Request

	Scenario: Add a concrete instance
    Given one hydra server containing an application without instances
    When the client requests POST /apps/app1/Instances
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
    Then the response should be 200 OK
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

  Scenario: Update a concrete instance
    Given one hydra server containing an application that contains one instance
    When the client requests POST /apps/app1/Instances
      """
      {
        "pc1001": {
          "cpuLoad": 11,
          "mem": 98,
          "connections": 24567,
          "state": 1,
          "cloud": "aws-oregon",
          "cost": 2
        }
      }
      """
    Then the response should be 200 OK
    When the client requests GET /apps/app1/Instances/pc1001
    Then the response should be JSON:
      """
      {
        "pc1001": {
          "cpuLoad": 11,
          "mem": 98,
          "connections": 24567,
          "state": 1,
          "cloud": "aws-oregon",
          "cost": 2
        }
      }
      """