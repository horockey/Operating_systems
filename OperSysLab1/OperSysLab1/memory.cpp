#include "memory.h"
#define nodesPair pair<Memory::Node*, Memory::Node*> 

Memory::Memory(int maxMemorySize) {
	this->head = new Memory::Node();
	this->head->info.startIndex = 0;
	this->head->info.length = maxMemorySize;
	
	this->tail = this->head;
}
Memory::~Memory() {
	Memory::Node* cur = nullptr;
	for (cur = this->head; cur != nullptr; cur = cur->next) {
		if (cur->prev != nullptr) {
			delete cur->prev;
		}
		if (cur == this->tail) {
			delete cur;
		}
	}
}

Memory::Node* Memory::getByBestFit(int size) {
	Memory::Node* cur;
	Memory::Node* bestFitNode = nullptr;
	for (cur = this->head; cur != nullptr; cur = cur->next) {
		if (cur->info.isOccupied || cur->info.length < size) {
			continue;
		}
		if (bestFitNode == nullptr || cur->info.length - size < bestFitNode->info.length - size) {
			bestFitNode = cur;
		}
	}
	return bestFitNode;
}
nodesPair Memory::splitNode(Node* node, int leftSize, string process) {
	auto splittedNodes = nodesPair();
	splittedNodes.first = new Memory::Node();
	splittedNodes.second = new Memory::Node();
	if (node->info.length < leftSize) {
		return nodesPair();
	}
	if (node->info.length == leftSize) {
		delete splittedNodes.second;
		splittedNodes.second = nullptr;
	}
	
	splittedNodes.first->info.startIndex = node->info.startIndex;
	splittedNodes.first->info.length = leftSize;
	splittedNodes.first->info.isOccupied = true;
	splittedNodes.first->info.occupiedBy = process;

	splittedNodes.first->prev = node->prev;

	if (splittedNodes.second != nullptr) {
		splittedNodes.second->info.startIndex = splittedNodes.first->info.getEndIndex() + 1;
		splittedNodes.second->info.length = node->info.length - leftSize;
		splittedNodes.second->prev = splittedNodes.first;
		splittedNodes.second->next = node->next;
		splittedNodes.first->next = splittedNodes.second;
	} else {
		splittedNodes.first->next = node->next;
	}
	
	if (node->prev != nullptr) {
		node->prev->next = splittedNodes.first;
	}
	if (node->next != nullptr) {
		if (splittedNodes.second != nullptr) {
			node->next->prev = splittedNodes.second;
		} else {
			node->next->prev = splittedNodes.first;
		}
	}
	return splittedNodes;
}
Memory::MemoryBlock Memory::addByBestFit(int size, string process) {
	
	auto node = Memory::getByBestFit(size);
	if (node == nullptr) {
		return Memory::MemoryBlock::MemoryBlock();
	}
	auto splittedNodes = Memory::splitNode(node, size, process);
	if (node == this->head) {
		this->head = splittedNodes.first;
	}
	if (node == this->tail) {
		if (splittedNodes.second != nullptr) {
			this->tail = splittedNodes.second;
		} else {
			this->tail = splittedNodes.first;
		}
	}
	return splittedNodes.first->info;
}

Memory::Node* Memory::joinWithNearestFreeNodes(Memory::Node* node) {
	if (node->info.isOccupied) {
		return nullptr;
	}
	Memory::Node* resNode = new Memory::Node();
	resNode->info.startIndex = node->info.startIndex;
	resNode->info.length = node->info.length;
	resNode->prev = node->prev;
	resNode->next = node->next;
	if (node->prev != nullptr) {
		node->prev->next = resNode;
		if(!node->prev->info.isOccupied) {
			resNode->info.startIndex = node->prev->info.startIndex;
			resNode->info.length += node->prev->info.length;
			resNode->prev = node->prev->prev;
			if (resNode->prev != nullptr) {
				resNode->prev->next = resNode;
			}
			if (node->prev == this->head) {
				this->head = resNode;
			}
			delete node->prev;
		}
	} else { // node == head
		this->head = resNode;
	}
	if (node->next != nullptr) {
		node->next->prev = resNode;
		if (!node->next->info.isOccupied) {
			resNode->info.length += node->next->info.length;
			resNode->next = node->next->next;
			if (resNode->next != nullptr) {
				resNode->next->prev = resNode;
			}
			if (node->next == this->tail) {
				this->tail = resNode;
			}
			delete node->next;
		}
	} else { // node== tail
		this->tail = resNode;
	}
	delete node;
	return resNode;
}
Memory::MemoryBlock Memory::free(int startIndex) {
	Node* cur = nullptr;
	for (cur = head; cur != nullptr; cur = cur->next) {
		if (cur->info.startIndex == startIndex) {
			break;
		}
	}
	if (cur == nullptr) {
		throw invalid_argument("no memory block with such startIndex");
		return Memory::MemoryBlock();
	}
	
	auto oldInfo = cur->info;

	cur->info.isOccupied = false;
	cur->info.occupiedBy = "";

	cur = Memory::joinWithNearestFreeNodes(cur);

	return oldInfo;
}

vector<Memory::Node*> Memory::getAllForProcess(string process) {
	vector<Memory::Node*> result;
	Memory::Node* cur = nullptr;
	for (cur = this->head; cur != nullptr; cur = cur->next) {
		if (cur->info.occupiedBy == process) {
			result.push_back(cur);
		}
	}
	return result;
}
vector<Memory::MemoryBlock> Memory::freeAllForProcess(string process) {
	auto nodes = Memory::getAllForProcess(process);
	auto oldInfo = vector<Memory::MemoryBlock>(0);
	for (auto node : nodes) {
		oldInfo.push_back(node->info);
		Memory::free(node->info.startIndex);
	}
	return oldInfo;
}

Memory::Statistics Memory::getStatistics() {
	Memory::Statistics stats;
	stats.freeMemory = 0;
	stats.occupiedMemory = 0;
	stats.freeBlocksCount = 0;
	stats.occupiedBlocksCount = 0;
	stats.blocks = vector<Memory::MemoryBlock>(0);

	Memory::Node* cur = nullptr;
	for (cur = head; cur != nullptr; cur = cur->next) {
		if (cur->info.isOccupied) {
			stats.occupiedMemory += cur->info.length;
			stats.occupiedBlocksCount++;
		} else {
			stats.freeMemory += cur->info.length;
			stats.freeBlocksCount++;
		}
		stats.blocks.push_back(cur->info);
	}
	return stats;
}
std::ostream& operator<< (std::ostream& out, Memory::Statistics stats) {
	out << "=============MEMORY STATS=============" << endl;
	
	out << "Memory (free/occupied): ";
	out << "(" << stats.freeMemory << " / " << stats.occupiedMemory << ")" << endl;

	out << "Blocks (free/occupied): ";
	out << "(" << stats.freeBlocksCount << " / " << stats.occupiedBlocksCount << ")" << endl;

	out << "Memory:" << endl;
	for (auto block : stats.blocks) {
		out << block << endl;
	}
	out << "=========END OF MEMORY STATS==========" << endl;
	return out;
}

Memory::Node::Node() {
	this->prev = nullptr;
	this->next = nullptr;
	this->info = MemoryBlock();
}
Memory::Node::Node(MemoryBlock info) {
	this->prev = nullptr;
	this->next = nullptr;
	this->info = info;
}
int Memory::MemoryBlock::getEndIndex() {
	return this->startIndex + this->length - 1;
}
std::ostream& operator<< (std::ostream& out, Memory::MemoryBlock block) {
	out << "Start index: " << block.startIndex << endl;
	out << "End index: " << block.getEndIndex() << endl;
	out << "Status: " << (block.isOccupied ? "Occupied" : "Free") << endl;
	out << "Occupier: " << block.occupiedBy << endl;
	return out;
}

Memory::MemoryBlock::MemoryBlock() {
	this->startIndex = -1;
	this->length = 0;
	this->isOccupied = false;
	this->occupiedBy = "";
}
Memory::MemoryBlock::MemoryBlock(int startIndex, int length, bool isOccupied, string occupiedBy) {
	this->startIndex = startIndex;
	this->length = length;
	this->isOccupied = isOccupied;
	this->occupiedBy = occupiedBy;
}

std::ostream& operator<< (std::ostream& out, vector<Memory::MemoryBlock> blocks) {
	for (auto block : blocks) {
		out << block << endl;
	}
	return out;
}