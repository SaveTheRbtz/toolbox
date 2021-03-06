# dev benchmark local

Create dummy folder structure in local file system. 

# Usage

This document uses the Desktop folder for command example.
## Run

Windows:
```
cd $HOME\Desktop
.\tbx.exe dev benchmark local -path /LOCAL/PATH/TO/PROCESS
```

macOS, Linux:
```
$HOME/Desktop/tbx dev benchmark local -path /LOCAL/PATH/TO/PROCESS
```

Note for macOS Catalina 10.15 or above: macOS verifies Developer identity. Currently, `tbx` is not ready for it. Please select "Cancel" on the first dialogue. Then please proceed "System Preference", then open "Security & Privacy", select "General" tab.
You may find the message like:
> "tbx" was blocked from use because it is not from an identified developer.

And you may find the button "Allow Anyway". Please hit the button with your risk. At second run, please hit button "Open" on the dialogue.

## Options:

| Option         | Description                | Default |
|----------------|----------------------------|---------|
| `-node-lambda` | Lambda parameter for nodes | 100     |
| `-node-max`    | Maximum number of nodes    | 1000    |
| `-node-min`    | Minimum number of nodes    | 100     |
| `-path`        | Path to create             |         |
| `-size-max`    | Maximum file size          | 2097152 |
| `-size-min`    | Minimum file size          | 0       |

## Common options:

| Option            | Description                                                                      | Default              |
|-------------------|----------------------------------------------------------------------------------|----------------------|
| `-auto-open`      | Auto open URL or artifact folder                                                 | false                |
| `-bandwidth-kb`   | Bandwidth limit in K bytes per sec for upload/download content. 0 for unlimited  | 0                    |
| `-budget-memory`  | Memory budget (limits some feature to reduce memory footprint)                   | normal               |
| `-budget-storage` | Storage budget (limits logs or some feature to reduce storage usage)             | normal               |
| `-concurrency`    | Maximum concurrency for running operation                                        | Number of processors |
| `-debug`          | Enable debug mode                                                                | false                |
| `-experiment`     | Enable experimental feature(s).                                                  |                      |
| `-lang`           | Display language                                                                 | auto                 |
| `-output`         | Output format (none/text/markdown/json)                                          | text                 |
| `-proxy`          | HTTP/HTTPS proxy (hostname:port)                                                 |                      |
| `-quiet`          | Suppress non-error messages, and make output readable by a machine (JSON format) | false                |
| `-secure`         | Do not store tokens into a file                                                  | false                |
| `-workspace`      | Workspace path                                                                   |                      |

# Proxy configuration

The executable automatically detects your proxy configuration from the environment. However, if you got an error or you want to specify explicitly, please add -proxy option, like -proxy hostname:port. Currently, the executable doesn't support proxies which require authentication.

