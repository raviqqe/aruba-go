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
