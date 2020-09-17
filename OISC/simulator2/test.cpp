#include <iostream>
struct test{
	static int i;
	test(){
		i++;
		std::cout<<"constructed"<<std::endl;
	}
	~test(){
		i--;
		std::cout<<"destroyed"<<std::endl;
	}
};
int test::i;
int main(){
	std::cout<<test::i<<std::endl;
}