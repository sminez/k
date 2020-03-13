k - Answers to questions you've already asked
=============================================

![demo](https://github.com/sminez/k/blob/master/k-1.gif)

### What is k?
k is a tool for managing snippets, links, answers, names... you get the idea.
It should make it relatively simple to quickly search through all of your
snippets based on title, tag and name-spacing through the filenames that you keep
them in.

It is essentially a wrapper that drives using fzf (so you'll need to have that
installed) which is a little more user friendly than simply grepping text files
on your file system. I'm a programmer so for me I want to be able to look at
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

Simply run `k` and then start typing to fuzzy match through your snippets to
select one. When you hit enter your chosen snippet will be printed to the
terminal.


### Snippet syntax
Snippet files are simple plain-text files that follow a simple convention for
delimiting individual snippets and marking lines with what type of content
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
* Copy to clipboard
* Pre-filter on certain tags
* Show all tags

### Non-goals / rejected ideas
* Allow for pulling snippets over the network
  * I've tried a couple of things to get this to work and it isn't proving to be
  easy to do in an efficient manner. For a git repo for example, you can just
  clone the thing (and then optionally, auto fetch on run) but that adds a
  pretty big start up cost to running `k` in the general case. You can also use
  the unauthenticated Github API for listing directory contents and then pulling
  the raw files straight into memory on the fly but that is also pretty slow,
  even if you spin out a goro for each file being pulled over the network.
  * An alternative would be to add some sort of flag for cloning a remote set of
  helpfiles and then auto-updating all local repos to ensure we're up to date
  and at that point we're just writing a full blown git client... So, in the
  spirit of minimalism, suckless and "I don't have time for this", lets just
  settle on "clone the damn repo!" as the accepted solution.
* Managing snippets from `k` itself.
  * Yes I could write a TUI / wizard for adding a new snippet to your local
  helpfiles, but that's what a text editor is for. `k` just wraps `fzf` to make
  it easier to search through your snippets.


### To write up
* https://ablagoev.github.io/linux/adventures/commands/2017/02/19/adventures-in-usr-bin.html
* https://danielmiessler.com/blog/collection-of-less-commonly-used-unix-commands/#gs.xW4fUl4
* https://jvns.ca/blog/2017/03/26/bash-quirks/
* https://kkovacs.eu/cool-but-obscure-unix-tools
