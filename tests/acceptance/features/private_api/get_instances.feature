Feature: Getting list of instances

	Scenario: Retrieve an empty application list
		Given one hydra server that is running without an application configuration file (apps.json)
		When the client requests GET /apps/app1/Instances
    Then the response should be 404 Not Found

  Scenario: Retrieve an empty instance list
    Given one hydra server containing an application without instances
    When the client requests GET /apps/app1/Instances
    Then the response should be JSON:
      """
      []
      """

  Scenario: Retrieve an instance list
    Given one hydra server containing an application that contains one instance
    When the client requests GET /apps/app1/Instances
    Then the response should be JSON:
      """
      [{
        "pc1001": {
          "cpuLoad": 99,
          "mem": 12,
          "connections": 12000,
          "state": 0,
          "cloud": "aws-ireland",
          "cost": 5
        }
      }]
      """

  Scenario: Retrieve an instance list
    Given one hydra server containing an application that contains two instances
    When the client requests GET /apps/app1/Instances
    Then the response should be JSON:
      """
      [{
        "pc1001": {
          "cpuLoad": 99,
          "mem": 12,
          "connections": 12000,
          "state": 0,
          "cloud": "aws-ireland",
          "cost": 5
        }
      }, {
        "pc1002": {
          "cpuLoad": 22,
          "mem": 43,
          "connections": 400,
          "state": 1,
          "cloud": "aws-oregon",
          "cost": 9
        }
      }]
      """

  # Scenario: Retrieve an application list
  #   Given a hydra server cluster that is running with an application configuration file (apps.json) containing only one application
  #   When the client requests GET /apps
  #   Then the response should be JSON:
  #     """
  #     [{
  #         "app1": {
  #             "Balancers": [
  #                 {
  #                     "worker": "MapAndSort",
  #                     "mapAttr": "cloud",
  #                     "mapSort": ["google", "amazon", "azure"]
  #                 },
  #                 {
  #                     "worker": "SortByNumber",
  #                     "sortAttr": "cpuLoad",
  #                     "order": 1
  #                 }
  #             ]
  #         }
  #     }]
  #     """