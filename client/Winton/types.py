from dataclasses import dataclass

@dataclass
class Agent:
    IP: str
    Hostname: str
    Sleep: str
    UID: str

    def winton(self) -> dict:
        return self.__dict__

@dataclass
class File:
    Filename: str
    Size: int
    IsDir: bool
    ModTime: str

    def winton(self) -> dict:
        return self.__dict__


@dataclass
class CommandData:
    CommandID: str
    Command: str

    def winton(self) -> dict:
        return self.__dict__

@dataclass
class Command:
    name: str
    description: str
    usage: str

    def __str__(self):
        return f"{self.name}\t\t{self.description}\nUsage: {self.usage}"

@dataclass
class Result:
    CommandID: str
    Result: str

    def winton(self) -> dict:
        return self.__dict__

@dataclass
class ResultList:
    results: list[Result]

    def winton(self) -> dict:
        return self.__dict__