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

  Scenario: Update a concrete application
    Given one hydra server that is running with an application configuration file (apps.json) containing only one application
    When the client requests POST /apps
      """
      {
        "app1": {
          "Balancers": [
              {
                  "worker": "MapAndSort",
                  "mapAttr": "cloud",
                  "mapSort": ["azure", "google", "amazon", "rackspace"]
              },
              {
                  "worker": "SortByNumber",
                  "sortAttr": "mem",
                  "order": 0,
                  "reverse": false,
                  "isValid": true,
                  "decimal": 22.73
              }
          ]
        }
      }
      """
    Then the response should be 200 OK
    When the client requests GET /apps/app1
    Then the response should be JSON:
      """
      {
        "app1": {
          "Balancers": [
              {
                  "worker": "MapAndSort",
                  "mapAttr": "cloud",
                  "mapSort": ["azure", "google", "amazon", "rackspace"]
              },
              {
                  "worker": "SortByNumber",
                  "sortAttr": "mem",
                  "order": 0,
                  "reverse": false,
                  "isValid": true,
                  "decimal": 22.73
              }
          ]
        }
      }
      """