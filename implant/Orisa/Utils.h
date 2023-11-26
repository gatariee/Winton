#pragma once

void init_string( struct string* s );

size_t writefunc( void* ptr, size_t size, size_t nmemb, struct string* s );

struct string http_get( char url[] );

struct string http_post_json( char url[], char json[] );