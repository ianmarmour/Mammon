## Introduction


## Development Setup

The only required setup and installation should be to set your Blizzard API secrets up for docker using hte following commands.

```
sudo docker swarm init

sudo docker build -t mammon .

echo "mysupersecureclientid" | docker secret create BLIZZARD_API_CLIENT_ID -
echo "mysupersecureclientsecret" | docker secret create BLIZZARD_API_CLIENT_SECRET -

sudo docker stack deploy -c docker-compose.yml mammo

// Track the service logs to make sure execution started properly
sudo docker service logs mammon_mammon
```