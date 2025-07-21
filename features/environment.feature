Feature: Environment

  Scenario: Set an environment variable
    Given I set an environment variable "foo" to "bar"
    When I successfully run `echo $foo`
    Then the stdout should contain exactly "bar"
