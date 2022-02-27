# go-glowMarkt
## Reads data from glowMarkt API and injects into influxDB

## Compiling Docker Image
docker build --tag goglowmarkt .
## Docker usage
docker run -e glowUsername="your glow username" -e glowPassword='your glow password' -e influxDbOrg="InfluxDb Org" -e influxDbBucker="influxDb bucket" -e influxDbToken="influxDB Token" -e influxDbUrl="url of your InfluxDb Instance including port" goglowmarkt
