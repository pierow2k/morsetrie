<!-- markdownlint-disable no-inline-html no-emphasis-as-heading -->
# morsetrie

![morsetrie Banner](./assets/morsetrie_banner-1200x400.png)
![Website](https://img.shields.io/website?url=https%3A%2F%2Fpkg.go.dev%2Fgithub.com%2Fpierow2k%2Fmorsetrie)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/pierow2k/morsetrie)
![License](https://img.shields.io/github/license/pierow2k/morsetrie)
![GitHub Tag](https://img.shields.io/github/v/tag/pierow2k/morsetrie)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/PROJECTID)](https://app.codacy.com/gh/pierow2k/morsetrie/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)

**morsetrie is a Go Package for Fast Morse Code Decoding**

<!-- TABLE OF CONTENTS -->
<details closed="closed">
<summary><h2 style="display: inline-block">Table of Contents</h2></summary>

- [Overview](#overview)
  - [What it does](#what-it-does)
  - [Why a trie?](#why-a-trie)
- [Usage](#usage)
- [Error handling behavior](#error-handling-behavior)
- [Contributing](#contributing)
- [License](#license)

</details>

## Overview

`morsetrie` is a small, fast Morse code decoder built around a compact trie
(prefix tree). It turns sequences of dots and dashes into text by walking a
pre-built decoding tree, making decoding efficient and predictable even for
long inputs.

The package ships with `MorseTable`, an International Morse code mapping
based on [ITU M.1677](https://www.itu.int/rec/R-REC-M.1677-1-200910-I/)
(letters, digits, and common punctuation), and provides helpers to build a
trie from any mapping table you supply.

### What it does

- Decodes `.` and `-` into runes using a trie-based lookup.
- Uses whitespace (`space`, `tab`, `newline`, `\r`) to delimit characters.
- Treats `/` as a word separator and emits a space in the decoded output.
- Represents unknown or invalid Morse sequences as `?` (rather than failing
mid-stream).

### Why a trie?

A trie is a natural fit for Morse: each dot/dash is a step down the tree.
This avoids repeatedly scanning a table or building strings for lookups.
Internally the implementation is **array-backed** and uses **int16 child
indices** to keep memory usage low while remaining cache-friendly.

![morsetrie diagram](./assets/morsetrie_diagram.png)

## Usage

Refer to the [package documentation on pkg.go.dev](https://pkg.go.dev/github.com/pierow2k/morsetrie)
and the [example_test.go](example_test.go) file for concrete examples.

Most users will build a trie once (often from `MorseTable`) and reuse it to
decode many messages.

- Build: `BuildTrie(morsetrie.MorseTable)`
- Decode: `trie.Decode("... --- ...")`

```go
import (
    "fmt"

    "github.com/pierow2k/morsetrie"
)

// Build the trie using the built-in morse code data from morsetrie's
// `MorseTable`.
trie, err := morsetrie.BuildTrie(morsetrie.MorseTable)
if err != nil {
  panic(err)
}

// Define morse code input.
morseCode := ".... . .-.. .-.. --- / .-- --- .-. .-.. -.."

// Call trie.Decode to decode the morse code data.
text, _ := trie.Decode(morseCode)

// Print the decoded text.
fmt.Println(text)
```

**Output**

```text
HELLO WORLD
```

## Error handling behavior

- If the input contains unsupported characters (anything other than `.`,
`-`, `/`, or whitespace), decoding fails with `ErrUnexpectedChar`.
- If the input contains a syntactically valid but unknown Morse sequence,
the decoder emits `?` for that symbol and continues.
- When building a trie, duplicate codes or invalid code strings are
rejected (`ErrDuplicate`, `ErrInvalidElement`), and extremely large tries
are prevented from overflowing the index type (`ErrTrieFull`).

## Contributing

We welcome contributions! Here's how you can help:

- Add a [GitHub Star](github.com/pierow2k/morsetrie) to the morsetrie project.
- Have an idea for a new feature or noticed something that isn’t working
quite right? [Open an issue](https://github.com/pierow2k/morsetrie/issues)
to let us know. Your feedback helps us keep morsetrie reliable
and feature-rich.
- **Submit a Pull Request**: If you’ve made improvements or fixed a bug,
we’d love to see your work.
[Submit a pull request](https://www.github.com/pierow2k/morsetrie/pulls)
and share your changes with the community.

We appreciate your support and contributions, which drive the continued
growth and success of morsetrie. Thank you for being part of the
journey!

## License

morsetrie is distributed under the MIT License. See the
[LICENSE](LICENSE) file for more details.
