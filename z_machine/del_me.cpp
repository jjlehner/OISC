#include <iostream>

struct Pair{
    int first, second;
};

Pair fib(int a){
    if( a == 1 ){
        return Pair{0,1};
    }
    else if( a == 2 ){
        return Pair{1,1};
    }
    else{
        Pair p = fib( a - 1 );
        return Pair{p.second, p.first + p.second};
    }
}

int main(){
	std::cout<<fib(5).second;
}