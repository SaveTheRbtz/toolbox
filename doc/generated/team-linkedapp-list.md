# team linkedapp list 

List linked applications




# Security

`watermint toolbox` stores credentials into the file system. That is located at below path:

| OS       | Path                                                               |
| -------- | ------------------------------------------------------------------ |
| Windows  | `%HOMEPATH%\.toolbox\secrets` (e.g. C:\Users\bob\.toolbox\secrets) |
| macOS    | `$HOME/.toolbox/secrets` (e.g. /Users/bob/.toolbox/secrets)        |
| Linux    | `$HOME/.toolbox/secrets` (e.g. /home/bob/.toolbox/secrets)         |

Please do not share those files to anyone including Dropbox support.
You can delete those files after use if you want to remove it.
If you want to make sure removal of credentials, revoke application access from setting or the admin console.

Please see below help article for more detail:
* Dropbox Business: https://help.dropbox.com/ja-jp/teams-admins/admin/app-integrations

This command use following access type(s) during the operation:
* Dropbox Business File access


# Usage

This document uses the Desktop folder for command example. 

## Run

Windows:

```powershell
cd $HOME\Desktop
.\tbx.exe team linkedapp list 
```

macOS, Linux:

```bash
$HOME/Desktop/tbx team linkedapp list 
```



## Options

| Option  | Description   | Default   |
|---------|---------------|-----------|
| `-peer` | Account alias | {default} |


Common options:

| Option         | Description                                                                      | Default              |
|----------------|----------------------------------------------------------------------------------|----------------------|
| `-concurrency` | Maximum concurrency for running operation                                        | Number of processors |
| `-debug`       | Enable debug mode                                                                | false                |
| `-proxy`       | HTTP/HTTPS proxy (hostname:port)                                                 |                      |
| `-quiet`       | Suppress non-error messages, and make output readable by a machine (JSON format) | false                |
| `-secure`      | Do not store tokens into a file                                                  | false                |
| `-workspace`   | Workspace path                                                                   |                      |



## Authentication

For the first run, `toolbox` will ask you an authentication with your Dropbox account. 
Please copy the link and paste it into your browser. Then proceed to authorization.
After authorization, Dropbox will show you an authorization code.
Please copy that code and paste it to the `toolbox`.

```
watermint toolbox xx.x.xxx
© 2016-2019 Takayuki Okazaki
Licensed under open source licenses. Use the `license` command for more detail.

Testing network connection...
Done

1. Visit the URL for the auth dialog:

https://www.dropbox.com/oauth2/authorize?client_id=xxxxxxxxxxxxxxx&response_type=code&state=xxxxxxxx

2. Click 'Allow' (you might have to login first):
3. Copy the authorisation code:
Enter the authorisation code
```



# Result

Report file path will be displayed last line of the command line output.
If you missed command line output, please see path below.
[job-id] will be the date/time of the run. Please see the latest job-id.

| OS      | Path                                                                                                      |
| ------- | --------------------------------------------------------------------------------------------------------- |
| Windows | `%HOMEPATH%\.toolbox\jobs\[job-id]\reports` (e.g. C:\Users\bob\.toolbox\jobs\20190909-115959.597\reports) |
| macOS   | `$HOME/.toolbox/jobs/[job-id]/reports` (e.g. /Users/bob/.toolbox/jobs/20190909-115959.597/reports)        |
| Linux   | `$HOME/.toolbox/jobs/[job-id]/reports` (e.g. /home/bob/.toolbox/jobs/20190909-115959.597/reports)         |



## Report: linked_app 

Report files are generated in `linked_app.csv`, `linked_app.xlsx` and `linked_app.json` format.
In case of a report become large, report in `.xlsx` format will be split into several chunks
like `linked_app_0000.xlsx`, `linked_app_0001.xlsx`, `linked_app_0002.xlsx`...   

| Column           | Description                                                                          |
|------------------|--------------------------------------------------------------------------------------|
| team_member_id   | ID of user as a member of a team.                                                    |
| email            | Email address of user.                                                               |
| status           | The user's status as a member of a specific team. (active/invited/suspended/removed) |
| given_name       | Also known as a first name                                                           |
| surname          | Also known as a last name or family name.                                            |
| familiar_name    | Locale-dependent name                                                                |
| display_name     | A name that can be used directly to represent the name of a user's Dropbox account.  |
| abbreviated_name | An abbreviated form of the person's name.                                            |
| external_id      | External ID that a team can attach to the user.                                      |
| account_id       | A user's account identifier.                                                         |
| app_id           | The application unique id.                                                           |
| app_name         | The application name.                                                                |
| is_app_folder    | Whether the linked application uses a dedicated folder.                              |
| publisher        | The publisher's URL.                                                                 |
| publisher_url    | The application publisher name.                                                      |
| linked           | The time this application was linked                                                 |


