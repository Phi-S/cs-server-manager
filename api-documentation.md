<!-- Generator: Widdershins v4.0.1 -->

<h1 id="github-com-phi-s-cs-server-manager-api">github.com/Phi-S/cs-server-manager API v1.0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Base URLs:

* <a href="/api/v1">/api/v1</a>

<h1 id="github-com-phi-s-cs-server-manager-api-server">server</h1>

## Send game-server command

`POST /command`

> Body parameter

```json
{
  "command": "string"
}
```

<h3 id="send-game-server-command-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[handlers.CommandRequest](#schemahandlers.commandrequest)|true|This command will be executed on the game server|
|» command|body|string|true|none|

> Example responses

> 200 Response

```
"string"
```

<h3 id="send-game-server-command-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|string|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Start the server

`POST /start`

Starts the server with the given start parameters

> Body parameter

```json
{
  "hostname": "string",
  "max_players": 128,
  "password": "string",
  "start_map": "string",
  "steam_login_token": "string"
}
```

<h3 id="start-the-server-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[handlers.StartBody](#schemahandlers.startbody)|false|You can provide no, all or only a few start parameters. The provided start parameters will overwrite the saved start parameters in the start-parameters.json file if the server started successfully.|
|» hostname|body|string|false|none|
|» max_players|body|integer|false|none|
|» password|body|string|false|none|
|» start_map|body|string|false|none|
|» steam_login_token|body|string|false|none|

> Example responses

> 400 Response

<h3 id="start-the-server-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Get the current status of the server

`GET /status`

> Example responses

> 200 Response

```json
{
  "hostname": "string",
  "ip": "string",
  "is_game_server_installed": true,
  "map": "string",
  "max_player_count": 0,
  "password": "string",
  "player_count": 0,
  "port": "string",
  "state": "idle"
}
```

<h3 id="get-the-current-status-of-the-server-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[status.InternalStatus](#schemastatus.internalstatus)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Stop the server

`POST /stop`

Stops the server of if the server is not running, returns 200 OK

> Example responses

> 400 Response

<h3 id="stop-the-server-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="github-com-phi-s-cs-server-manager-api-files">files</h1>

## Get editable files

`GET /files`

> Example responses

> 200 Response

```json
[
  {
    "files": [
      "string"
    ]
  }
]
```

<h3 id="get-editable-files-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<h3 id="get-editable-files-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[handlers.FilesResponse](#schemahandlers.filesresponse)]|false|none|none|
|» files|[string]|false|none|none|

<aside class="success">
This operation does not require authentication
</aside>

## Get files content

`GET /files/{file}`

<h3 id="get-files-content-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|file|path|string|true|file to get content for|

> Example responses

> 200 Response

```
"string"
```

<h3 id="get-files-content-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|string|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Set files content

`PATCH /files{file}`

> Body parameter

```
string

```

<h3 id="set-files-content-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|file|path|string|true|file to set content for|
|body|body|string|true|file content|

> Example responses

> 400 Response

<h3 id="set-files-content-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="github-com-phi-s-cs-server-manager-api-logs">logs</h1>

## Get logs

`GET /logs/{count}`

<h3 id="get-logs-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|count|path|integer|true|Get the last X logs|

> Example responses

> 200 Response

```json
[
  {
    "log_type": "string",
    "message": "string",
    "timestamp": "string"
  }
]
```

<h3 id="get-logs-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<h3 id="get-logs-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[logwrt.LogEntry](#schemalogwrt.logentry)]|false|none|none|
|» log_type|string|false|none|none|
|» message|string|false|none|none|
|» timestamp|string|false|none|none|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="github-com-phi-s-cs-server-manager-api-plugins">plugins</h1>

## Get all available plugins

`GET /plugins`

> Example responses

> 200 Response

```json
{
  "description": "string",
  "name": "string",
  "url": "string",
  "versions": [
    {
      "dependencies": [
        {
          "dependencies": [
            {}
          ],
          "download_url": "string",
          "install_dir": "string",
          "name": "string",
          "version": "string"
        }
      ],
      "installed": true,
      "name": "string"
    }
  ]
}
```

<h3 id="get-all-available-plugins-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[handlers.PluginResponse](#schemahandlers.pluginresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Install given plugin

`POST /plugins`

> Body parameter

```json
{
  "name": "string",
  "version": "string"
}
```

<h3 id="install-given-plugin-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[handlers.InstallPluginRequest](#schemahandlers.installpluginrequest)|true|The plugin and version that should be installed|
|» name|body|string|false|none|
|» version|body|string|false|none|

> Example responses

> 400 Response

<h3 id="install-given-plugin-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Uninstall the currently installed plugin

`DELETE /plugins`

> Example responses

> 400 Response

<h3 id="uninstall-the-currently-installed-plugin-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="github-com-phi-s-cs-server-manager-api-settings">settings</h1>

## Get the current settings

`GET /settings`

> Example responses

> 200 Response

```json
{
  "hostname": "string",
  "max_players": 128,
  "password": "string",
  "start_map": "string",
  "steam_login_token": "string"
}
```

<h3 id="get-the-current-settings-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[handlers.SettingsModel](#schemahandlers.settingsmodel)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Update settings

`POST /settings`

> Body parameter

```json
{
  "hostname": "string",
  "max_players": 128,
  "password": "string",
  "start_map": "string",
  "steam_login_token": "string"
}
```

<h3 id="update-settings-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[handlers.SettingsModel](#schemahandlers.settingsmodel)|true|The updated settings|
|» hostname|body|string|true|none|
|» max_players|body|integer|true|none|
|» password|body|string|false|none|
|» start_map|body|string|true|none|
|» steam_login_token|body|string|false|none|

> Example responses

> 200 Response

```json
{
  "hostname": "string",
  "max_players": 128,
  "password": "string",
  "start_map": "string",
  "steam_login_token": "string"
}
```

<h3 id="update-settings-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[handlers.SettingsModel](#schemahandlers.settingsmodel)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="github-com-phi-s-cs-server-manager-api-update">update</h1>

## Start server update

`POST /update`

> Example responses

> 400 Response

<h3 id="start-server-update-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Cancel the server update

`POST /update/cancel`

Cancel the currently running server update or if no update is currently running, returns 200 OK

> Example responses

> 400 Response

<h3 id="cancel-the-server-update-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_handlers.CommandRequest">handlers.CommandRequest</h2>
<!-- backwards compatibility -->
<a id="schemahandlers.commandrequest"></a>
<a id="schema_handlers.CommandRequest"></a>
<a id="tocShandlers.commandrequest"></a>
<a id="tocshandlers.commandrequest"></a>

```json
{
  "command": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|command|string|true|none|none|

<h2 id="tocS_handlers.ErrorResponse">handlers.ErrorResponse</h2>
<!-- backwards compatibility -->
<a id="schemahandlers.errorresponse"></a>
<a id="schema_handlers.ErrorResponse"></a>
<a id="tocShandlers.errorresponse"></a>
<a id="tocshandlers.errorresponse"></a>

```json
{
  "message": "string",
  "request_id": "string",
  "status": 0
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|message|string|false|none|none|
|request_id|string|false|none|none|
|status|integer|false|none|none|

<h2 id="tocS_handlers.FilesResponse">handlers.FilesResponse</h2>
<!-- backwards compatibility -->
<a id="schemahandlers.filesresponse"></a>
<a id="schema_handlers.FilesResponse"></a>
<a id="tocShandlers.filesresponse"></a>
<a id="tocshandlers.filesresponse"></a>

```json
{
  "files": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|files|[string]|false|none|none|

<h2 id="tocS_handlers.InstallPluginRequest">handlers.InstallPluginRequest</h2>
<!-- backwards compatibility -->
<a id="schemahandlers.installpluginrequest"></a>
<a id="schema_handlers.InstallPluginRequest"></a>
<a id="tocShandlers.installpluginrequest"></a>
<a id="tocshandlers.installpluginrequest"></a>

```json
{
  "name": "string",
  "version": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|false|none|none|
|version|string|false|none|none|

<h2 id="tocS_handlers.PluginDependencyResponse">handlers.PluginDependencyResponse</h2>
<!-- backwards compatibility -->
<a id="schemahandlers.plugindependencyresponse"></a>
<a id="schema_handlers.PluginDependencyResponse"></a>
<a id="tocShandlers.plugindependencyresponse"></a>
<a id="tocshandlers.plugindependencyresponse"></a>

```json
{
  "dependencies": [
    {
      "dependencies": [],
      "download_url": "string",
      "install_dir": "string",
      "name": "string",
      "version": "string"
    }
  ],
  "download_url": "string",
  "install_dir": "string",
  "name": "string",
  "version": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|dependencies|[[handlers.PluginDependencyResponse](#schemahandlers.plugindependencyresponse)]|false|none|none|
|download_url|string|false|none|none|
|install_dir|string|false|none|none|
|name|string|false|none|none|
|version|string|false|none|none|

<h2 id="tocS_handlers.PluginResponse">handlers.PluginResponse</h2>
<!-- backwards compatibility -->
<a id="schemahandlers.pluginresponse"></a>
<a id="schema_handlers.PluginResponse"></a>
<a id="tocShandlers.pluginresponse"></a>
<a id="tocshandlers.pluginresponse"></a>

```json
{
  "description": "string",
  "name": "string",
  "url": "string",
  "versions": [
    {
      "dependencies": [
        {
          "dependencies": [
            {}
          ],
          "download_url": "string",
          "install_dir": "string",
          "name": "string",
          "version": "string"
        }
      ],
      "installed": true,
      "name": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|description|string|false|none|none|
|name|string|false|none|none|
|url|string|false|none|none|
|versions|[[handlers.PluginVersionResponse](#schemahandlers.pluginversionresponse)]|false|none|none|

<h2 id="tocS_handlers.PluginVersionResponse">handlers.PluginVersionResponse</h2>
<!-- backwards compatibility -->
<a id="schemahandlers.pluginversionresponse"></a>
<a id="schema_handlers.PluginVersionResponse"></a>
<a id="tocShandlers.pluginversionresponse"></a>
<a id="tocshandlers.pluginversionresponse"></a>

```json
{
  "dependencies": [
    {
      "dependencies": [
        {}
      ],
      "download_url": "string",
      "install_dir": "string",
      "name": "string",
      "version": "string"
    }
  ],
  "installed": true,
  "name": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|dependencies|[[handlers.PluginDependencyResponse](#schemahandlers.plugindependencyresponse)]|false|none|none|
|installed|boolean|false|none|none|
|name|string|false|none|none|

<h2 id="tocS_handlers.SettingsModel">handlers.SettingsModel</h2>
<!-- backwards compatibility -->
<a id="schemahandlers.settingsmodel"></a>
<a id="schema_handlers.SettingsModel"></a>
<a id="tocShandlers.settingsmodel"></a>
<a id="tocshandlers.settingsmodel"></a>

```json
{
  "hostname": "string",
  "max_players": 128,
  "password": "string",
  "start_map": "string",
  "steam_login_token": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|hostname|string|true|none|none|
|max_players|integer|true|none|none|
|password|string|false|none|none|
|start_map|string|true|none|none|
|steam_login_token|string|false|none|none|

<h2 id="tocS_handlers.StartBody">handlers.StartBody</h2>
<!-- backwards compatibility -->
<a id="schemahandlers.startbody"></a>
<a id="schema_handlers.StartBody"></a>
<a id="tocShandlers.startbody"></a>
<a id="tocshandlers.startbody"></a>

```json
{
  "hostname": "string",
  "max_players": 128,
  "password": "string",
  "start_map": "string",
  "steam_login_token": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|hostname|string|false|none|none|
|max_players|integer|false|none|none|
|password|string|false|none|none|
|start_map|string|false|none|none|
|steam_login_token|string|false|none|none|

<h2 id="tocS_logwrt.LogEntry">logwrt.LogEntry</h2>
<!-- backwards compatibility -->
<a id="schemalogwrt.logentry"></a>
<a id="schema_logwrt.LogEntry"></a>
<a id="tocSlogwrt.logentry"></a>
<a id="tocslogwrt.logentry"></a>

```json
{
  "log_type": "string",
  "message": "string",
  "timestamp": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|log_type|string|false|none|none|
|message|string|false|none|none|
|timestamp|string|false|none|none|

<h2 id="tocS_status.InternalStatus">status.InternalStatus</h2>
<!-- backwards compatibility -->
<a id="schemastatus.internalstatus"></a>
<a id="schema_status.InternalStatus"></a>
<a id="tocSstatus.internalstatus"></a>
<a id="tocsstatus.internalstatus"></a>

```json
{
  "hostname": "string",
  "ip": "string",
  "is_game_server_installed": true,
  "map": "string",
  "max_player_count": 0,
  "password": "string",
  "player_count": 0,
  "port": "string",
  "state": "idle"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|hostname|string|false|none|none|
|ip|string|false|none|none|
|is_game_server_installed|boolean|false|none|none|
|map|string|false|none|none|
|max_player_count|integer|false|none|none|
|password|string|false|none|none|
|player_count|integer|false|none|none|
|port|string|false|none|none|
|state|[status.State](#schemastatus.state)|false|none|none|

<h2 id="tocS_status.State">status.State</h2>
<!-- backwards compatibility -->
<a id="schemastatus.state"></a>
<a id="schema_status.State"></a>
<a id="tocSstatus.state"></a>
<a id="tocsstatus.state"></a>

```json
"idle"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|idle|
|*anonymous*|server-starting|
|*anonymous*|server-started|
|*anonymous*|server-stopping|
|*anonymous*|steamcmd-updating|
|*anonymous*|plugin-installing|
|*anonymous*|plugin-uninstalling|

