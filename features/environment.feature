Feature: Environment

  Scenario: Set an environment variable
    Given I set the environment variable "foo" to "bar"
    When I run the following script:
      """sh
      echo $foo
      """
    Then the stdout should contain exactly "bar"

  Scenario: Change a directory
    Given a directory named "foo"
    And a file named "foo/bar.txt" with "foo"
    When I cd to "foo"
    And I successfully run `cat bar.txt`
    Then the stdout should contain exactly "foo"
