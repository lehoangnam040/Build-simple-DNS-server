import socket
import sys


server_addr_port = ("0.0.0.0", 10053)
buffer_size = 1024

# Create a UDP socket at client side
udp_client_socket = socket.socket(family=socket.AF_INET, type=socket.SOCK_DGRAM)
# Send to server using created UDP socket
udp_client_socket.sendto(sys.argv[1].encode(), server_addr_port)
response = udp_client_socket.recvfrom(buffer_size)
print(f"Response from Server {response[0].decode()}")
