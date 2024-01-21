class SocketConnection:
    def __init__(self, socket, address):
        self.socket = socket
        self.address = address
        self.received = b""
        self.to_send = b""

    def receive_data_from_peer(self, data: bytes) -> None:
        self.received += data

    def read_received_data(self) -> bytes:
        data = str(self.received)
        self.received = b""
        return data

    def send_data_to_peer(self, data: bytes) -> None:
        self.to_send += data

    def get_data_to_send(self) -> bytes:
        data = self.to_send
        self.to_send = b""
        return data
