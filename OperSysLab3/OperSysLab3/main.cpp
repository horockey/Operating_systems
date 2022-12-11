#include <fstream>
#include <string>
#include <windows.h>
#include <tchar.h>
#include <iostream>
#include <filesystem>
#include <fstream>

#define BUF_SIZE 4096


using namespace std;

void executeProcess(string processNameString, int timeoutSeconds) {
	STARTUPINFO si = { sizeof(si) };
	PROCESS_INFORMATION pi;
	ZeroMemory(&si, sizeof(si));
	ZeroMemory(&pi, sizeof(pi));

	wstring processName(processNameString.length(), L' ');
	copy(processNameString.begin(), processNameString.end(), processName.begin());

	if (!CreateProcess(
		processName.c_str(), // process name
		NULL,					// args
		NULL,					// process security attributes
		NULL,					// thread security attributes
		false,					// handlers inheritage
		NORMAL_PRIORITY_CLASS,	// creation parameters
		NULL,					// env
		NULL,					// cwd
		&si,					// return: StartupInfo
		&pi						// return: ProcessInfo
	)) {
		cout << "failed to create process" << endl;
		return;

	}
	DWORD result;
	if (timeoutSeconds > 0) {
		result = WaitForSingleObject(pi.hProcess, 1'000 * timeoutSeconds);
	}
	else {
		result = WaitForSingleObject(pi.hProcess, INFINITE);
	}
	if (result == WAIT_TIMEOUT) {
		cout << "Terminated by timeout" << endl;
		TerminateProcess(pi.hProcess, 0);
	}
	CloseHandle(pi.hThread);
	CloseHandle(pi.hProcess);
}

bool setStreams() {
	ifstream fin; fin.open("D:\\Repos\\Operating_systems\\OperSysLab3\\OperSysLab3\\problem_3_files\\input.txt");
	ofstream fout; fout.open("D:\\Repos\\Operating_systems\\OperSysLab3\\OperSysLab3\\problem_3_files\\output.txt");
	while (!fin.eof()) {
		string s;
		fin >> s;
		fout << s;
	}
	fin.close();
	fout.close();
	return false;
}

void executeFile(string fileNameString) {
	STARTUPINFO si = { sizeof(si) };
	PROCESS_INFORMATION pi;
	ZeroMemory(&si, sizeof(si));
	ZeroMemory(&pi, sizeof(pi));

	fileNameString = "/c " + fileNameString;
	string psProsessname = "C:\\Windows\\System32\\cmd.exe";

	wstring processName(psProsessname.length(), L' ');
	wstring fileName(fileNameString.length(), L' ');
	copy(psProsessname.begin(), psProsessname.end(), processName.begin());
	copy(fileNameString.begin(), fileNameString.end(), fileName.begin());

	if (!CreateProcess(
		processName.c_str(), // process name
		(LPWSTR)(fileName.c_str()),	// args
		NULL,					// process security attributes
		NULL,					// thread security attributes
		false,					// handlers inheritage
		NORMAL_PRIORITY_CLASS,	// creation parameters
		NULL,					// env
		NULL,					// cwd
		&si,					// return: StartupInfo
		&pi						// return: ProcessInfo
	)) {
		cout << "failed to create process" << endl;
		return;

	}
	DWORD result;
	result = WaitForSingleObject(pi.hProcess, INFINITE);
	if (result != WAIT_OBJECT_0) {
		cout << "Non-zero return code for " << pi.dwProcessId << endl;
	}
	CloseHandle(pi.hThread);
	CloseHandle(pi.hProcess);
}

void problem1() {
	ifstream fin; fin.open("processes.ini");
	while (!fin.eof()) {
		string str;
		getline(fin, str);
		string processName = str.substr(0, str.find(" "));
		int timeoutSeconds = atoi(str.substr(str.find(" ") + 1).c_str());
		executeProcess(processName, timeoutSeconds);
	}
	return;
}

void problem2() {
	string dirPath = "D:\\Repos\\Operating_systems\\OperSysLab3\\OperSysLab3\\problem_2_processes";
	for (const auto& entry : filesystem::directory_iterator(dirPath)) {
		if (entry.is_directory()) {
			continue;
		}
		string fileName = entry.path().string();
		if (
			fileName.substr(fileName.length() - 4, 4) != ".exe" &&
			fileName.substr(fileName.length() - 4, 4) != ".cmd" &&
			fileName.substr(fileName.length() - 4, 4) != ".bat"
			) {
			continue;
		}
		executeFile(entry.path().string());
		remove(entry.path().string().c_str());
	}
	return;
}

void problem3() {
	STARTUPINFO si = { sizeof(si) };
	PROCESS_INFORMATION pi;
	ZeroMemory(&si, sizeof(si));
	ZeroMemory(&pi, sizeof(pi));

	SECURITY_ATTRIBUTES saAttr;
	saAttr.nLength = sizeof(SECURITY_ATTRIBUTES);
	saAttr.bInheritHandle = TRUE;
	saAttr.lpSecurityDescriptor = NULL;

	HANDLE hFinR = NULL;
	HANDLE hFinW = NULL;
	HANDLE hFoutR = NULL;
	HANDLE hFoutW = NULL;
	
	HANDLE inFile = NULL;
	HANDLE outFile = NULL;

	BOOL isPipeProblem = false;

	isPipeProblem |= CreatePipe(&hFinR, &hFinW, &saAttr, 0);
	isPipeProblem |= SetHandleInformation(hFinR, HANDLE_FLAG_INHERIT, 0);
	isPipeProblem |= CreatePipe(&hFoutR, &hFoutW, &saAttr, 0);
	isPipeProblem |= SetHandleInformation(hFoutR, HANDLE_FLAG_INHERIT, 0);
	
	if (!setStreams()) {
		return;
	}
	if (isPipeProblem) {
		cout << "pipe problem" << endl;
		return;
	}

	si.cb = sizeof(STARTUPINFO);
	si.hStdInput = hFinR;
	si.hStdOutput = hFoutW;
	si.hStdError = hFoutW;
	si.dwFlags |= STARTF_USESTDHANDLES;

	inFile = CreateFile(
		(LPCWSTR)"D:\\Repos\\Operating_systems\\OperSysLab3\\OperSysLab3\\problem_3_files\\input.txt",
		GENERIC_READ,
		0,
		NULL,
		OPEN_EXISTING,
		FILE_ATTRIBUTE_READONLY,
		NULL
	);
	outFile = CreateFile(
		(LPCWSTR)"D:\\Repos\\Operating_systems\\OperSysLab3\\OperSysLab3\\problem_3_files\\output.txt",
		GENERIC_WRITE,
		0,
		NULL,
		OPEN_EXISTING,
		FILE_ATTRIBUTE_NORMAL,
		NULL
	);
	TCHAR szCmdline[] = TEXT("child");
	

	if (!CreateProcess(
		NULL, // process name
		szCmdline,					// args
		NULL,					// process security attributes
		NULL,					// thread security attributes
		true,					// handlers inheritage
		0,	// creation parameters
		NULL,					// env
		NULL,					// cwd
		&si,					// return: StartupInfo
		&pi						// return: ProcessInfo
	)) {
		cout << "failed to create process" << endl;
		return;

	}
	CloseHandle(pi.hProcess);
	CloseHandle(pi.hThread);

	
	DWORD dwRead, dwWritten;
	CHAR chBuf[BUF_SIZE];

	if (!ReadFile(inFile, chBuf, BUF_SIZE, &dwRead, NULL) || dwRead == 0) {
		cout << "failed to read file" << endl;
		return;
	}
	if (!WriteFile(outFile, chBuf, dwRead, &dwWritten, NULL)) {
		cout << "failed to write file" << endl;
		return;
	}

	CloseHandle(hFinR);
	CloseHandle(hFinW);
	CloseHandle(hFoutR);
	CloseHandle(hFoutW);
}

void problem5(string srcCatalogName, string dstCatalogName) {
	STARTUPINFO si = { sizeof(si) };
	PROCESS_INFORMATION pi;
	ZeroMemory(&si, sizeof(si));
	ZeroMemory(&pi, sizeof(pi));

	string psProcessName = "C:\\Windows\\System32\\cmd.exe";
	string commandString = "/c copy";

	for (const auto& entry : filesystem::directory_iterator(srcCatalogName)) {
		if (entry.is_directory()) {
			continue;
		}
		commandString = "/c copy";
		commandString += " " + entry.path().string() + " " + dstCatalogName;

		wstring processName(psProcessName.length(), L' ');
		wstring command(commandString.length(), L' ');
		copy(psProcessName.begin(), psProcessName.end(), processName.begin());
		copy(commandString.begin(), commandString.end(), command.begin());

		if (!CreateProcess(
			processName.c_str(), // process name
			(LPWSTR)(command.c_str()),	// args
			NULL,					// process security attributes
			NULL,					// thread security attributes
			false,					// handlers inheritage
			NORMAL_PRIORITY_CLASS,	// creation parameters
			NULL,					// env
			NULL,					// cwd
			&si,					// return: StartupInfo
			&pi						// return: ProcessInfo
		)) {
			cout << "failed to create process" << endl;
			return;

		}
		DWORD result;
		result = WaitForSingleObject(pi.hProcess, INFINITE);
		if (result != WAIT_OBJECT_0) {
			cout << "Non-zero return code for " << pi.dwProcessId << endl;
		}
		CloseHandle(pi.hThread);
		CloseHandle(pi.hProcess);
	}
}

int main() {
	//problem1();
	
	//problem2();
	
	problem3();
	
	/*problem5(
		"D:\\Repos\\Operating_systems\\OperSysLab3\\OperSysLab3\\problem_5\\src",
		"D:\\Repos\\Operating_systems\\OperSysLab3\\OperSysLab3\\problem_5\\dst"
	);*/
	return 0;
}