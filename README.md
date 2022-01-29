# Korero: Send, recieve and manage messages on multiple platforms
Korero is a terminal based messaging service, being built to interact with popular direct and group messaging services like discord, telegram, slack, and twitter. 

The motivation for this project comes primarily from the intensive and bloating nature of messaging services, especially for lightweight machines. 

**How to install Korero**
- 
- Korero is supported on Fedora 34, 45 and Ubuntu. 
- For Fedora, it can be installed via dnf/yum by running: `sudo dnf copr enable cdoern/korero` followed by `sudo dnf install korero`.
- For Ubuntu, you can download and install the .deb package directly from the repository, [korero_0.1-1_all.deb](releases/korero_0.1-1_all.deb)
- Korero can be built from source by cloning this repository and running `make korero` within the project folder. The Korero binary will then be placed in `bin/korero`.

**Usage**
-
Korero is a command like tool and can be invoked by simply running `korero` this will show you the available commands:

- `korero discord setup` -- a step by step setup process for creating a bot and getting its token
- `korero discord login` -- test your current configuration by logging into the bot
- `korero discord messages` -- stream all messages from servers that the bot is in.

Once getting your token and creating your discord bot, you can set the environmental variable `$KORERO_DISCORD_TOKEN` so the program can automatically log in.

To permanently set environmental variables on Fedora `cd ~` and then `nano .bashrc` inside of .bashrc add the following line:
`export KORERO_DISCORD_TOKEN=<YOUR_TOKEN>`

**Contributions** 
- 
Contributions to Korero are welcome on the github repository. 

- If something in the app is broken please file an issue and place in the title **[BUG]** 
- If there is an enhancement that you would like to be made please file a bug and put **[FEATURE]** in the title. 

**Future Plans**
- 
Korero currently supports Discord messaging, the future development will be aimed at services like:

- Twitter
- Reddit
- Telegram
- Whatsapp
- Slack

This list has no preference level and will be worked on as community interest grows.