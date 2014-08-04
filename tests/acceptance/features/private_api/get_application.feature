Feature: Getting an application

	Scenario: Retrieve a Not Found response 
		Given one hydra server that is running without an application configuration file (apps.json)
		When the client requests GET /apps/app3
    Then the response should be 404 Not Found

  Scenario: Retrieve a concrete application
    Given one hydra server that is running with an application configuration file (apps.json) containing only one application
    When the client requests GET /apps/app1
    Then the response should be JSON:
      """
      {
          "app1": {
              "Balancers": [
                  {
                      "worker": "MapAndSort",
                      "mapAttr": "cloud",
                      "mapSort": ["google", "amazon", "azure"]
                  },
                  {
                      "worker": "SortByNumber",
                      "sortAttr": "cpuLoad",
                      "order": 1
                  }
              ]
          }
      }
      """

  Scenario: Retrieve a concrete application
    Given a hydra server cluster that is running with an application configuration file (apps.json) containing only one application
    When the client requests GET /apps/app1
    Then the response should be JSON:
      """
      {
          "app1": {
              "Balancers": [
                  {
                      "worker": "MapAndSort",
                      "mapAttr": "cloud",
                      "mapSort": ["google", "amazon", "azure"]
                  },
                  {
                      "worker": "SortByNumber",
                      "sortAttr": "cpuLoad",
                      "order": 1
                  }
              ]
          }
      }
      """