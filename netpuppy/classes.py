import threading


class SocketConnection:
    def __init__(self, socket, address):
        self.socket = socket
        self.address = address
        self.received = b""
        self.to_send = b""
        self.running = True

        # Start threads for reading and writing to socket:
        threading.Thread(target=self.read_stream).start()
        threading.Thread(target=self.write_stream).start()

    def read_stream(self) -> None:
        while self.running:
            data: bytes = self.socket.recv(1024)
            if not data:
                continue

            print(data.decode("utf-8"))
            data = b""

        return

    def write_stream(self) -> None:
        while self.running:
            try:
                data: str = input()
            except EOFError:
                self.running = False
            if not data:
                continue

            self.socket.sendall(data.encode("utf-8"))
            data = b""

        return
