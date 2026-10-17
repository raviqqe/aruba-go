Feature: Command

  Scenario: Run a command
    When I run `echo hello world`

  Scenario: Run a command successfully
    When I successfully run `echo hello world`

  Scenario: Check an exit status
    When I run `true`
    Then the exit status should be 0

  Scenario: Check an exit status of 1
    When I run `false`
    Then the exit status should be 1

  Scenario: Check a non-zero exit status
    When I run `false`
    Then the exit status should not be 0

  Scenario: Run a command interactively
    When I run `echo` interactively
    Then the exit status should be 0

  Scenario: Run a command in the background
    When I run `echo foo` in the background
    Then the stdout should contain exactly "foo"

  Scenario: Wait for a command to start up
    When I wait 1 second for the command to start up
    And I run `sh -c 'sleep 0.1; echo foo > foo.txt'` in the background
    Then the file "foo.txt" should exist

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

  Scenario: Run a command with a single-quoted argument
    When I successfully run `sh -c 'echo foo bar'`
    Then the stdout should contain exactly "foo bar"

  Scenario: Run a command with a double-quoted argument
    When I successfully run `sh -c "echo foo bar"`
    Then the stdout should contain exactly "foo bar"

  Scenario: Run the following commands
    When I run the following commands:
      """sh
      echo foo
      """
    Then the stdout should contain exactly "foo"
