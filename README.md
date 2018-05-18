# Proxy Checker
It's just simple proxy checker which has created in few hours as a small home project.
So I didn't have time to implement command line interface or write tests.

## Configuration
You can override most of options in this program, just make `config.json` near the program:
```json
{
    "httpbin_host": "https://httpbin.example.org/ip",
    "concurency": 300,
    "cookies": {
        "private-access-key": "private-access-key",
        "private-access-key2": "private-access-key2"
    },
    "protos": ["http", "tcp"],
    "result": "./result.list",
    "source": "./source.list",
    "forbid_mixed_ip": true
}
```
