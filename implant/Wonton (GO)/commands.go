package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type File struct {
	Filename string
	Size     int64
	IsDir    bool
	ModTime  time.Time
}

func get_folder_size(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}

		return nil
	})

	return size, err
}

func pwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return err.Error()
	}

	return dir
}

func ls(path string) ([]File, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var files []File

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			size, err := get_folder_size(path + "/" + fileInfo.Name())
			if err != nil {
				return nil, err
			}

			file := File{
				Filename: fileInfo.Name(),
				Size:     size,
				IsDir:    fileInfo.IsDir(),
				ModTime:  fileInfo.ModTime(),
			}

			files = append(files, file)

		} else {
			file := File{
				Filename: fileInfo.Name(),
				Size:     fileInfo.Size(),
				IsDir:    fileInfo.IsDir(),
				ModTime:  fileInfo.ModTime(),
			}

			files = append(files, file)
		}
	}

	return files, nil
}

func whoami() string {
	user, err := user.Current()
	if err != nil {
		return err.Error()
	}

	return user.Username
}

func shell(command string) (string, error) {
	cmd := exec.Command("cmd.exe", "/c", command)
	stdout, err := cmd.Output()
	if err != nil {
		// return err.Error()
		return string(stdout), err
	}

	return string(stdout), nil
}

func getProcessArch(processHandle windows.Handle) (string, error) {
	var isWow64 bool
	err := windows.IsWow64Process(processHandle, &isWow64)
	if err != nil {
		return "", err
	}

	if isWow64 {
		return "x86", nil
	}

	return "x64", nil
}

func getProcessSession(processID uint32) (string, error) {
	var sessionID uint32

	res := windows.ProcessIdToSessionId(processID, &sessionID)
	if res != nil {
		return "", res
	}

	return fmt.Sprintf("%d", sessionID), nil
}

func getProcessOwner(processHandle windows.Handle) (string, error) {
	var token windows.Token
	err := windows.OpenProcessToken(processHandle, windows.TOKEN_QUERY, &token)
	if err != nil {
		if errno, ok := err.(syscall.Errno); ok && errno == 5 {
			return "NA", nil
		} else {
			return "", err
		}
	}

	defer token.Close()

	user, _ := token.GetTokenUser()

	userName, domainName, _, _ := user.User.Sid.LookupAccount("")

	return fmt.Sprintf("%s\\%s", domainName, userName), nil
}	

func ps() (string, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return "", err
	}
	defer windows.CloseHandle(snapshot)

	var pe32 windows.ProcessEntry32
	pe32.Size = uint32(unsafe.Sizeof(pe32))
	if err := windows.Process32First(snapshot, &pe32); err != nil {
		return "", err
	}

	var output string
	output += " PID   PPID  Name                                   Arch  Session     User\n"
	output += " ---   ----  ----                                   ----  -------     ----\n"

	for {
		pid := pe32.ProcessID
		ppid := pe32.ParentProcessID
		exeFile := syscall.UTF16ToString(pe32.ExeFile[:])

		if pid < 100 {
			arch := "NA"
			session := "NA"
			user := "NA"
			output += fmt.Sprintf("%5d %5d %-40s %-5s %-10s %-20s\n", pid, ppid, exeFile, arch, session, user)
		} else {
			processHandle, _ := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, pid)

			defer windows.CloseHandle(processHandle)

			arch, _ := getProcessArch(processHandle)
			session, _ := getProcessSession(pid)


			user, err := getProcessOwner(processHandle)
			if err != nil {
				user = "NA"
			}
			output += fmt.Sprintf("%5d %5d %-40s %-5s %-10s %-20s\n", pid, ppid, exeFile, arch, session, user)

		}

		if err := windows.Process32Next(snapshot, &pe32); err != nil {
			break
		}
	}

	return output, nil
}

func get_pid() string {
	return fmt.Sprintf("%d", os.Getpid())
}