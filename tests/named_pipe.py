"""
execute-assembly data transfer over named pipe test
"""

import win32pipe, win32file, pywintypes

def pipe_server():
    pipe_name = r'\\.\pipe\temp'

    while True:
        print("Waiting for client...")
        pipe = win32pipe.CreateNamedPipe(
            pipe_name,
            win32pipe.PIPE_ACCESS_INBOUND,
            win32pipe.PIPE_TYPE_MESSAGE | win32pipe.PIPE_READMODE_MESSAGE | win32pipe.PIPE_WAIT,
            1, 65536, 65536,
            0,
            None
        )

        try:
            win32pipe.ConnectNamedPipe(pipe, None)
            print("Client connected")

            while True:
                resp = win32file.ReadFile(pipe, 64*1024)
                print(f"Received: {resp[1].decode()}")
        except pywintypes.error as e:
            if e.args[0] == 109: 
                print("Client disconnected")
        finally:
            win32file.CloseHandle(pipe)

if __name__ == '__main__':
    pipe_server()
