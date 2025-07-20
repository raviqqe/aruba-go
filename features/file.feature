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

  Rule: Contain strings

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

  Rule: Contain doc-strings

    Scenario: Check a file to contain a string
      When a file named "foo.txt" with:
        """foo
        a
        b
        c
        d
        """
      Then a file named "foo.txt" should contain:
        """
        b
        c
        """

    Scenario: Check a file to contain an exact string
      When a file named "foo.txt" with:
        """foo
        a
        b
        """
      Then a file named "foo.txt" should contain exactly:
        """
        a
        b
        """

    Scenario: Check a file to contain an exact string with trailing spaces
      When a file named "foo.txt" with:
        """foo
        a

        """
      Then a file named "foo.txt" should contain exactly:
        """
        a
        """

    Scenario: Check a file not to contain an exact string with surrounding spaces
      When a file named "foo.txt" with:
        """foo

        a

        """
      Then a file named "foo.txt" should not contain exactly:
        """
        a
        """

    Scenario: Check a file not to contain a string
      When a file named "foo.txt" with:
        """foo
        a
        b
        """
      Then a file named "foo.txt" should not contain:
        """
        a
        c
        """
