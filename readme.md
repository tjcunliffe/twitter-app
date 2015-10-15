# Twitter Client demo for Mirage

This is an example client application implementation for service virtualization tool [Mirage](https://github.com/SpectoLabs/mirage). 
It can either call Twitter API directly and by calling [Mirage Twitter Proxy](https://github.com/SpectoLabs/twitter-proxy)
(which provides recording and playback functionality). It expects to find proxy on port 8300 (localhost).
  
 
## Installation

* glide.yaml file holds all required dependencies. You can use glide to initialize your environment (glide install)  
* build a binary file (go build).
* Register your application in Twitter to get key and secret.

* Preferred way of storing Twitter auth details is environment variables:
  export TwitterKey=your_key
  export TwitterSecret=your_secret
  
  You can also create a configuration file for this:
  Rename _conf.json.example_ to _conf.json_ and provide your previously acquired key and secret.

* Run application (./twitter-app

## Configuration

To override default proxy location (localhost:8300) you can export environment variable "MirageProxyAddress":
export MirageProxyAddress=http://proxyhost:8300

To override default port (8080) supply flag during startup:
./twitter-app -port=":8888"