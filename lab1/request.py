import ssl
import socket


def request(url):
    host, path = url.replace("https://", "").split("/", 1)
    port = 443

    context = ssl.create_default_context()

    sock = socket.create_connection((host, port))
    wrapped_socket = context.wrap_socket(sock, server_hostname=host)

    http_request = f"""GET /{path} HTTP/1.0\r
Host: {host}\r
Connection: close\r\n\r\n"""

    wrapped_socket.sendall(http_request.encode())

    response = b""
    while True:
        part = wrapped_socket.recv(4096)
        if not part:
            break
        response += part

    wrapped_socket.close()

    response_str = response.decode()
    headers, html_content = response_str.split("\r\n\r\n", 1)
    return html_content
