#include "menu.h"
#include<string>
using namespace std;

Menu::Menu() {
	this->commands = vector<string>(0);
	this->commands.push_back("occupy <start_index> <process_name>");
	this->commands.push_back("free <start_index>");
	this->commands.push_back("free_all <process_name>");
	this->commands.push_back("stats");

	this->memory = Memory(256);
}
Menu::~Menu() {
	this->memory.~Memory();
}

bool Menu::show() {
	cout << "Avaliable commands:" << endl;
	for (int i = 0; i < this->commands.size(); i++) {
		cout << i + 1 << ". " << commands[i] << endl;
	}

	string str;
	getline(cin, str);
	auto args = vector<string>(0);
	while (str.find(" ") != string::npos) {
		args.push_back(str.substr(0, str.find(" ")));
		str = str.substr(str.find(" ") + 1);
	}
	args.push_back(str);
	
	if (args.size() < 1) {
		cout << "Empty command!" << endl;
		return false;
	}
	
	if (args[0] == "stats") {
		cout << this->memory.getStatistics() << endl;
		return true;
	}
	if (args[0] == "occupy") {
		if (args.size() < 3) {
			cout << "Not enougth args!" << endl;
			return false;
		}
		int startIndex = atoi(args[1].c_str());
		string process = args[2];
		cout << "Occupy block:" << endl << this->memory.addByBestFit(startIndex, process) << endl;
		return true;
	}
	if (args[0] == "free") {
		if (args.size() < 2) {
			cout << "Not enougth args!" << endl;
			return false;
		}
		int startIndex = atoi(args[1].c_str());
		cout << "Free block:" << endl << this->memory.free(startIndex) << endl;
		return true;
	}
	if (args[0] == "free_all") {
		if (args.size() < 2) {
			cout << "Not enougth args!" << endl;
			return false;
		}
		string process = args[1];
		cout << "Free blocks:" << this->memory.freeAllForProcess(process);
		return true;
	}
	cout << "Unknown command!" << endl;
	return true;
}