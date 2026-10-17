Feature: Character escape

  Scenario: Create a file with a backslash
    Given a file named "foo.txt" with:
      """
      a\b
      """
    When I successfully run `cat foo.txt`
    Then the stdout should contain exactly "a\b"

  Scenario: Create a file with a double quote
    Given a file named "foo.txt" with:
      """
      a"b
      """
    When I successfully run `cat foo.txt`
    Then the stdout should contain exactly "a\"b"

  Scenario: Create a file with a single quote
    Given a file named "foo.txt" with:
      """
      a'b
      """
    When I successfully run `cat foo.txt`
    Then the stdout should contain exactly "a\'b"

  Scenario: Create a file with an escaped double quote
    Given a file named "foo.txt" with:
      """
      a\"b
      """
    When I successfully run `cat foo.txt`
    Then the stdout should contain:
      """
      a\"b
      """
    And the stdout should not contain "a\"b"

  Scenario: Unescape a backslash in a command
    When I successfully run `echo \\`
    Then the stdout should contain exactly:
      """
      \
      """

  Scenario: Create a file with an escaped newline
    Given a file named "foo.py" with:
      """python
      print("foo\nbar")
      """
    When I successfully run `python3 foo.py`
    Then the stdout should contain exactly:
      """
      foo
      bar
      """

  Scenario Outline: Create a file with an escaped example value
    Given a file named "foo.py" with:
      """python
      print("<value>")
      """
    When I successfully run `python3 foo.py`
    Then the stdout should contain exactly:
      """
      foo
      bar
      """

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

  @go
  Scenario Outline: Compare escape sequences in examples
    Given a file named "foo.py" with:
      """python
      print(r"<value>", end="")
      """
    When I successfully run `python3 foo.py`
    Then the stdout should contain exactly "<value>"

    Examples:
      | value |
      | \\a   |
      | \\n   |
      | \\t   |
      | \\r   |

  @go
  Scenario Outline: Escape special characters in examples
    Given a file named "foo.py" with:
      """python
      print("{!r}".format("<value>"))
      """
    When I successfully run `python3 foo.py`
    Then the stdout should contain exactly "'<value>'"

    Examples:
      | value |
      | \\n   |
      | \\t   |
      | \\r   |

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
    Then the stdout should not contain exactly "a\nb"
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
    Then a file named "foo.txt" should contain:
      """
      a
      b
      """
    And a file named "foo.txt" should not contain "a\nb"

  @go
  Scenario: Check a file to contain escaped quotes
    When a file named "foo.txt" with:
      """
      a"b'c
      """
    Then a file named "foo.txt" should contain "a\"b\'c"

  Scenario: Pipe in a file with an escaped single quote
    Given a file named "foo'.txt" with:
      """
      foo
      """
    When I run `cat` interactively
    And I pipe in the file "foo\'.txt"
    Then the stdout should contain exactly "foo"

  @go
  Scenario: Create files with escaped quotes in their names
    Given a file named "a\"b\'c" with "foo"
    And a file named "d\"e\'f" with:
      """
      bar
      """
    And a file named "g\"h\'i" with mode "0600" and with:
      """
      baz
      """
    And an executable named "j\"k\'l" with:
      """
      #!/bin/sh
      """
    And a directory named "m\"n\'o"
    When I successfully run `ls`
    Then the stdout should contain exactly:
      """
      a"b'c
      d"e'f
      g"h'i
      j"k'l
      m"n'o
      """
    And the file "a\"b\'c" should exist
    And the file "a\"b\'c" should contain "foo"
    And the file "d\"e\'f" should contain:
      """
      bar
      """
    And the directory "m\"n\'o" should exist

  @go
  Scenario: Remove files with escaped quotes in their names
    Given a file named "a\"b\'c" with "foo"
    And a file named "d\"e\'f" with "bar"
    When I remove the file "a\"b\'c"
    And the file "d\"e\'f" does not exist
    And I successfully run `ls`
    Then the stdout should contain exactly ""

  @go
  Scenario: Change to a directory with escaped quotes in its name
    Given a directory named "a\"b\'c"
    When I cd to "a\"b\'c"
    And I successfully run `pwd`
    Then the stdout should contain "/a\"b'c"

  @go
  Scenario: Set an environment variable with escaped quotes
    Given I set the environment variable "FOO" to "a\"b"
    And I append "\'c" to the environment variable "FOO"
    When I successfully run `sh -c 'printf %s "$FOO"'`
    Then the stdout should contain exactly "a\"b'c"
