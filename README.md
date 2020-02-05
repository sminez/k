k - Answers to questions you've already asked
=============================================

### What is k?
k is a tool for managing snippets, links, answers, names... you get the idea.
It should make it relatively simple to quickly search through all of your
snippets based on title, tag and namespacing through the filenames that you keep
them in.

It is essentially a wrapper that drives using fzf (so you'll need to have that
installed) which is a little more user friendly than simply grepping text files
on your filesystem. I'm a programmer so for me I want to be able to look at
text, code and links: so that's what k supports. I may end up adding some more
complicated actions around this in the future but for now, it simply hunts down
that bit of information that you are sure you knew at some point but can't quite
put your finger on.

k is named for my amazing wife (Katie) who has spent far too long doing this
sort of thing for me herself (remembering things I mean: not coding up utility
scripts). It's also kind of an in joke we have around me answering 'Que?' when
am asked a simple question that I'm drawing a blank on.

For more information, please see the rest of this README file or find me
somewhere online. If there is a user named 'sminez' there's a decent chance that
it's me.


### How to get started
You will also need to install the `fzf` program (available via your package
manager of choice) in order for `k` to run.

By default, k will look in `~/.helpfiles` for help snippets to display so you
can either clone this repo to that location or add the clone destination to the
`$HELPFILE_PATH` environment variable (same semantics as `$PATH` - `:` delimited
absolute paths) and k will scan all directories it finds on that path.

Simply run `k` and then fuzzy match through your snippets to select one. By
default the selected snippet will pretty print to the terminal with ANSI colors
but you can instead send the output to you system clipboard with the `-clip`
flag if desired. In this case, the ANSI color escape codes will be dropped.


### Snippet syntax
Snippet files are simple plain-text files that follow a simple convention for
deliminating individual snippets and marking lines with what type of content
they provide.
```
# This is a tile or description for a snippet. It's what you will see when searching
? these,are,comma,delimited,tags,you'll,see,them,as,well

% https://this.is.a.url.com

> This is a comment or note. You can have as many of these in an individual
> snippet as you like and they can be mixed and matched with code blocks and
> URLs. Note that each line can only be marked as a single content type however
> (so no nesting URLs in comments I'm afraid)

$ def code_block(example):
$     print("Regardless of language, code blocks are marked with a shell style $")
--
```

### Todo
* Copy only code selection to clipboard
* Pre-filter on certain tags?
* Show all tags
* Pretty print code snippets with bat?
  * Could also allow marking the language at the head of the snippet?
* Also allow for pulling snippets over the network?...


### To write up
* https://ablagoev.github.io/linux/adventures/commands/2017/02/19/adventures-in-usr-bin.html
* https://danielmiessler.com/blog/collection-of-less-commonly-used-unix-commands/#gs.xW4fUl4
* https://jvns.ca/blog/2017/03/26/bash-quirks/
* https://kkovacs.eu/cool-but-obscure-unix-tools
