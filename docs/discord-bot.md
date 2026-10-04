# Debugging the Discord bot

The bot runs inside the `cmd/server` process (`pkg/discord`). It isn't a separate service, so when the bot goes offline the website usually stays up.

## Where it runs

- **Host:** Linode, `ssh casd`
- **Service:** `ddstats-server.service`, unit file at `/lib/systemd/system/ddstats-server.service`
- **Binary:** `/home/alex/ddstats/server`, deployed by `.github/workflows/release.yml`
- **Bot token:** passed as `-discord-token=...` in the unit's `ExecStart`. Don't paste the unit file anywhere public, because it also contains the DB password.
- **Discord application:** <https://discord.com/developers/applications>, bot user `ddstatsv2`

## Read the logs

The server logs to stdout and stderr, which go to the systemd journal. Reading the journal needs `sudo`.

```bash
ssh casd
sudo journalctl -u ddstats-server -f                    # follow live
sudo journalctl -u ddstats-server --since "1 hour ago"
sudo journalctl -u ddstats-server | grep -i discord     # bot lines only
```

What the bot logs:

| Log line | Meaning |
| --- | --- |
| `Discord bot connected as ddstatsv2 to N guild(s)` | Gateway login succeeded. Printed on startup and after every full reconnect. |
| `Discord bot broadcasting to N ddstats channel(s)` | Number of channels that alerts go to. `0` means no channel name contains `ddstats`. |
| `Discord bot resumed gateway session` | Short drop, recovered automatically. |
| `Discord bot disconnected from gateway` | Connection lost. Look for a `connected` or `resumed` line right after it. |
| `discordgo: ...` | The library's own warnings and errors, such as failed reconnects and rate limits. |
| `discord command .rank args=["1"] from user (id) in guild G channel C took 210ms` | Every command the bot handled. |
| `discord command .x from ... rejected: on cooldown` | Command was ignored because of its cooldown. |
| `sending .x reply to ...: <err>` | Discord rejected the reply, usually because of a missing channel permission. |
| `panic handling discord message ...` | A command crashed. The log includes a stack trace. |
| `Discord bot failed to start, continuing without it: ...` | The website is up but the bot never connected. Check the error text. |

## Bot shows offline: checklist

Work through these in order.

1. **Is the process running?** Run `systemctl status ddstats-server`. If it isn't running, the site is down too, so start there.
2. **Does the process have a gateway connection?** Run `ss -tn | grep -E '162\.159\.'`. No connection to a Discord IP means the gateway is disconnected.
3. **Check the logs** for `discordgo:` and `Discord bot` lines, using the table above.
4. **Is the token still valid?** Run this on `casd`:

   ```bash
   T=$(systemctl cat ddstats-server | grep -o -- '-discord-token=[^ ]*' | cut -d= -f2-)
   curl -s -o /dev/null -w '%{http_code}\n' -H "Authorization: Bot $T" https://discord.com/api/v10/users/@me
   ```

   `200` means the token is fine. `401` means the token was reset in the Developer Portal: update the unit file, then run `sudo systemctl daemon-reload && sudo systemctl restart ddstats-server`.
5. **Match the gateway close code** in the logs (`websocket: close 40xx`):
   - `4004`: the token is invalid. See step 4.
   - `4013`: invalid intents.
   - `4014`: disallowed intents. The bot requests the privileged **Message Content** intent, which must be enabled in the Developer Portal under *Bot → Privileged Gateway Intents*. Without it, the bot can't read `.commands`.
6. **Are commands ignored while the bot shows online?** Either the Message Content intent is missing, or the bot lacks *Send Messages* or *Embed Links* in that channel. Look for `sending ... reply` errors.

## Restart

```bash
sudo systemctl restart ddstats-server
```

A restart also briefly takes down the website and the live socket.io and websocket feeds, so avoid it during busy hours. A Discord startup failure is logged and doesn't stop the server.

## Run the bot locally

Use a separate test bot application. If you run a second process with the production token, the two gateway sessions fight over the bot's connection.

```bash
go run ./cmd/server --discord-token "<test bot token>" --dsn "host=localhost ..."
```

## Known failure: outdated discordgo

In October 2026 the bot went offline while the server stayed up. discordgo v0.20.2 parsed Discord's `X-RateLimit-Reset` header as an integer. Discord now sends a float (`1791055029.503`), so every reconnect attempt failed with `strconv.ParseInt: parsing "...": invalid syntax`. Before this change, discordgo's logs weren't wired into the server logs, so nothing showed the failure. If the bot fails in a way that suggests the Discord API changed, check for a newer discordgo first: `go list -m -u github.com/bwmarrin/discordgo`.
