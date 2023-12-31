package commands

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

func GetFolderSize(path string) (int64, error) {
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

func Pwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return err.Error()
	}

	return dir
}

func Ls(path string) ([]File, error) {
    dir, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer dir.Close()

    fileInfos, err := dir.Readdir(-1)
    if err != nil {
        return nil, err
    }

    files := make([]File, len(fileInfos))
    for i, fileInfo := range fileInfos {
        size := fileInfo.Size()
        if fileInfo.IsDir() {
            size, err = GetFolderSize(filepath.Join(path, fileInfo.Name()))
            if err != nil {
                return nil, err
            }
        }
        files[i] = File{
            Filename: fileInfo.Name(),
            Size:     size,
            IsDir:    fileInfo.IsDir(),
            ModTime:  fileInfo.ModTime(),
        }
    }
    return files, nil
}

func Whoami() string {
	user, err := user.Current()
	if err != nil {
		return err.Error()
	}

	return user.Username
}

func Shell(command string) (string, error) {
    cmd := exec.Command("cmd.exe", "/c", command)
    output, err := cmd.CombinedOutput()
    return string(output), err
}

func Cat(filename string) (string, error) {
	// read file without using cmd

	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	defer file.Close()

	buf := make([]byte, 1024)
	var output string

	for {
		n, err := file.Read(buf)
		if err != nil {
			break
		}

		output += string(buf[:n])
	}

	return output, nil
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

func Ps() (string, error) {
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

func Get_pid() string {
	return fmt.Sprintf("%d", os.Getpid())
}