Feature: Environment

  Scenario: Set an environment variable
    Given I set the environment variable "foo" to "bar"
    When I run the following script:
      """sh
      echo $foo
      """
    Then the stdout should contain exactly "bar"
