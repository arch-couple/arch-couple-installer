from pathlib import Path


class User:
    def __init__(self, username: str, password: str, homepath: Path, sudoer: bool):
        self.username = username
        self.password = password
        self.homepath = homepath
        self.sudoer = sudoer

    def to_string(self) -> str:
        return f"{self.username}:\nPassword: {self.password}\nHomepath: {self.homepath}\nIs Sudoer?: {self.sudoer}"
