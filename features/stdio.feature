Feature: Standard I/O
  Scenario: Check stdout
    When I run `echo foo bar baz`
    Then the stdout should contain "bar"

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
