# go-glowMarkt
## Reads electricity/gas usage data from glowMarkt API from your SMETS1 or SMETS2 and injects into influxDB
For real-time data you need Glowmarkt compatible hardware - which works with either smart or traditional meters. People with SMETS2 (or E&A SMETS 1) meters can access their half hourly data via Bright for free.
Sign up using one of the links below and this will give you access to the API which is used to inject into influxDb

## Android App
https://play.google.com/store/apps/details?id=uk.co.hildebrand.brightionic&hl=en_GB&gl=US

## IOS App
https://apps.apple.com/gb/app/bright/id1369989022

## Compiling Docker Image
docker build --tag goglowmarkt .

## Docker usage
docker run -e glowUsername="your glow username" -e glowPassword='your glow password' -e influxDbOrg="InfluxDb Org" -e influxDbBucker="influxDb bucket" -e influxDbToken="influxDB Token" -e influxDbUrl="url of your InfluxDb Instance including port" -e defaultInterval="5" goglowmarkt

## Usage
By default go-glowmarkt will run every 30 minutes, unless you specify the **defaultInterval** variable in which case you can enter any period you like with a minimum value of 5 minutes.  NB to retrieve data of less than 30 minute intervals you must purchase a device from https://shop.glowmarkt.com/

You can alter whether to send a catchup by sending a post to the endpoint as per below:

GET : http://<yourcontainer>:8888/catchup    ~   Gets the current catchup status
POST: http://<yourcontainer>:8888/catchup    ~   Post the below body with either true or false
    {
        "setCatchup": true
    }

You can retrieve the last x days worth of data ( and populate InluxDb by using the following endpoint and body)
POST: http://<yourcontainer>:8888/metrics       
{
    "numberOfDays": 3
}