Feature: File

  Scenario: Create a file
    Given a file named "foo.txt" with "foo"
    When I successfully run `test -r foo.txt`
    Then a file named "foo.txt" should contain "foo"

  Scenario: Create a file with a doc-string
    Given a file named "foo.txt" with:
      """
      foo
      """
    When I successfully run `test -r foo.txt`
    Then a file named "foo.txt" should contain "foo"

  Scenario: Create a file with a content type
    Given a file named "foo.txt" with:
      """foo
      foo
      """
    When I successfully run `cat foo.txt`
    Then the stdout should contain exactly "foo"

  Scenario: Create a file in a directory
    Given a file named "foo/bar.txt" with "foo"
    When I successfully run `cat foo/bar.txt`
    Then the stdout should contain exactly "foo"

  Scenario Outline: Check file existence
    Given a file named "foo.txt" with:
      """
      foo
      """
    Then <article> file named "foo.txt" should exist
    And <article> file "foo.txt" should exist
    And <article> file named "bar.txt" should not exist
    And <article> file "bar.txt" should not exist

    Examples:
      | article |
      | a       |
      | the     |

  Scenario Outline: Remove a file and a directory
    Given a file named "foo.txt" with "foo"
    And a directory named "bar"
    When <article> file named "foo.txt" does not exist
    And <article> directory "bar" does not exist
    Then a file named "foo.txt" should not exist
    And a directory named "bar" should not exist

    Examples:
      | article |
      | a       |
      | the     |

  Scenario: Remove a directory with a file
    Given a file named "foo/bar.txt" with "foo"
    When a directory named "foo" does not exist
    Then a directory named "foo" should not exist

  Scenario: Remove a missing file
    When a file named "foo.txt" does not exist
    Then a file named "foo.txt" should not exist

  Scenario: Create a directory
    Given a directory named "foo"
    Then the directory named "foo" should exist
    And the directory named "bar" should not exist

  Scenario: Create an executable
    Given an executable named "foo.sh" with:
      """
      #!/bin/sh

      echo foo
      """
    When I successfully run `sh -c ./foo.sh`
    Then the stdout should contain exactly "foo"

  Scenario: Create a file with a mode
    Given a file named "foo.sh" with mode "0755" and with:
      """
      #!/bin/sh

      echo foo
      """
    When I successfully run `sh -c ./foo.sh`
    Then the stdout should contain exactly "foo"

  Scenario: Check a file with a definite article
    Given a file named "foo.txt" with:
      """
      foo
      """
    Then the file "foo.txt" should contain "foo"
    And the file named "foo.txt" should contain exactly:
      """
      foo
      """

  Rule: Remove files

    Scenario Outline: Remove a file and a directory
      Given a file named "foo.txt" with "foo"
      And a directory named "bar"
      When I remove <article> file named "foo.txt"
      And I remove <article> directory "bar"
      Then a file named "foo.txt" should not exist
      And a directory named "bar" should not exist

      Examples:
        | article |
        | a       |
        | the     |

    Scenario: Remove a directory with a file
      Given a file named "foo/bar.txt" with "foo"
      When I remove the directory "foo"
      Then a directory named "foo" should not exist

    Scenario: Remove a missing file with full force
      When I remove a file named "foo.txt" with full force
      Then a file named "foo.txt" should not exist

  Rule: Contain strings

    Scenario: Check a file to contain a string
      When a file named "foo.txt" with:
        """
        foo
        """
      Then a file named "foo.txt" should contain "foo"

    Scenario: Check a file not to contain a string
      When a file named "foo.txt" with:
        """
        foo
        """
      Then a file named "foo.txt" should not contain "bar"

  Rule: Contain doc-strings

    Scenario: Check a file to contain a string
      When a file named "foo.txt" with:
        """
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
        """
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
        """
        a

        """
      Then a file named "foo.txt" should contain exactly:
        """
        a
        """

    Scenario: Check a file to contain an exact string with surrounding spaces
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
        """
        a
        b
        """
      Then a file named "foo.txt" should not contain:
        """
        a
        c
        """

    Scenario: Check a file not to end with a newline
      When a file named "foo.txt" with:
        """
        a
        """
      Then a file named "foo.txt" should contain "a"
      And I successfully run `wc -c foo.txt`
      And the stdout should contain "1 foo.txt"

    Scenario: Check a file to keep trailing blank lines
      When a file named "foo.txt" with:
        """
        a


        """
      Then a file named "foo.txt" should contain "a"
      And I successfully run `wc -c foo.txt`
      And the stdout should contain "3 foo.txt"

    Scenario: Check a file to contain leading spaces
      When a file named "foo.txt" with:
        """
          a
        """
      Then a file named "foo.txt" should contain "  a"
