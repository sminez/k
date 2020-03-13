#!/usr/bin/env zsh
# k is a tool for managing snippets, links, answers, names... you get the idea.
# It should make it relatively simple to quickly search through all of your
# snippets based on title, tag and namespacing through the filenames that you keep
# them in.
#
# It is essentially a wrapper that drives using fzf (so you'll need to have that
# installed) which is a little more user friendly than simply grepping text files
# on your filesystem. I'm a programmer so for me I want to be able to look at
# text, code and links: so that's what k supports. I may end up adding some more
# complicated actions around this in the future but for now, it simply hunts down
# that bit of information that you are sure you knew at some point but can't quite
# put your finger on.
#
# k is named for my amazing wife (Katie) who has spent far too long doing this
# sort of thing for me herself (remembering things I mean: not coding up utility
# scripts). It's also kind of an in joke we have around me answering 'Que?' when
# am asked a simple question that I'm drawing a blank on.

# == vars ==
typeset -AH C
C=(
  red   '\033[31m' green   '\033[32m' yellow  '\033[33m'
  blue  '\033[34m' purple  '\033[35m' cyan    '\033[36m'
  white '\033[37m' bold    '\033[1m'  nc      '\033[0m'
)
missing_snippet="$C[red]$C[bold]-- ERROR --\n\nunable to find that snippet$C[nc]"
default_helpfile_dir="$HOME/.helpfiles"
hpath=$HELPFILE_PATH
if [ -d $default_helpfile_dir ]; then
  hdirs+=":$default_helpfile_dir"
fi
hdirs=(${(@s/:/)hpath})


# == helpers ==
function error_and_exit {
  cecho red "$C[red]$C[bold]error:$C[nc] $*" 2>&1
  exit 1
}

function require_external {
  for prog in $*; do
    if ! [ -x "$(command -v $prog)" ]; then
      error_and_exit "'$prog' is required for k to run"
    fi
  done
}

# strip surrounding whitespace without messing up quotes like xargs does
alias trim="sed 's/^[[:space:]]*//; s/[[:space:]]*$//'"


# == functions ==
function generate_stubs {
  local d abs_path snippets

  for d in $hdirs; do
    for abs_path in $d/*(.); do
      snippets+="$(
        awk -v f=${$(basename $abs_path)%.*} '
          /^#/  { printf f "ζ|  "; for (i=2; i<=NF; i++) printf $i " "; printf "ζ|  " };
          /^\?/ { for (i=2; i<NF; i++) printf $i " "; print $NF };
        ' $abs_path
      )\n"
    done
  done

  [ -z "$snippets" ] && error_and_exit "no snippets found"
  # NOTE: we _should_ be able to just specify an output delimiter, but the BSD/OSX
  #       version of column doesn't have that flag...(using ζ to hopefully avoid clashes!)
  echo $snippets | column -t -s 'ζ' | sed 's/ | /|/g' | sort
}


function ansi_snippet {
  local d p title abs_path fname_target raw stub=$1

  fname_target="$(echo $stub | cut -d'|' -f1 | trim)"
  title="$(echo $stub | cut -d'|' -f2 | trim | sed 's./.\\/.')"

  for d in $hdirs; do
    for p in $d/*(.); do
      if [[ "${$(basename $p)%.*}" == $fname_target ]]; then
        raw=$(sed -n "/# *$title/,/^--/p" $p)
        [ -n "$raw" ] && echo $raw | awk \
          -v heading="$C[yellow]$C[bold]" \
          -v url="$C[blue]$C[bold]" \
          -v desc="$C[bold]" \
          -v nc="$C[nc]" '
          /^#/  { printf heading; for (i=2; i<=NF; i++) printf $i " "; print nc }
          /^%/  { printf url; for (i=2; i<=NF; i++) printf $i " "; print nc }
          /^>/  { printf desc; for (i=2; i<=NF; i++) printf $i " "; print nc }
          /^\$/ { print }
          /^$/  { print "" }
        '
        return
      fi
    done
  done
  echo $missing_snippet
}

# == main ==
require_external fzf column

if [[ "$1" == "___preview" ]]; then
  ansi_snippet $2
else
  selection=$(
    generate_stubs |
    fzf \
      --preview-window=up:70% \
      --preview "$0 ___preview {}"
  )

  [ -n "$selection" ] && ansi_snippet $selection
fi
