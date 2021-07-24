import requests

payload = {
    "username": "john",
    "content": "This is a test right here right now",
}

r = requests.post("http://localhost:10000/send", json=payload)