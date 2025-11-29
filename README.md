# GracefulShutdown

Implementation of graceful shutdown presented by [FlorianWoelki](https://github.com/FlorianWoelki) in [this video](https://www.youtube.com/watch?v=UPVSeZXBTxI) 

## Table of contents
* [General info](#general-info)
* [Technologies](#technologies)
* [Status](#status)
* [Inspiration](#inspiration)

## General info

### Docker 

Just download

`docker pull arturmareknowak/gracefulshutdown-amd64` 

or 

`docker pull arturmareknowak/gracefulshutdown-arm64`

and run

`docker run -p 8080:8080 <image id>`

### Kubernetes

Just do run

`kubectl run gracefulshutdown --image=arturmareknowak/gracefulshutdown-amd64 --restart=Never`

or

`kubectl run gracefulshutdown --image=arturmareknowak/gracefulshutdown-arm64 --restart=Never`


## Technologies
* go 1.25

## Status
Project is: _finished_

## Inspiration
[FlorianWoelki](https://github.com/FlorianWoelki) and just life 