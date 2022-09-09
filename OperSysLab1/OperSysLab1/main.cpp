#include <iostream>
#include "memory.h"
#include "menu.h"

using namespace std;


int main() {
	/*Memory mem = Memory(256);
	cout << mem.getStatistics() << endl;
	
	cout << "Occupy block:" << endl << mem.addByBestFit(20, "process_a") << endl;
	cout << "Occupy block:" << endl << mem.addByBestFit(30, "process_b") << endl;
	cout << "Occupy block:" << endl << mem.addByBestFit(10, "process_b") << endl;
	cout << "Occupy block:" << endl << mem.addByBestFit(25, "process_c") << endl;
	
	cout << mem.getStatistics() << endl;

	cout << "Free block:" << endl << mem.free(0) << endl;

	cout << mem.getStatistics() << endl;

	cout << "Free blocks:" << mem.freeAllForProcess("process_b");

	cout << mem.getStatistics() << endl;*/

	auto menu = Menu();
	while (true) {
		if (!menu.show()) {
			break;
		}
	}
	menu.~Menu();
	return 0;
}