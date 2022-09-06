#pragma once
#include <iostream>
#include <vector>
using namespace std;

class Memory {
public:
	class MemoryBlock {
	public:
		int startIndex;
		int length;
		bool isOccupied;
		string occupiedBy;

		MemoryBlock();
		MemoryBlock(int startIndex, int length, bool isOccupied, string occupiedBy);

		int getEndIndex();
	};

	struct Statistics {
	public:
		int freeBlocksCount;
		int occupiedBlocksCount;
		int freeMemory;
		int occupiedMemory;
		vector<MemoryBlock> blocks;
	};
private:
	class Node {
	public:
		Node* prev;
		Node* next;
		MemoryBlock info;

		Node();
		Node(MemoryBlock info);
	};

	Node* head;
	Node* tail;
	Statistics statistics;

	Node* getByBestFit(int size);
	vector<Node*> getAllForProcess(string process);
	pair<Node*, Node*> splitNode(Node* node, int leftSize, string process);
	Node* joinWithNearestFreeNodes(Node* node);
public:
	Memory(int maxMemorySize);
	MemoryBlock addByBestFit(int size, string process);
	void free(int startIndex);
	void freeAllForProcess(string process);
	Statistics getStatistics();
};

std::ostream& operator<< (std::ostream&, Memory::Statistics);
std::ostream& operator<< (std::ostream&, Memory::MemoryBlock);