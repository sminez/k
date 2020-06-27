" Vim syntax file
" Language: k helpfiles
" Maintainer: Innes Anderson-Morrison
" Latest Revision: 27 June 2020

if exists("b:current_syntax")
  finish
endif

let b:current_syntax = "helpfile"

syn match hfPrefix contained '^[#?%>$]'
syn match hfTitle '^#.*$'  contains=hfPrefix
syn match hfTags  '^?.*$'  contains=hfPrefix
syn match hfLink  '^%.*$'  contains=hfPrefix
syn match hfText  '^>.*$'  contains=hfPrefix
syn match hfCode  '^\$.*$' contains=hfPrefix
syn match hfSplit '^--$'

hi def link hfPrefix Comment
hi def link hfTitle  Todo
hi def link hfTags   Type
hi def link hfLink   Identifier
hi def link hfText   String
hi def link hfCode   None
hi def link hfSplit  Keyword
" hi def link hfSplit  Error
