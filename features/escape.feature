Feature: Character escape
  Scenario: Create a file with an escaped backslash
    Given a file named "foo.txt" with:
      """
      a\\b
      """
    When I successfully run `cat foo.txt`
    Then the stdout should contain exactly "a\\b"

  Scenario: Create a file with an escaped double quote
    Given a file named "foo.txt" with:
      """
      a\"b
      """
    When I successfully run `cat foo.txt`
    Then the stdout should contain exactly "a\"b"

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

    Examples:
      | value     |
      | foo\\nbar |

  Scenario: Create a file with many backslashes
    Given a file named "foo.py" with:
      """python
      print("\\\\\\\\")
      """
    When I successfully run `python3 foo.py`
    Then the stdout should contain "\\\\"

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

  Scenario Outline: Compar asymmetric escapes in examples
    Given a file named "foo.py" with:
      """python
      print("{!r}".format(<input>).replace("'", '"'))
      """
    When I successfully run `python3 foo.py`
    Then the stdout should contain exactly "<output>"

    Examples:
      | input | output    |
      | "foo" | \\"foo\\" |
