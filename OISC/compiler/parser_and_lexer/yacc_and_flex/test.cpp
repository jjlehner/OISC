#include <iostream>
extern "C"{
int yyparse();
int yylex();
} 
int main()
{
    std::cout<<"out "<<yyparse()<<std::endl;
	std::cout<<"out "<<yyparse()<<std::endl;
}