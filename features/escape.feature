@go
Feature: Character escape

  Scenario: Create a file with a backslash
    Given a file named "foo.txt" with:
      """
      a\b
      """
    When I successfully run `cat foo.txt`
    Then the stdout should contain exactly "a\\b"

  Scenario: Create a file with a double quote
    Given a file named "foo.txt" with:
      """
      a"b
      """
    When I successfully run `cat foo.txt`
    Then the stdout should contain exactly "a\"b"

  Scenario: Create a file with an escaped double quote
    Given a file named "foo.txt" with:
      """
      a\"b
      """
    When I successfully run `cat foo.txt`
    Then the stdout should not contain "a\\\"b"

  Scenario Outline: Escape a normal character
    Given a file named "foo.txt" with:
      """
      <value>
      """
    When I successfully run `cat foo.txt`
    Then the stdout should contain exactly "<value>"

    Examples:
      | value |
      | \\a   |

  Scenario: Check stdout with a blank character
    When I successfully run `echo \\\\\\\\`
    Then the stdout should contain exactly "\\\\"

  Scenario: Create a file with an escaped newline
    Given a file named "foo.py" with:
      """python
      print("foo\nbar")
      """
    When I successfully run `python3 foo.py`
    Then the stdout should contain exactly "foo\nbar"

  Scenario Outline: Create a file with an escaped example value
    Given a file named "foo.py" with:
      """python
      print("<value>")
      """
    When I successfully run `python3 foo.py`
    Then the stdout should contain exactly "<value>"
    And the stdout should contain exactly "foo\nbar"

    Examples:
      | value     |
      | foo\\nbar |

  Scenario: Create a file with many backslashes
    Given a file named "foo.py" with:
      """python
      print("\\\\\\\\")
      """
    When I successfully run `python3 foo.py`
    Then the stdout should contain "\\\\\\\\"

  Scenario Outline: Compare special characters in examples
    Given a file named "foo.py" with:
      """python
      print("<value>", end="")
      """
    When I successfully run `python3 foo.py`
    Then the stdout should contain exactly "<value>"

    Examples:
      | value |
      | \\n   |
      | \\t   |
      | \\"   |

  Scenario Outline: Escape special characters in examples
    Given a file named "foo.py" with:
      """python
      print("{!r}".format("<value>"))
      """
    When I successfully run `python3 foo.py`
    Then the stdout should contain exactly "'<value>'"

    Examples:
      | value     |
      | \\n       |
      | \\t       |
      | \\r       |
      | \\n\\t\\r |

  Scenario Outline: Compare asymmetric escapes in examples
    Given a file named "foo.py" with:
      """python
      print("{!r}".format(<input>).replace("'", '"'))
      """
    When I successfully run `python3 foo.py`
    Then the stdout should contain exactly "<output>"

    Examples:
      | input | output    |
      | "foo" | \\"foo\\" |

  Scenario: Match a file content with itself
    Given a file named "foo.txt" with:
      """
      a
      b
      """
    When I successfully run `cat foo.txt`
    Then the stdout should contain exactly "a\nb"
    And the stdout should not contain exactly "a\\nb"
    And the stdout should contain exactly:
      """
      a
      b
      """

  Scenario: Check a file to contain a string with a newline
    When a file named "foo.txt" with:
      """foo
      a
      b
      """
    Then a file named "foo.txt" should not contain "a\nb"

  Scenario: Check a file to contain an exact string with surrounding spaces
    When a file named "foo.txt" with:
      """foo

      a

      """
    Then a file named "foo.txt" should contain exactly:
      """
      a
      """
