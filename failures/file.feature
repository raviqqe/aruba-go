Feature: File

  Scenario: Create a file
    Given a file named "foo.txt" with "foo"
    When I successfully run `test -r foo.txt`
    Then a file named "foo.txt" should not contain "foo"

  Scenario: Create a file with a doc-string
    Given a file named "foo.txt" with:
      """
      foo
      """
    When I successfully run `test -r foo.txt`
    Then a file named "foo.txt" should not contain "foo"

  Scenario: Create a file with a content type
    Given a file named "foo.txt" with:
      """foo
      foo
      """
    When I successfully run `cat foo.txt`
    Then the stdout should not contain "foo"

  Scenario: Check file existence
    Given a file named "foo.txt" with ""
    Then <article> file named "foo.txt" should not exist

    Examples:
      | article |
      | a       |
      | the     |
  # Scenario: Create a directory
  #   Given a directory named "foo"
  #   Then the directory named "foo" should exist
  #   And the directory named "bar" should not exist
  #
  # Rule: Contain strings
  #
  #   Scenario: Check a file to contain a string
  #     When a file named "foo.txt" with:
  #       """
  #       foo
  #       """
  #     Then a file named "foo.txt" should contain "foo"
  #
  #   Scenario: Check a file not to contain a string
  #     When a file named "foo.txt" with:
  #       """
  #       foo
  #       """
  #     Then a file named "foo.txt" should not contain "bar"
  #
  # Rule: Contain doc-strings
  #
  #   Scenario: Check a file to contain a string
  #     When a file named "foo.txt" with:
  #       """
  #       a
  #       b
  #       c
  #       d
  #       """
  #     Then a file named "foo.txt" should contain:
  #       """
  #       b
  #       c
  #       """
  #
  #   Scenario: Check a file to contain an exact string
  #     When a file named "foo.txt" with:
  #       """
  #       a
  #       b
  #       """
  #     Then a file named "foo.txt" should contain exactly:
  #       """
  #       a
  #       b
  #       """
  #
  #   Scenario: Check a file to contain an exact string with trailing spaces
  #     When a file named "foo.txt" with:
  #       """
  #       a
  #
  #       """
  #     Then a file named "foo.txt" should contain exactly:
  #       """
  #       a
  #       """
  #
  #   Scenario: Check a file not to contain a string
  #     When a file named "foo.txt" with:
  #       """
  #       a
  #       b
  #       """
  #     Then a file named "foo.txt" should not contain:
  #       """
  #       a
  #       c
  #       """
  #
  #   Scenario: Check a file to contain a newline
  #     When a file named "foo.txt" with:
  #       """
  #       a
  #       """
  #     Then a file named "foo.txt" should contain "a\n"
  #
  #   Scenario: Check a file to contain two newlines
  #     When a file named "foo.txt" with:
  #       """
  #       a
  #
  #
  #       """
  #     Then a file named "foo.txt" should contain "a\n"
  #     And a file named "foo.txt" should not contain "a\n\n"
  #
  #   Scenario: Check a file to contain two newlines
  #     When a file named "foo.txt" with:
  #       """
  #         a
  #       """
  #     Then a file named "foo.txt" should contain "  a"
