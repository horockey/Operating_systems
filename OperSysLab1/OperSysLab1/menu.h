#pragma once
#include <iostream>
#include<vector>
#include "memory.h"

using namespace std;

class Menu {
private:
	vector<string> commands;
	Memory memory;
public:
	Menu();
	~Menu();
	bool show();
};