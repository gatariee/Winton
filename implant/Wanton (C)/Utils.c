#include <stdio.h>
#include <stdlib.h>
#include <curl/curl.h>
#include "types.h"

void init_string( struct string* s ) {
	s->len = 0;
	s->ptr = malloc( s->len + 1 );
	if (s->ptr == NULL) {
		fprintf( stderr, "malloc() failed\n" );
		exit( EXIT_FAILURE );
	}
	s->ptr[0] = '\0';
}

size_t writefunc( void* ptr, size_t size, size_t nmemb, struct string* s )
{
	size_t new_len = s->len + size * nmemb;
	s->ptr = realloc( s->ptr, new_len + 1 );
	if (s->ptr == NULL) {
		fprintf( stderr, "realloc() failed\n" );
		exit( EXIT_FAILURE );
	}
	memcpy( s->ptr + s->len, ptr, size * nmemb );
	s->ptr[new_len] = '\0';
	s->len = new_len;

	return size * nmemb;
}


struct string http_get( char url[] ) {
	CURL* curl;
	CURLcode res;

	curl = curl_easy_init();
	if (curl) {
		struct string s;
		init_string( &s );

		curl_easy_setopt( curl, CURLOPT_URL, url );
		curl_easy_setopt( curl, CURLOPT_WRITEFUNCTION, writefunc );
		curl_easy_setopt( curl, CURLOPT_WRITEDATA, &s );
		res = curl_easy_perform( curl );

		curl_easy_cleanup( curl );
		return s;
	}
}

struct string http_post_json( char url[], char json[] ) {
	CURL* curl;
	CURLcode res;

	curl = curl_easy_init();

	if (curl) {
		struct string s;
		init_string( &s );

		struct curl_slist* headers = NULL;
		headers = curl_slist_append( headers, "Content-Type: application/json" );

		curl_easy_setopt( curl, CURLOPT_URL, url );
		curl_easy_setopt( curl, CURLOPT_POSTFIELDS, json );
		curl_easy_setopt( curl, CURLOPT_HTTPHEADER, headers );
		curl_easy_setopt( curl, CURLOPT_WRITEFUNCTION, writefunc );
		curl_easy_setopt( curl, CURLOPT_WRITEDATA, &s );
		res = curl_easy_perform( curl );

		curl_easy_cleanup( curl );
		return s;
	}

}