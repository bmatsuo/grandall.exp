##grandall [![Build Status](https://travis-ci.org/bmatsuo/grandall.exp.svg?branch=master)](https://travis-ci.org/bmatsuo/grandall.exp)

Grandall is a url aliasing system inspired by Justin Abrahms'
[randall](https://github.com/justinabrahms/randall).

##Usage

Bind a short URL to a frequented website.

    $ cat > ~/.config/grandall/sites/go-playground
    bind = "/play"
    url = "http://play.golang.org/"
    description = "go playground"
    ^D
    $ sudo service grandalld restart

After grandalld is restarted visiting the aliased location will redirect to the
destination URL.

    http://<grandalld-host-addr>/play

On the command line the
[grandall](http://godoc.org/github.com/bmatsuo/grandall.exp/cmd/grandall)
utility may also be used to open aliased locations.

    $ grandall play

##Service setup

Download a [release
archive](http://github.com/bmatsuo/grandall.exp/releases/latest) or use `go
get` to install grandalld.

    go get github.com/bmatsuo/grandall.exp/cmd/grandalld

See the grandalld
[documentation](http://godoc.org/github.com/bmatsuo/grandall.exp/cmd/grandalld)
for help configuring the service.

##Command line client

The `grandall` command line interface is an optional component that launches
aliases from the command line.  For more information see the grandall
[documentation](http://godoc.org/github.com/bmatsuo/grandall.exp/cmd/grandall).
The command line interface can be installed from a [release
archive](http://github.com/bmatsuo/grandall.exp/releases/latest) with `go get`.

    go get github.com/bmatsuo/grandall.exp/cmd/grandall

##Integration examples

The `examples/` directory contains examples for integrating grandalld and
ensuring it is always running.

- [lsb-init](https://github.com/bmatsuo/grandall.exp/tree/master/examples/lsb-init)
  provides an [LSB init script](https://wiki.debian.org/LSBInitScripts) for
  grandalld.

- [dockerfile](https://github.com/bmatsuo/grandall.exp/tree/master/examples/dockerfile)
  provides a simple Dockerfile configuration for grandalld.

- OS X? I don't like launchd...

##Adding sites to a remote/distributed environment

If running grandalld locally adding sites is trivial, create a new entry in the
sites directory and restart grandalld.  Running grandalld remotely or on
multiple servers/workstations makes managing sites more difficult.  I don't
think the optimal way to manage aliases is completely clear.

This section explores possible implementations for an alias management system.

###Read-write HTTP API

If the HTTP interface allows updates then securing the API becomes a (more)
serious issue.  Securing API access may be inevitable and thus it would be best
to just get this done.  It's not clear to me that is the case.

An read-write API also makes the current deployment scenario less ideal.
Allowing grandalld to create files itself doesn't really sound optimal.
Further, in order to support the sites-enabled/-available deployment scenario
grandalld would need to create relative symlinks.

Making an API work in a distributed enviroment is not trivial.  Distributing
change is orthogonal to allowing change and an API does aid in change
distribution.

###Confd

The idea of using confd intrigues me.  Grandalld would be agnostic to confd's
presence, keeping things simple and allowing it to keep it's file-based config
strategy.  The problem of API security is deferred to etcd/consul.

A confd based approach might make distribution quite easy.  It may even still
work sufficiently when machines are not part of the same private subnet.

I think there are some questions around this.

- Does etcd support authentication and encryption? Yes. Server/Client certs.
  HTTP basic auth via nginx proxy (see
[issues/245](https://github.com/coreos/etcd/issues/254)).

- Does consul support authentication and encryption?

- Is this a completely sufficient configuration mechanism.

- How much pain is required to run machines across an insecure network.
  Specifically, for an etcd/consul client to communicate with a server.

###Dropbox/Syncthing?

For a "local distributed" setup (one grandalld on each personal computer).  It
may be convenient to do some filesystem watching or triggered action and use a
sync service to distribute alias definitions.
