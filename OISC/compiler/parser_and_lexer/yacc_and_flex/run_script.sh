#!/bin/bash
cd outputfiles
yacc -d ../grammar.y -o grammar.c
flex ../lexer.flex
cc -c lex.yy.c -o lex.yy.o
cd ..
g++ test.cpp outputfiles/lex.yy.o outputfiles/grammar.c -lm
