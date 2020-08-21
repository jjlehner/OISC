#include <iostream>
#include <fstream>
#include <sstream>
#include "simulator.hpp"
#include <cassert>
int main(){
	OISC_SIM sim;
	std::fstream file("simulator/test_programs/basics/initial.subleq");
	if( file ){
		std::stringstream buffer;
		buffer << file.rdbuf();
		sim.load_program(buffer);
		sim.dump_ram(std::cerr);

		sim.run();
		sim.load_program(buffer);
		sim.dump_ram(std::cerr);

	}
	else{
		assert(0 && "File not found");
	}
	return 0;
}