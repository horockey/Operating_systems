#pragma once
#include <map>
#include <iostream>
using namespace std;

#define Pap pair<MemoryManager::Address, MemoryManager::Memory::Page>

const int PAGE_SIZE = 32;

class MemoryManager {
private:
	class Address {
	public:
		int pageIndex;
		int delta;

		Address(int index);
		Address() {};
	};
	friend bool operator<(Address a, Address b);
	friend bool operator>(Address a, Address b);
	friend ostream& operator<<(ostream& out, MemoryManager::Address addr);
	class Memory {
	public:
		class Page {
		private:
			int size;
			char* storage;
			int ts;
		public:
			Page(int size);
			Page() {};

			char getByte(Address addr);
			void setByte(Address addr, char val);
			int getTs();
			void setTs(int val);
		};

		Memory(int size, int indexDelta);
		Memory() {};

		Pap popLRUPap();
		Pap getPap(Address addr);
		Pap popPap(Address addr);
		void setPap(Pap pairAdressPage);
		int size;
		map<Address, Page> pages;
		int currentTs;
	};

	Memory ram;
	Memory rom;
public:
	MemoryManager(int size);

	char getByte(int index);
	void setByte(int index, char val);
	void printStats();
};