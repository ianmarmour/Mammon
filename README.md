## Introduction


## Development Setup

The only required setup and installation should be to set your Blizzard API secrets up for docker using hte following commands.

```
echo "mysupersecureclientid" | docker secret create BLIZZARD_API_CLIENT_ID -
echo "mysupersecureclientsecret" | docker secret create BLIZZARD_API_CLIENT_SECRET -
```