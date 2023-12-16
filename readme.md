# Google Home Notifier

## About

This program is designed for letting Google Home Speak text which was sent from Slack.

## Install

1. Clone this repository and compile it.
2. Download VOICEVOX Core and Open JTalk (UTF-8)
3. Move every files in voicevox_core to the GoogleHomeNotifier binary directory
4. Specify settings on `settings/settings.yaml`
5. Start

### Point

If you are using Linux, you have to add libvoicevox_core.so to Library path.

```bash
export LD_LIBRARY_PATH=/path/to/so/directory:$LD_LIBRARY_PATH
```
