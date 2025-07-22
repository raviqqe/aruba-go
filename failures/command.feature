Feature: Command

  Scenario: Run a command successfully
    When I successfully run `false`

  Scenario: Check an exit status
    When I run `true`
    Then the exit status should not be 0

  Scenario: Check an exit status of 1
    When I run `false`
    Then the exit status should be 0

  Scenario: Check a non-zero exit status
    When I run `false`
    Then the exit status should not be 0

  Scenario: Run a command interactively
    When I run `echo`
    Then the exit status should be 0

  Scenario: Pipe in a file
    Given a file named "foo.txt" with:
      """
      foo
      """
    When I run `cat` interactively
    And I pipe in the file named "foo.txt"
    Then the exit status should be 0
    And the stdout should contain exactly "foo"

  Scenario: Pipe in a file without named
    Given a file named "foo.txt" with:
      """
      foo
      """
    When I run `cat` interactively
    And I pipe in the file "foo.txt"
    Then the exit status should be 0
    And the stdout should contain exactly "foo"
