# Vote to Mute
An experiment with [bwmarrin's discordgo library.](https://github.com/bwmarrin/discordgo)

## Building from source
1. Install [Golang](https://golang.org/dl)
2. Fetch dependencies and build (`go build`)

The binary `vote-to-mute` is then able to be executed with environment variables. **You will need to manually export and set these variables, as the .env file will not be read.** You can learn more about how to set environment variables on [Windows](https://docs.microsoft.com/powershell/module/microsoft.powershell.core/about/about_environment_variables), [macOS](https://support.apple.com/guide/terminal/apd382cc5fa-4f58-4449-b20a-41c53c006f8f), or [Linux](https://www.redhat.com/sysadmin/linux-environment-variables) with their respective links.

## Usage
This tool requires only one variables:
- `DISCORD_TOKEN`: This is required to connect to Discord's servers. You can fetch one by:
  1. Going to [Discord's developer portal](https://discord.com/developers)
  2. Creating a new application
  3. Going to the « Bot » tab and creating a new bot
  4. Copying the `Token` value

## FAQ
**Q: What permissions should the bot have?** The bot only needs `MANAGE_ROLES` (permission integer `268435456`). It needs:
  - `MANAGE_ROLES`: to create a `Muted` role, if one doesn't already exist, that mutes the user across channels
**Q: Is there a version of this bot already hosted?** I do not host a public bot, but some people do host instances for personal servers.