# Google Home Notifier

## About

This program is designed for letting Google Home Speak text which was sent from Slack.

## Install

1. Download VOICEVOX Core and Open JTalk dict (UTF-8)
2. Clone this repository and compile it.
3. Move every files in voicevox_core to the GoogleHomeNotifier binary directory
4. Specify settings on `settings/settings.yaml`
5. Start

☆ If you are using Linux, you have to add libvoicevox_core.so to Library path.

```bash
export LD_LIBRARY_PATH=/path/to/so/directory:$LD_LIBRARY_PATH
```


## Settings

```yaml
GoogleHome:
  Addr: # Google Home IP address
  Port: 8009 # GoogleHome port number 
  Detach: true  # Optional
  ForceDetach: true # Optional

Voicevox:
  SpeakerID: 3 
  OpenJtalkDictDir: "open_jtalk_dic_utf_8-1.11" # You have to specify Open JTalk's dict path

Slack:
  Token: # Slack bot token, which has permissions of app_mentions:read, chat:write, and users:read, (and optinally, chat:write.customize)
  AppLevelToken: # Slack App level token, which has a scopeof connections:write.
  Icon: # (optional) icon emoji You can use this if you add chat:write.customize permission.
```
