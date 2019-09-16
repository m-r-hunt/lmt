# lmt - Max Hunt ver.

This is a hard fork of the original lmt by driusan/Dave MacFarlane, which is to say I'm changing the program in incompatible ways with no intention of reconciliation.

I like lmt a lot, but I have an issue with how code blocks are tagged. There are two problems:

1) It's in a grey area of the standard. Although the GFM spec doesn't specify format or interpretation of the "info string" after a backtick code block header, the only common use is for specifying the language. Some parsers may fall down on more info given (or have some other interpretation). This has caused me problems.

2) Code blocks lack headers. See the original README.md defining lmt for an example. The code blocks lack any indication of where they go and it's hard to follow which block is being defined, especially when skimming the text.

For this reason I'm changing to use H6 headers (i.e. ###### lines) to define code block tags. This idea is swiped from knot, another markdown based literate programming tool. I don't want to use knot for other, unrelated reasons, which is why I'm modding lmt to fit my needs better.

The original README.md (which is a key part of lmt's source) is moved to README_orig.md. My new modification are in h6_tags.md. I've also fixed it for Windows by tidying up treatment of line endings to allow for \r\n style. This has gone into the original sources.

```shell
lmt README_orig.md WhitespacePreservation.md SubdirectoryFiles.md LineNumbers.md h6_tags.md
```

# TODO

I have a few things I'd like to do with lmt going forward.

* Improved warnings. I'd like to warn on unused code sections and on a file which produces no code sections.
* Rewrite prose into one coherent article. The current form of README + patch files make maintainence and reading a little tricky. I think rewriting the prose without changing the code would be an improvement.
