#define _CRT_SECURE_NO_WARNINGS
#include <stdio.h>
#include <windows.h>
#include <tchar.h>

#include "Types.h"


void Dir( const TCHAR* file_path, FileList* file_list ) {
	TCHAR Path[MAX_PATH];
	TCHAR SearchPattern[MAX_PATH];
	HANDLE hFind;
	WIN32_FIND_DATA FindFileData;

	_tcscpy( Path, file_path );
	_tcscpy( SearchPattern, _T( "*" ) );

	if (Path[_tcslen( Path ) - 1] != _T( '\\' ))
		_tcscat( Path, _T( "\\" ) );

	_tcscat( Path, SearchPattern );

	hFind = FindFirstFile( Path, &FindFileData );

	if (hFind == INVALID_HANDLE_VALUE)
	{
		_tprintf( _T( "Path not found: [%s]\n" ), file_path );
		return;
	}

	do
	{
		if (FindFileData.dwFileAttributes & FILE_ATTRIBUTE_DIRECTORY)
			continue;

		file_list->file_count++;
		file_list->file_names = realloc( file_list->file_names, file_list->file_count * sizeof( TCHAR* ) );
		file_list->file_names[file_list->file_count - 1] = malloc( (_tcslen( FindFileData.cFileName ) + 1) * sizeof( TCHAR ) );
		_tcscpy( file_list->file_names[file_list->file_count - 1], FindFileData.cFileName );


	} while (FindNextFile( hFind, &FindFileData ));

	FindClose( hFind );
}