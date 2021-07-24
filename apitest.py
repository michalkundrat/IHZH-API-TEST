import requests

payload = {
    "username": "PyTester0",
    "content": "Testing",
}

r = requests.post("http://localhost:10000/send", json=payload)