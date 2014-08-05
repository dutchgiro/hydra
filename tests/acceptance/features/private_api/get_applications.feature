Feature: Getting list of applications

	Scenario: Retrieve an empty application list
		Given one hydra server that is running without an application configuration file (apps.json)
		When the client requests GET /apps
    Then the response should be 404 Not Found

  Scenario: Retrieve an empty application list
    Given one hydra server that is running with an application configuration file (apps.json) containing an empty array of applications
    When the client requests GET /apps
    Then the response should be 404 Not Found

  Scenario: Retrieve an application list
    Given one hydra server that is running with an application configuration file (apps.json) containing only one application
    When the client requests GET /apps
    Then the response should be JSON:
      """
      [{
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
                      "order": 1,
                      "reverse": true,
                      "isValid": false,
                      "decimal": 55.19,
                      "sortAttr": "cpuLoad",
                      "worker": "SortByNumber"
                  }
              ]
          }
      }]
      """

  Scenario: Retrieve an application list
    Given one hydra server that is running with an application configuration file (apps.json) containing two applications
    When the client requests GET /apps
    Then the response should be JSON:
      """
      [{
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
      }, {
          "app2": {
              "Balancers": [
                  {
                      "worker": "MapByLimit",
                      "limitAttr": "limit",
                      "limitValue": 50,
                      "mapSort": "reverse"
                  },
                  {
                      "worker": "RoundRobin",
                      "simple": "OK"
                  }
              ]
          }
      }]
      """

  Scenario: Retrieve an application list
    Given a hydra server cluster that is running with an application configuration file (apps.json) containing only one application
    When the client requests GET /apps
    Then the response should be JSON:
      """
      [{
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
      }]
      """