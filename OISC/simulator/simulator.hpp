#ifndef GUARD_SIMULATOR_HPP
#define GUARD_SIMULATOR_HPP

#include "simulator.hpp"
#include <bitset>
#include <cassert>
#include <iomanip>
#include <algorithm>
#include <vector>

constexpr unsigned int ADDRESS_WIDTH = 16;
constexpr unsigned int RAM_WORD_WIDTH = ADDRESS_WIDTH * 3;
constexpr unsigned int RAM_SIZE = 1 << ADDRESS_WIDTH;

class OISC_SIM
{
	using word = std::bitset<RAM_WORD_WIDTH>;
	using address = std::bitset<ADDRESS_WIDTH>;
	word ram[RAM_SIZE];

	word &pc = ram[RAM_SIZE - 5];
	word &carry_value = ram[RAM_SIZE - 4];
	word &carry_enable = ram[RAM_SIZE - 3];
	word &link_register = ram[RAM_SIZE - 2];
	word &stack_pointer = ram[RAM_SIZE - 1];

	bool subleq()
	{
		word a_ram_word = ram[pc.to_ulong()];
		word b_ram_word = a_ram_word;
		word c_ram_word = a_ram_word;

		if ( pc.to_ulong() == 0 )
		{
			return false;
		}

		a_ram_word >>= ( ADDRESS_WIDTH * 2 );
		( b_ram_word <<= ADDRESS_WIDTH ) >>= ( ADDRESS_WIDTH * 2 );
		( c_ram_word <<= ( ADDRESS_WIDTH * 2 ) ) >>= ( ADDRESS_WIDTH * 2 );
		address a( a_ram_word.to_ulong() );
		address b( b_ram_word.to_ulong() );
		address c( c_ram_word.to_ulong() );
		if ( !minus_equals( b, a ) )
		{
			pc = c.to_ulong();
		}
		else
		{
			pc = pc.to_ulong() + 1ul;
		}
		dump_ram(std::cout);
		return true;
	}
	// returns false if negative
	bool minus_equals( const address &x, const address &y )
	{
		word &xVal = ram[x.to_ulong()];
		word yVal = ram[y.to_ulong()];
		yVal.flip();

		bool carry = false;
		word one = {1};

		for ( int i = 0; i < RAM_WORD_WIDTH; i++ )
		{
			bool b = one[i];
			bool a = yVal[i];
			bool bit = a ^ ( b ^ carry  );
			carry = a && b || a && carry || b && carry;
			yVal[i] = bit;
		}
		carry = carry_value[0] && carry_enable[0];
		for ( int i = 0; i < RAM_WORD_WIDTH; i++ )
		{
			bool b = xVal[i];
			bool a = yVal[i];
			bool bit = a ^ ( b ^ carry  );
			carry = a && b || a && carry || b && carry;
			xVal[i] = bit;
		}
		carry_value[0] = carry;
		
		if ( xVal[RAM_WORD_WIDTH - 1] || xVal.to_ulong() == 0 )
		{
			return false;
		}
		return true;
	}

  public:
	OISC_SIM(){
		pc = 1;

	}

	void run()
	{
		// clang-format off
		while ( subleq() );
		// clang-format on
	}

	void load_program( std::stringstream &program )
	{
		std::string line;
		while ( std::getline( program, line ) )
		{
			std::stringstream stream( line );
			long address_l;
			long op_l;
			stream >> std::hex >> address_l >> op_l;

			assert( op_l >= 0 && "Program Error -  Op can't be negative " );
			assert( op_l < ( 2ul << RAM_WORD_WIDTH ) && "Program Error = Op can't be greater than RAM_WORD_WIDTH" );
			assert( address_l <= RAM_SIZE && "Program Error - Address of Op too large" );

			ram[address_l] = word( op_l );
		}
	}
	void dump_ram( std::ostream &stream, const unsigned long &limit = 20 )
	{
		for ( unsigned long i = 0; i < limit && i < RAM_SIZE; i++ )
		{
			stream << std::hex << "0x" << i << " " << ram[i].to_ulong() << std::endl;
		}

		stream << "-----------------" << std::endl;
		stream << "program counter - " << pc.to_ulong() << std::endl;
	}
};

#endif
