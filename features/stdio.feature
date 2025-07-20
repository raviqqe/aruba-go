Feature: Standard I/O

  Scenario: Check stdout
    When I run `echo foo bar`
    Then the stdout should contain "foo"

  Scenario: Check exact stdout
    When I run `echo foo`
    Then the stdout should contain exactly "foo"

  Scenario: Check stderr
    When I run `rm foo`
    Then the stderr should contain "file"

  Scenario: Check stdout with a blank character
    When I successfully run `echo foo`
    Then the stdout should contain "\n"

  Scenario: Check stdout to contain nothing
    When I successfully run `echo foo`
    Then the stdout should contain ""

  Scenario: Check stdout to contain exactly nothing
    When I successfully run `echo`
    Then the stdout should contain exactly ""

  Rule: Containing strings

    Scenario: Check stdout to contain a string
      When I run `echo foo bar`
      Then the stdout should contain "foo"

    Scenario: Check stdout to contain an exact string
      When I run `echo foo`
      Then the stdout should contain exactly "foo"

    Scenario: Check stdout not to contain a string
      When I run `echo foo`
      Then the stdout should not contain "bar"

    Scenario: Check stdout not to contain an exact string
      When I run `echo foo`
      Then the stdout should contain exactly "bar"

  Rule: Containing doc-strings

    Scenario: Check stdout to contain a string
      When I successfully run `echo foo`
      Then the stdout should contain:
        """
        foo
        """

    Scenario: Check stdout to contain an exact string
      When I successfully run `echo foo`
      Then the stdout should contain exactly:
        """
        foo
        """

    Scenario: Check stdout not to contain a string
      When I successfully run `echo foo`
      Then the stdout should not contain:
        """
        bar
        """

    Scenario: Check stdout not to contain an exact string
      When I successfully run `echo foo`
      Then the stdout should not contain exactly:
        """
        bar
        """
