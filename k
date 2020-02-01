#!/usr/bin/env python3
"""
k - Answers to questions that you have already asked

k is a tool for managing snippets, links, answers, names... you get the idea.
It should make it relatively simple to quickly search through all of your
snippets based on title, tag and namespacing through the filenames that you
keep them in.

It is essentially a wrapper that drives using fzf (so you'll need to have that
installed) which is a little more user friendly than simply grepping text files
on your filesystem. I'm a programmer so for me I want to be able to look at
text, code and links: so that's what k supports. I may end up adding some more
complicated actions around this in the future but for now, it simply hunts down
that bit of information that you are sure you knew at some point but can't
quite put your finger on.

k is named for my amazing wife (Katie) who has spent far too long doing this
sort of thing for me herself (remembering things I mean: not coding up utility
scripts). It's also kind of an in joke we have around me answering 'Que?' when
am asked a simple question that I'm drawing a blank on.

For more information, please see the README file in the repo or find me
somewhere online. If there is a user named 'sminez' there's a decent chance
that it's me.

Enjoy!

- Innes Anderson-Morrison (sminez)
"""
import subprocess
import argparse
import sys
import os


PATH = os.path.normpath(os.path.join(os.path.abspath(__file__), os.pardir))
HELPFILE_DIRECTORY = os.environ.get("HELPFILE_DIR", os.path.join(PATH, "helpfiles"))
END_OF_SECTION_MARKER = "--"
COMMENT_MARKER = ">"
TITLE_MARKER = "#"
CODE_MARKER = "$"
URL_MARKER = "%"
TAG_MARKER = "?"
SEP = " | "


class Snippet:
    def __init__(self, fname, lines):
        self.file = fname
        self.title = None
        self.tags = None

        for line in lines:
            if line.startswith(TITLE_MARKER):
                self.title = line.strip(f"{TITLE_MARKER} ")
            elif line.startswith(TAG_MARKER):
                self.tags = line.strip(f"{TAG_MARKER} ")

        if self.title is None or self.tags is None:
            section = "\n".join(lines)
            print(f"{fname} contains an invalid section: {section}", file=sys.stderr)
            sys.exit(1)

    def fzf_format(self, file_len: int, title_len: int, tag_len: int):
        return f"{SEP}".join(
            [self.file.ljust(file_len), self.tags.ljust(tag_len), self.title.ljust(title_len)]
        )


def snippet_selections():
    def get_len(snippets, attr):
        return max(map(lambda s: len(getattr(s, attr)), snippets))

    helpfiles = sorted(os.listdir(HELPFILE_DIRECTORY), reverse=True)
    snippets = []

    for fname in helpfiles:
        _snippets, current_snippet = [], []

        with open(os.path.join(HELPFILE_DIRECTORY, fname), "r") as f:
            for line in f:
                line = line.strip()
                if line == END_OF_SECTION_MARKER:
                    _snippets.append(Snippet(fname, current_snippet))
                    current_snippet = []
                else:
                    current_snippet.append(line)

        snippets.extend(sorted(_snippets, key=lambda s: s.title, reverse=True))

    file_len, title_len, tag_len = map(lambda s: get_len(snippets, s), ["file", "title", "tags"])
    return "\n".join(s.fzf_format(file_len, title_len, tag_len) for s in snippets)


def run_interactive(cmd):
    proc = subprocess.run(cmd, stdin=sys.stdin, stdout=subprocess.PIPE, shell=True)
    return proc.stdout.decode("utf-8").strip()


def get_section_from_file(selection):
    file, _, title = selection.split(SEP)
    section = []

    with open(os.path.join(HELPFILE_DIRECTORY, file), "r") as f:
        while True:
            line = next(f)
            if title in line:
                section.append(line.strip())
                break
        for line in f:
            if line.strip() == END_OF_SECTION_MARKER:
                break
            section.append(line.strip())
    return section


def pretty_print(section):
    def pprint(color, line):
        print(f"\033[1;{color}m{line}\033[0m")

    for line in section:
        if line.startswith(TITLE_MARKER):
            pprint(33, line[2:])  # yellow
        elif line.startswith(URL_MARKER):
            pprint(34, line[2:])  # blue
        elif line.startswith(COMMENT_MARKER):
            pprint(37, line[2:])  # grey
        elif line.startswith(CODE_MARKER):
            print(line)
        elif line == "":
            print()


if __name__ == "__main__":
    parser = argparse.ArgumentParser("k", description="Answers for questions you've already asked")
    parser.add_argument("--pprint", default="", help=argparse.SUPPRESS)
    args = parser.parse_args()

    if args.pprint:
        pretty_print(get_section_from_file(args.pprint))
    else:
        selections = snippet_selections()
        selection = run_interactive(
            f'echo "{selections}" | fzf --preview "{PATH}/k ' + '--pprint {}"'
        )
        pretty_print(get_section_from_file(selection))
