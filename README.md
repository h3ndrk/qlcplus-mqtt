# qlcplus-mqtt

## QLC+ WebSocket protocol reverse engineering

- (empty): ignore
- `QLC+CMD|opMode`: Toggle between design/operate mode
- `QLC+IO|(INPUT|OUTPUT|FB)|<universe>|<pluginName>|<input>`: Patch the given universe to go through the given input/output plugin (and output with feedback)
    - `<universe>`: The input universe to patch (starting at zero)
    - `<pluginName>`: The name of the plugin to patch to the universe
    - `<input>`: An input universe provided by the plugin to patch to
- `QLC+IO|PROFILE|<universe>|<profileName>`: Patch the given universe via a given profile
    - `<universe>`: The input universe to patch (starting at zero)
    - `<profileName>`: The name of an input profile
- `QLC+IO|PASSTHROUGH|<universe>|<state>`: Set/unset the universe with the given index in passthrough mode
    - `<universe>`: The universe index (starting at zero)
    - `<state>`: true = passthrough, false = normal mode
- `QLC+IO|AUDIOIN|<device>`: Set input capture device
    - `<device>`: The device or `__qlcplusdefault__` to disable
- `QLC+IO|AUDIOOUT|<device>`: Set output capture device
    - `<device>`: The device or `__qlcplusdefault__` to disable
- `QLC+AUTH|ADD_USER|<username>|<password>|<level>`: Adds user to password table. If given username already exists it is replaced.
    - `<username>`: The username to add
    - `<password>`: The password to add the user with
    - `<level>`: GUEST_LEVEL = 0, LOGGED_IN_LEVEL = 1, VC_ONLY_LEVEL = 10, SIMPLE_DESK_AND_VC_LEVEL = 20, SUPER_ADMIN_LEVEL = 100, NOT_PROVIDED_LEVEL = 100
- `QLC+AUTH|DEL_USER|<username>`: Deletes a user from password table
    - `<username>`: The username to delete
- `QLC+AUTH|SET_USER_LEVEL|<username>|<level>`: Set user level for a user in password table
    - `<username>`: The username to delete
    - `<level>`: GUEST_LEVEL = 0, LOGGED_IN_LEVEL = 1, VC_ONLY_LEVEL = 10, SIMPLE_DESK_AND_VC_LEVEL = 20, SUPER_ADMIN_LEVEL = 100, NOT_PROVIDED_LEVEL = 100
- `QLC+SYS|(NETWORK|AUTOSTART|REBOOT|HALT)`: System functionalities (skipped because not the scope)
- `QLC+API|isProjectLoaded`: Query and clear/disable project loaded state
    - answer: `QLC+API|isProjectLoaded|<state>` with `<state>` equals to `true` or `false`
- `QLC+API|getFunctionsNumber`: Returns the number of functions
    - answer: `QLC+API|getFunctionsNumber|<number>`
- `QLC+API|getFunctionsList`: Returns all functions
    - answer: `QLC+API|getFunctionsList|[<id>|<name>]...`
- `QLC+API|getFunctionType|<id>`: Returns the type of the function with the given `<id>` (or `Undefined`)
    - answer: `QLC+API|getFunctionType|<type>` with `<type>` one of `Scene`, `Chaser`, `EFX`, `Collection`, `Script`, `RGBMatrix`, `Show`, `Sequence`, `Audio`, `Video`
- `QLC+API|getFunctionStatus|<id>`: Returns the state of the function with the given `<id>`
    - answer: `QLC+API|getFunctionStatus|<state>` with `<state>` equals to `Running` or `Stopped` (or `Undefined`)
- `QLC+API|setFunctionStatus|<id>|<state>`: Sets the state of the function with the given `<id>` to the given `<state>` (`0` is `Stopped`, `1` is `Running`)
- `QLC+API|getWidgetsNumber`: Returns the number of widgets
    - answer: `QLC+API|getWidgetsNumber|<number>`
- `QLC+API|getWidgetsList`: Returns all widgets
    - answer: `QLC+API|getWidgetsList|[<id>|<caption>]...`
- `QLC+API|getWidgetType|<id>`: Returns the type of the function with the given `<id>` (or `Unknown`)
    - answer: `QLC+API|getWidgetType|<type>` with `<type>` one of `Button`, `Slider`, `XYPad`, `Frame`, `SoloFrame`, `SpeedDial`, `CueList`, `Label`, `AudioTriggers`, `Animation`, `Clock`
- `QLC+API|getWidgetStatus|<id>`: Returns the state of the function with the given `<id>`
    - answer for type `Button`: `QLC+API|setWidgetStatus|<value>` with `<value>` equals to `255` when active, `127` when monitoring or `0` when inactive
    - answer for type `Slider`: `QLC+API|setWidgetStatus|<value>` with `<value>` in the range `[0,255]`
    - answer for type playing `CueList`: `QLC+API|setWidgetStatus|PLAY|<cueIndex>` with `<cueIndex>` equals to the current cue index
    - answer for type stopped `CueList`: `QLC+API|setWidgetStatus|STOP`
- `QLC+API|getChannelsValues|<universe>|<startAddr>|[<count>]`: Returns the channel values of a given `<universe>` (starting at one) beginning at given `<startAddr>` (starting at one) and returning one or given `<count>` channels
    - answer: `QLC+API|getChannelsValues|[<channel>|<value>|<type>]...` with `<channel>` starting at one, value in range `[0,255]` and `<type>` the fixture type
- `QLC+API|sdResetChannel|<channel>`: Reset given `<channel>` of the current universe (skipped because not the scope)
- `QLC+API|sdResetUniverse`: Reset all channels of the current universe (skipped because not the scope)
- `CH|<absAddr>|<value>`: Sets the channel described with given `<absAddr>` to given `<value>` in range `[0,255]`
- `POLL`: Poll this API
- `<id>[|<value>]`: Sets given `<value>` of widget with given `<id>`
    - `Button`: `<value>` equals to `1` for press and `0` for release
    - `Slider`: `<value>` is in range of slider
    - `AudioTrigger`: `<value>` equals to `1` for active and `0` for inactive
    - `CueList`: `<value>` is `<command>[|<index>]`:
        - `PLAY`: Plays cue
        - `STOP`: Stops cue
        - `PREV`: Selects previous cue
        - `NEXT`: Selects next cue
        - `STEP|<index>`: Selects cue with given `<index>`
    - `Frame`, `SoloFrame`: `<value>` is one of `NEXT_PG` (next page) or `PREV_PG` (previous page)

## MQTT topics and payloads

The following topics are subscribed/published to by qlcplus-mqtt.

### Subscribe to topics

- `/qlcplus`
    - `/opmode`: Toggle between design/operate mode
    - `/io`
        - `/input`: `{"universe": number, "pluginName": string, "input": number}`
        - `/output`: `{"universe": number, "pluginName": string, "input": number}`
        - `/feedback`: `{"universe": number, "pluginName": string, "input": number}`
        - `/profile`: `{"universe": number, "profileName": string}`
        - `/passthrough`: `{"universe": number, "enable": boolean}`
        - `/audioin`: `{"device": string/null}`
        - `/audioout`: `{"device": string/null}`
    - `/auth`
        - `/add`: `{"username": string, "password": string, "level": number}`
        - `/delete`: `{"username": string}`
        - `/modify`: `{"username": string, "level": number}`
    - `/api`
        - `/function`
            - `/<id>`
                - `/enable`: `boolean`
        - `/channel`
            - `/<absAddr>`
                - `/value`: `number`
        - `/widget`
            - `/<id>`
                - `/value`: `boolean` for `Button` and `AudioTrigger` type, `number` for `Slider` type, `{"command": string, ["index": number]}` for `CueList` type, `string` for `Frame` and `SoloFrame` type

### Publish to topics

- `/qlcplus`
    - `/api`
        - `/function`: `[{"id": number, "type": string, "enabled": boolean}]` (retained)
        - `/channel`: `[{"absAddr": number, "value": number}]` (retained)
        - `/widget`: `[{"id": number, "type": string, "value": <value>}]` with value see above (retained)
