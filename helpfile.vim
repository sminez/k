" Vim syntax file
" Language: k helpfiles
" Maintainer: Innes Anderson-Morrison
" Latest Revision: 27 June 2020

if exists("b:current_syntax")
  finish
endif

syn keyword hfPrefix contained # ? % > $
syn match hfTitle '^#.*$'  contains=hfPrefix
syn match hfTags  '^?.*$'  contains=hfPrefix
syn match hfLink  '^%.*$'  contains=hfPrefix
syn match hfText  '^>.*$'  contains=hfPrefix
syn match hfCode  '^\$.*$' contains=hfPrefix
syn match hfSplit '^--$'

let b:current_syntax = "helpfile"

hi def link hfPrefix Type
hi def link hfTitle  Todo
hi def link hfTags   Constant
hi def link hfLink   Comment
hi def link hfText   Statement
hi def link hfCode   PreProc
hi def link hfSplit  Conditional
