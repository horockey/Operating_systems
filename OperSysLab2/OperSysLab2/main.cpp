#include <iostream>
#include "memory_manager.h"
using namespace std;

int main() {
	/*auto mm = MemoryManager(1024);
	mm.printStats();
	mm.setByte(14, 120);
	mm.getByte(720);
	mm.printStats();
	cout << (int)(mm.getByte(14)) << endl;*/
	auto mm = MemoryManager(1024);
	int choise = 0;
	while (choise != 4) {
		cout << "1. Put value" << endl;
		cout << "2. Get value" << endl;
		cout << "3. Show RAM stats" << endl;
		cout << "4. Exit" << endl;

		cin >> choise;
		if (choise == 1) {
			int index, val;
			cout << "Index: ";
			cin >> index;
			cout << "Value: ";
			cin >> val;
			mm.setByte(index, val);
		}
		else if (choise == 2) {
			int index;
			cout << "Index: ";
			cin >> index;
			cout << (int)(mm.getByte(index)) << endl;
		}
		else if (choise == 3) {
			mm.printStats();
		}
		else if (choise == 4) {
			break;
		}
		else {
			cout << "Invalid command!";
		}
	}
	return 0;
}