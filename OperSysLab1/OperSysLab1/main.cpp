#include <iostream>
#include "memory.h"

using namespace std;


int main() {
	Memory mem = Memory(100);
	cout << mem.getStatistics() << endl;
	
	mem.addByBestFit(20, "process_a");
	mem.addByBestFit(30, "process_b");
	mem.addByBestFit(10, "process_b");
	mem.addByBestFit(25, "process_c");
	
	cout << mem.getStatistics() << endl;

	mem.free(0);

	cout << mem.getStatistics() << endl;

	mem.freeAllForProcess("process_b");

	cout << mem.getStatistics() << endl;

	return 0;
}