# whoisenum

Whoisenum is a tool for querying Regional Internet registrys
(ARIN, RIPE NCC, etc) for information they maintain with regards
to IPs or Domains. Whoisenum leverages the emerging RDAP protocol to query these entities APIs.

### Help Information
```bash
Usage:
  whoisenum [command]

Available Commands:
  domain      Search internet registries for a given domain
  help        Help about any command
  ip          Search internet registries by ip address

Flags:
  -f, --file string     File with single IP/Domain per line
  -h, --help            help for whoisenum
      --json            Output results in JSON
  -l, --lookup string   Single IP/Domain to lookup. Has no effect when --file is also specified.
  -q, --quiet           Quiet extra program output and only print results
  -r, --rate int        The number of seconds between queries
  -v, --verbose         Include verbose messages from program execution
      --version         version for whoisenum

Use "whoisenum [command] --help" for more information about a command.
```

### Example Usage

```bash
whoisenum ip -l <IP> --json
```
The above command will query the appropriate API to determine information about the IP address provided and return results in a json format.

### Todo
- Ban Detection (from over querying)
- Throttle excessive concurrent API requests to the same RIR API

### Contributers
https://github.com/brutalgg