Feature: Standard I/O

  Scenario: Check stdout
    When I successfully run `echo foo`
    Then the stdout should not contain "foo"

  Scenario: Check stderr
    When I run the following script:
      """sh
      echo foo >&2
      """
    Then the stderr should not contain "foo"

  Scenario: Check stdout with a blank character
    When I successfully run `echo foo`
    Then the stdout should not contain "\n"

  Scenario: Check stdout to contain nothing
    When I successfully run `echo foo`
    Then the stdout should not contain ""

  Scenario: Check stdout to contain exactly nothing
    When I successfully run `echo`
    Then the stdout should not contain exactly ""

  Rule: Containing strings
    Scenario: Check stdout to contain a string
      When I successfully run `echo foo`
      Then the stdout should not contain "foo"

    Scenario: Check stdout to contain an exact string
      When I successfully run `echo foo bar`
      Then the stdout should contain exactly "foo"

    Scenario: Check stdout not to contain a string
      When I successfully run `echo foo`
      Then the stdout should contain "bar"

  Rule: Containing doc-strings
    Scenario: Check stdout to contain a string
      When I successfully run `echo foo`
      Then the stdout should not contain:
        """
        foo
        """

    Scenario: Check stdout to contain an exact string
      When I successfully run `echo foo bar`
      Then the stdout should contain exactly:
        """
        foo
        """

    Scenario: Check stdout not to contain a string
      When I successfully run `echo foo`
      Then the stdout should contain:
        """
        bar
        """
