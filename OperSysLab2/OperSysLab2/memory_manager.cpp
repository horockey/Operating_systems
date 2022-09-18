#include "memory_manager.h"

MemoryManager::Address::Address(int index) {
	this->pageIndex = (index & 0b11111111111111111111111111100000)>>5;
	this->delta = index &      0b00000000000000000000000000011111;
}
bool operator<(MemoryManager::Address a, MemoryManager::Address b) {
	return (a.pageIndex < b.pageIndex);
	
}
bool operator>(MemoryManager::Address a, MemoryManager::Address b) {
	return (a.pageIndex > b.pageIndex);
}
ostream& operator<<(ostream& out, MemoryManager::Address addr) {
	out << addr.pageIndex;
	return out;
}

MemoryManager::Memory::Page::Page(int size) {
	this->size = size;
	this->ts = 0;
	this->storage = new char[size];
	for (int i = 0; i < size; i++) {
		this->storage[i] = 0;
	}
}
char MemoryManager::Memory::Page::getByte(MemoryManager::Address addr) {
	if (addr.delta < 0) {
		throw invalid_argument("addr delta is negative");
	}
	if (addr.delta > this->size) {
		throw out_of_range("addr.delta too big");
	}
	return this->storage[addr.delta];
}
void MemoryManager::Memory::Page::setByte(MemoryManager::Address addr, char val) {
	if (addr.delta < 0) {
		throw invalid_argument("addr delta is negative");
	}
	if (addr.delta > this->size) {
		throw out_of_range("addr.delta too big");
	}
	this->storage[addr.delta] = val;
}
int MemoryManager::Memory::Page::getTs() {
	return this->ts;
}
void MemoryManager::Memory::Page::setTs(int val) {
	this->ts = val;
}

MemoryManager::Memory::Memory(int size, int indexDelta) {
	if (size < 0) {
		throw invalid_argument("memory size is negative");
	}
	if (size % PAGE_SIZE != 0) {
		throw invalid_argument("memory size must be miltiple to page size");
	}
	this->size = size;
	this->currentTs = 0;
	this->pages = map<MemoryManager::Address, MemoryManager::Memory::Page>();
	int numberOfPages = size / PAGE_SIZE;
	for (int i = 0; i < numberOfPages; i++) {
		int index = (i * PAGE_SIZE) + indexDelta;
		this->pages[MemoryManager::Address::Address(index)] = MemoryManager::Memory::Page::Page(PAGE_SIZE);
	}
}
Pap MemoryManager::Memory::popLRUPap() {
	Pap lruPap = *(this->pages.begin());
	int lruScore = 0;
	for (auto pap : this->pages) {
		if (this->currentTs - pap.second.getTs() > lruScore) {
			lruScore = this->currentTs - pap.second.getTs();
			lruPap = pap;
		}
	}
	this->pages.erase(lruPap.first);
	return lruPap;
}
Pap MemoryManager::Memory::getPap(MemoryManager::Address addr) {
	if (addr.pageIndex < 0) {
		throw invalid_argument("addr.pageIndex is negative");
	}
	auto it = this->pages.find(addr);
	if (it != this->pages.end()) {
		it->second.setTs(this->currentTs+1);
		this->currentTs++;

		return Pap(it->first, it->second);
	}
	throw out_of_range("no such page");
}
Pap MemoryManager::Memory::popPap(MemoryManager::Address addr) {
	if (addr.pageIndex < 0) {
		throw invalid_argument("addr.pageIndex is negative");
	}
	auto it = this->pages.find(addr);
	if (it != this->pages.end()) {
		auto pap = *it;
		this->pages.erase(it);
		return pap;
	}
	throw out_of_range("no such page");
}
void MemoryManager::Memory::setPap(Pap pap) {
	auto addr = pap.first;
	auto val = pap.second;
	if (addr.pageIndex < 0) {
		throw invalid_argument("addr.pageIndex is negative");
	}	
	this->pages[addr] = val;
		
	this->pages[addr].setTs(this->currentTs+1);
	this->currentTs++;
	
}

MemoryManager::MemoryManager(int size) {
	if (size < 512) {
		throw invalid_argument("too small MM size");
	}
	this->ram = Memory(256, 0);
	this->rom = Memory(size-256, 256);
}
char MemoryManager::getByte(int index) {
	auto addr = MemoryManager::Address::Address(index);
	auto pap = Pap(MemoryManager::Address::Address(), MemoryManager::Memory::Page::Page());
	try {
		pap = this->ram.getPap(addr);
	}
	catch (const out_of_range& e) {
		pap = this->rom.popPap(addr);
		auto movingPap = this->ram.popLRUPap();
		this->rom.setPap(movingPap);
		this->ram.setPap(pap);
	}
	return pap.second.getByte(addr);
}
void MemoryManager::setByte(int index, char val) {
	auto addr = MemoryManager::Address::Address(index);
	auto pap = Pap(MemoryManager::Address::Address(), MemoryManager::Memory::Page::Page());
	try {
		pap = this->ram.getPap(addr);
	}
	catch (const out_of_range& e) {
		pap = this->rom.popPap(addr);
		auto movingPap = this->ram.popLRUPap();
		this->rom.setPap(movingPap);
		this->ram.setPap(pap);
	}
	pap.second.setByte(addr, val);
}
void MemoryManager::printStats() {
	cout << "======= RAM STATS =======" << endl;
	cout << "RAM:" << endl;
	for (auto pap : this->ram.pages) {
		cout << pap.first << " (ts: " << pap.second.getTs() << ")" << endl;
	}
	cout << "=== END OF RAM STATS ====" << endl;
}