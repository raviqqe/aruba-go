Feature: File

  Scenario: Create a file
    Given a file named "foo.txt" with:
      """
      foo
      """
    When I successfully run `test -r foo.txt`

  Scenario: Create a file with a content type
    Given a file named "foo.txt" with:
      """foo
      foo
      """
    When I successfully run `cat foo.txt`
    Then the stdout should contain exactly "foo"

  Scenario: Check a file to contain a string
    When a file named "foo.txt" with:
      """foo
      foo
      """
    Then a file named "foo.txt" should contain "foo"

  Scenario: Check a file not to contain a string
    When a file named "foo.txt" with:
      """foo
      foo
      """
    Then a file named "foo.txt" should not contain "bar"

  Scenario: Check a file not to contain a string with a newline
    When a file named "foo.txt" with:
      """foo
      a
      b
      """
    Then a file named "foo.txt" should not contain "a\\nb"
