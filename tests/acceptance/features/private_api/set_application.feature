Feature: Setting an application

	Scenario: Add a concrete application
    Given one hydra server that is running without an application configuration file (apps.json)
    When the client requests POST /apps
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
    Then the response should be 200 OK

  Scenario: Override a concrete application
    Given one hydra server that is running with an application configuration file (apps.json) containing only one application
    When the client requests POST /apps
      """
      {
          "app1": {
              "Balancers": [
                  {
                      "worker": "MapByLimit",
                      "limitAttr": "limit",
                      "limitValue": 50,
                      "mapSort": "reverse"
                  },
                  {
                      "worker": "RoundRobin"
                      "simple": "OK"
                  }
              ]
          }
      }
      """
    Then the response should be 200 OK