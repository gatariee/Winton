BINARY_NAME = teamserver
BIN_FOLDER = ./bin

.PHONY: all windows linux clean

all: teamserver-win64
	@echo [*] Default build is for Windows x64, use 'make windows' or 'make linux' to build for other platforms
	@echo [*] Done compiling for Windows x64, binary is located at: $(BIN_FOLDER)/$(BINARY_NAME)-win64.exe


windows: teamserver-win64 teamserver-win32
	@echo [*] Done compiling for Windows x64 and x86, binary is located at: $(BIN_FOLDER)


linux: teamserver-x64 teamserver-x86
	@echo [*] Done compiling for Linux x64 and x86, binary is located at: $(BIN_FOLDER)

teamserver-x64:
	@echo [*] Building $(BINARY_NAME)-linux-x64...
	@set GOOS=linux&& set GOARCH=amd64&& go build -o $(BIN_FOLDER)/$(BINARY_NAME)-linux-x64

teamserver-x86:
	@echo [*] Building $(BINARY_NAME)-linux-x86...
	@set GOOS=linux&& set GOARCH=386&& go build -o $(BIN_FOLDER)/$(BINARY_NAME)-linux-x86

teamserver-win64:
	@echo [*] Building $(BINARY_NAME)-win64.exe...
	@set GOOS=windows&& set GOARCH=amd64&& go build -o $(BIN_FOLDER)\$(BINARY_NAME)-win64.exe

teamserver-win32:
	@echo [*] Building $(BINARY_NAME)-win32.exe...
	@set GOOS=windows&& set GOARCH=386&& go build -o $(BIN_FOLDER)\$(BINARY_NAME)-win32.exe

clean:
	@if exist "$(BIN_FOLDER)\$(BINARY_NAME)-linux-x64" del /f "$(BIN_FOLDER)\$(BINARY_NAME)-linux-x64"
	@if exist "$(BIN_FOLDER)\$(BINARY_NAME)-linux-x86" del /f "$(BIN_FOLDER)\$(BINARY_NAME)-linux-x86"
	@if exist "$(BIN_FOLDER)\$(BINARY_NAME)-win64.exe" del /f "$(BIN_FOLDER)\$(BINARY_NAME)-win64.exe"
	@if exist "$(BIN_FOLDER)\$(BINARY_NAME)-win32.exe" del /f "$(BIN_FOLDER)\$(BINARY_NAME)-win32.exe"
	@echo [*] Done.
