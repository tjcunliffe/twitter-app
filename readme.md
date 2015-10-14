# Twitter Client demo for Mirage

This is an example client application implementation for service virtualization tool [Mirage](https://github.com/SpectoLabs/mirage). 
It can either call Twitter API directly and by calling [Mirage Twitter Proxy](https://github.com/SpectoLabs/twitter-proxy)
(which provides recording and playback functionality). It expects to find proxy on port 8300 (localhost).
  
 
## Installation

* glide.yaml file holds all required dependencies. You can use glide to initialize your environment (glide install)  
* build a binary file (go build).
* Register your application in Twitter to get key and secret.
* Rename _conf.json.example_ to _conf.json_ and provide your previously acquired key and secret.
* Run application (./twitter-app

