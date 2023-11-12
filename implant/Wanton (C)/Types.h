#pragma once

struct string {
	char* ptr;
	size_t len;
};

typedef struct {
	TCHAR** file_names;
	DWORD file_count;
} FileList;
