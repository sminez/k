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
heading="$C[yellow]$C[bold]"
url="$C[blue]$C[bold]"
desc="$C[bold]"
nc="$C[nc]"

case $(uname -s) in
  Darwin) 
    alias clipboard='pbcopy';
  ;;
  Linux)
    alias clipboard='xclip -selection clipboard';
    alias open='xdg-open';
  ;;
  *)
    alias clipboard='cat && echo -e "\n\nunknown OS: no clipboard support"';
    alias open='cat && echo -e "\n\nunknown OS: no open in browser support"';
  ;;
esac

hpath=$HELPFILE_PATH
if [ -d $default_helpfile_dir ]; then
  hdirs+=":$default_helpfile_dir"
fi
hdirs=(${(@s/:/)hpath})

# == helpers ==
function usage {
  cat <<EOF
usage: k <flags>

flags:
  --no-color|-n     don't use ANSI color codes when printing the selected entry
  --open|-o         open the first url listed for the selected snippet in your
                    default web browser
  --clip|-c         copy the selected entry to the system clipboard instead of
                    printing it to the terminal (implies --no-color)
  --help|-h         display this help message
EOF
}

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
  local stub=$1 _heading=$2 _url=$3 _desc=$4 _nc=$4
  local d p title abs_path fname_target raw

  fname_target="$(echo $stub | cut -d'|' -f1 | trim)"
  title="$(echo $stub | cut -d'|' -f2 | trim | sed 's./.\\/.')"

  for d in $hdirs; do
    for p in $d/*(.); do
      if [[ "${$(basename $p)%.*}" == $fname_target ]]; then
        raw=$(sed -n "/#[[:space:]]*$title/,/^--/p" $p)
        if [ -n "$raw" ]; then
          echo $raw |
          awk -v heading="$_heading" -v url="$_url" -v desc="$_desc" -v nc="$_nc" '
            /^#/  { printf heading; for (i=2; i<=NF; i++) printf $i " "; print nc }
            /^%/  { printf url; for (i=2; i<=NF; i++) printf $i " "; print nc }
            /^>/  { printf desc; for (i=2; i<=NF; i++) printf $i " "; print nc }
            /^\$/ { print }
            /^$/  { print "" }
          '
          return
        fi
      fi
    done
  done
  echo $missing_snippet
}

# == main ==
require_external fzf column


if [[ "$1" == "___preview" ]]; then
  ansi_snippet $2 $heading $url $desc
  exit 0
fi

selection=$(
  generate_stubs |
  fzf \
    --preview-window=up:70% \
    --preview "$0 ___preview {}"
)

[ -z "$selection" ] && exit 0

case $1 in
  --no-color|-n) ansi_snippet $selection;;
      --open|-o) open "$(ansi_snippet $selection | grep 'http' | head -n1)";;
      --clip|-c) ansi_snippet $selection | clipboard;;
              *) ansi_snippet $selection $heading $url $desc $nc;;
esac
