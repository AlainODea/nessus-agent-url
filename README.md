# nessus-agent-url
Get the URL of the latest matching Nessus Agent Download.

## Motivation
Tenable has no "latest" URL for downloading a Nessus Agent installer. They also somewhat frequently drop previously available downloads from their site.

## Usage
Get the URL for the latest Amazon Linux / Amazon Linux 2 agent RPM installer for x86_64:
```shell
./nessus-agent-url '-amzn.x86_64.rpm'
```

Returns something like this:
```shell
https://www.tenable.com/downloads/api/v1/public/pages/nessus-agents/downloads/16733/download?i_agree_to_tenable_license_agreement=true
```

Scripted download in Bash:
```bash
curl --location --remote-name --remote-header-name "$(./nessus-agent-url '-amzn.x86_64.rpm')"
```

Or more compactly if you prefer:
```bash
curl -LOJ "$(./nessus-agent-url '-amzn.x86_64.rpm')"
```

## Build
For your platform only:
```shell
go build .
```

For a configurable selection of platforms in the Makefile:
```shell
make compile
```
