# discord-matic-ticker
Live matic tickers for your discord server.

This repo is a clone/redo of another project of mine. The documentation here might not be the best, and if you need further help please join the support server linked below.

[![Releases](https://github.com/rssnyder/discord-stock-ticker/workflows/Build%20and%20Publish%20Container%20Image/badge.svg)](https://github.com/rssnyder/discord-matic-ticker/releases)
[![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](https://github.com/tterb/atomic-design-ui/blob/master/LICENSEs)

[![GitHub last commit](https://img.shields.io/github/last-commit/rssnyder/discord-matic-ticker.svg?style=flat)]()
[![GitHub stars](https://img.shields.io/github/stars/rssnyder/discord-matic-ticker.svg?style=social&label=Star)]()
[![GitHub watchers](https://img.shields.io/github/watchers/rssnyder/discord-matic-ticker.svg?style=social&label=Watch)]()

## Preview

![Discord Sidebar w/ Bots](https://s3.cloud.rileysnyder.org/public/assets/sidebarmatic.png)

## Join the discord server
[![Discord Chat](https://logo.clearbit.com/discord.com)](https://discord.gg/CQqnCYEtG7)

## Support this project
<a href='https://ko-fi.com/rileysnyder' target='_blank'><img height='35' style='border:0px;height:46px;' src='https://az743702.vo.msecnd.net/cdn/kofi3.png?v=0' border='0' alt='Buy Me a Coffee' /></a>

Love these bots? You can support this project by subscribing to the [premium version](https://github.com/rssnyder/discord-stock-ticker/blob/master/README.md#premium) or maybe [buy me a coffee](https://ko-fi.com/rileysnyder) or [hire me](https://github.com/rssnyder) to write/host **your** discord bot!

## Related Projects

This project was created as a fork from my main [discord-stock-ticker](https://github.com/rssnyder/discord-stock-ticker) project

## Premium

![Discord Sidebar w/ Premium Bots](https://s3.cloud.rileysnyder.org/public/assets/sidebar-premium.png)

For advanced features like faster update times and color changing names on price changes you can subscribe to my premuim offering. I will host individual instances for your discord server at a cost of $1 per bot per month. You can choose a mix of cryptos and stocks and cancel at any time.

If you wish to host your bots on your own hardware, but need help getting set up, I also offer setup services for $20. I will install the service on your hardware and set you up with my internal tools to help manage your instances. This requires a running linux server.

If you are interested please see the [contact info on my github page](https://github.com/rssnyder) and send me a messgae via your platform of choice (discord perferred). For a live demo, join the support discord linked at the top or bottom of this page.

### Self-Hosting

#### Running in a simple shell

Pull down the latest release for your OS [here](https://github.com/rssnyder/discord-matic-ticker/releases).

```
wget https://github.com/rssnyder/discord-matic-ticker/releases/download/v0.1.0/discord-matic-ticker-v0.1.0-linux-amd64.tar.gz

tar zxf discord-matic-ticker-v2.0.0-linux-amd64.tar.gz

./discord-matic-ticker
```

### Systemd service

The below script (ran as root) will download and install a `discrod-matic-ticker` service on your linux machine with the API avalible on port `8080` to manage bots.

```
wget https://github.com/rssnyder/discord-matic-ticker/releases/download/v2.2.0/discord-matic-ticker-v2.2.0-linux-amd64.tar.gz

tar zxf discord-matic-ticker-v2.2.0-linux-amd64.tar.gz

mkdir -p /etc/discord-matic-ticker

mv discord-matic-ticker /etc/discord-matic-ticker/

wget https://raw.githubusercontent.com/rssnyder/discord-matic-ticker/master/discord-matic-ticker.service

mv discord-matic-ticker.service /etc/systemd/system/

systemctl daemon-reload

systemctl start discord-matic-ticker.service
```

### Adding multiple bots

A new feature in v2 is having one instance of the discord-matic-ticker manage multiple bots for different matics and cryptos.

To add another bot to your instance, you need to use the API exposed on port 8080 of the host you are running on:

#### List current running bots

```
curl localhost:8080/ticker
```

#### Add a new bot

Matic Payload: 

```
{
  "contract": "0x0000000000000000",
  "discord_bot_token": "xxxxxxxxxxxxxxxxxxxxxxxx",
  "name": "XDO",  # string/OPTIONAL: overwrites display name of bot
  "frequency": 10,  # int/OPTIONAL: default 60
  "set_nickname": true,  # bool/OPTIONAL
}
```

#### Remove a bot

```
curl -X DELETE localhost:8080/ticker/0x0000000000000000
```

## Support

If you have a request for a new ticker or issues with a current one, please open a github issue or find me on discord at `jonesbooned#1111` or [join the support server](https://discord.gg/CQqnCYEtG7).

Love these bots? Maybe [buy me a coffee](https://ko-fi.com/rileysnyder)! Or send some crypto to help keep these bots running:

eth: 0x27B6896cC68838bc8adE6407C8283a214ecD4ffE

doge: DTWkUvFakt12yUEssTbdCe2R7TepExBA2G

bch: qrnmprfh5e77lzdpalczdu839uhvrravlvfr5nwupr

btc: 1N84bLSVKPZBHKYjHp8QtvPgRJfRbtNKHQ
