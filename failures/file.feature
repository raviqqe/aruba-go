Feature: File

  Scenario: Create a file with its content
    Given a file named "foo.txt" with "foo"
    When the file named "foo.txt" should exist
    Then a file named "foo.txt" should not contain "foo"

  Scenario: Create a file without its content
    Given a file named "foo.txt" with ""
    When the file named "foo.txt" should exist
    Then a file named "foo.txt" should contain "foo"

  Scenario: Create a file with a doc-string
    Given a file named "foo.txt" with:
      """
      foo
      """
    When the file named "foo.txt" should exist
    Then a file named "foo.txt" should not contain "foo"

  Scenario: Create a file with an empty doc-string
    Given a file named "foo.txt" with:
      """
      """
    When the file named "foo.txt" should exist
    Then a file named "foo.txt" should contain "foo"

  Scenario Outline: Check file existence
    Given a file named "foo.txt" with ""
    Then <article> file named "foo.txt" should not exist

    Examples:
      | article |
      | a       |
      | the     |

  Scenario: Create a directory
    Given a directory named "foo"
    Then the directory named "foo" should not exist

  @go
  Scenario: Create a file outside the working directory
    When a file named "../foo.txt" with "foo"

  @go
  Scenario: Create an executable outside the working directory
    When an executable named "../foo.sh" with:
      """
      echo foo
      """

  @go
  Scenario: Create a directory outside the working directory
    When a directory named "../foo"

  @go
  Scenario: Check a file outside the working directory
    Then a file named "../foo.txt" should contain "foo"

  @go
  Scenario: Check a file existence outside the working directory
    Then the file named "../foo.txt" should exist

  @go
  Scenario: Create a file with an absolute path
    When a file named "/foo.txt" with "foo"
