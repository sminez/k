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
```

### Todo
* show all known tags
* pretty print code snippets with bat?
  * could also allow marking the language at the head of the snippet?
* helper for creating new snippets? (probably not needed really)
* default snippets to include with the repo
* a way of keeping personal snippets in their own gitignored directory
