#include <array>
#include <cassert>
#include <cinttypes>
#include <cmath>
#include <filesystem>
#include <fstream>
#include <iostream>
#include <limits>
#include <sstream>

constexpr unsigned int word_size = 48;
constexpr unsigned int bit_system = word_size / 3;

struct word
{
	__uint128_t internal;
	static constexpr unsigned int R_SHIFT_DISTANCE = 128 - bit_system;
	uint32_t a() const
	{
		static constexpr unsigned int L_SHIFT_DISTANCE = 128 - 3 * bit_system;
		return ( internal << L_SHIFT_DISTANCE ) >> ( R_SHIFT_DISTANCE );
	}
	uint32_t b() const
	{
		static constexpr unsigned int L_SHIFT_DISTANCE = 128 - 2 * bit_system;
		return ( internal << L_SHIFT_DISTANCE ) >> ( R_SHIFT_DISTANCE );
	}
	uint32_t c() const
	{
		static constexpr unsigned int L_SHIFT_DISTANCE = 128 - bit_system;
		return ( internal << L_SHIFT_DISTANCE ) >> ( R_SHIFT_DISTANCE );
	}
	bool leq() const
	{
		return internal >> ( word_size - 1 ) || internal == 0;
	}
	word( uint64_t internal = 0 )
	{
		this->internal = internal;
	}
	word( __uint128_t internal )
	{
		this->internal = internal;
	}
};

uint32_t pc = 1;
constexpr unsigned long long RAM_SIZE = 1ll << ( bit_system );
std::array<word, RAM_SIZE> memory{};

void operator-=( word &first, const word &second )
{
	first.internal -= second.internal;
}

std::ostream &operator<<( std::ostream &stream, const word &num )
{
	uint32_t a = num.a();
	uint32_t b = num.b();
	uint32_t c = num.c();
	char buf[100];
	static constexpr unsigned int lead = bit_system / 4;
	sprintf( buf, "%0*" PRIX32 "%0*" PRIX32 "%0*" PRIX32, lead, a, lead, b, lead, c );
	stream << buf;
	return stream;
}
word &operator<<( word &w, unsigned int num )
{
	w.internal = w.internal << num;
	w.internal %= (__uint128_t)1 << word_size;
	return w;
}
void dumpRam( std::ostream &stream, unsigned long int end, unsigned long int start = 0 )
{

	for ( size_t a = start; a < RAM_SIZE && a <= end; a++ )
	{
		word temp = memory[a];
		stream << std::hex << a << "\t" << temp << std::endl;
	}
}

void loadProgram( std::istream &file )
{
	memory.fill( (uint64_t)0 );
	std::string line;
	while ( std::getline( file, line ) )
	{
		std::stringstream ss( line );
		std::string addr;
		std::string data;
		ss >> addr >> data;

		int addr_pow = 10;
		int data_pow = 10;

		if ( addr.size() > 2 )
		{
			addr_pow = ( addr.substr( 0, 2 ) ) == "0x" ? 16 : 10;
			addr = ( addr_pow == 16 ) ? addr.substr( 2 ) : addr;
		}
		if ( data.size() > 2 )
		{
			data_pow = ( data.substr( 0, 2 ) ) == "0x" ? 16 : 10;
			data = ( data_pow == 16 ) ? data.substr( 2 ) : data;
		}
		uint64_t addr_val{ 0 };
		__uint128_t data_val{ 0 };

		for ( char c : addr )
		{
			int num = 0;
			num = ( isdigit( c ) ) ? ( c - '0' ) : num;
			num = ( isalpha( c ) && islower( c ) ) ? ( c - 'a' + 10 ) : num;
			num = ( isalpha( c ) && isupper( c ) ) ? ( c - 'A' + 10 ) : num;
			addr_val *= addr_pow;
			addr_val += num;
		}
		for ( char c : data )
		{
			int num = 0;
			num = ( isdigit( c ) ) ? ( c - '0' ) : num;
			num = ( isalpha( c ) && islower( c ) ) ? ( c - 'a' + 10 ) : num;
			num = ( isalpha( c ) && isupper( c ) ) ? ( c - 'A' + 10 ) : num;
			data_val *= addr_pow;
			data_val += num;
		}
		memory[addr_val] = word( data_val );
	}
}
int main()
{

	static_assert( ( word_size % 3 ) == 0, "Word size must be a multiple of 3 & 4" );
	static_assert( ( word_size % 4 ) == 0, "Word size must be a multiple of 3 & 4" );
	static_assert( word_size <= 128, "Word size must be <= 128" );

	std::ifstream file;
	file.open( "simulator2/test_programs/basics/initial.subleq" );
	if ( file.is_open() )
	{
		loadProgram( file );
	}
	else
	{
		std::cout << "could not find file" << std::endl;
	}
	dumpRam( std::cout, 0xc );
	while ( pc != 0 )
	{
		word &pc_word = memory[pc];
		memory[pc_word.b()] -= memory[pc_word.a()];
		pc = ( pc_word.leq() ) ? pc_word.c() : pc + 1;
	}
	std::cout<<std::endl<<std::endl;
	dumpRam( std::cout, 0x100 );
}