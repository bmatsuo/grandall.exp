##grandall

Grandall is a url aliasing system inspired by Justin Abrahms'
[randall](https://github.com/justinabrahms/randall).

##Usage

Bind a short URL to a frequented website.

    $ cat > ~/.config/grandall/sites/go-playground
    bind = "/play"
    url = "http://play.golang.org/"
    ^D
    $ sudo service grandalld restart

After grandalld is restarted visiting the aliased location will redirect to the
destination URL.

    http://<grandalld-host-addr>/play

On the command line the
[grandall](http://godoc.org/github.com/bmatsuo/grandall.exp/cmd/grandall)
utility may also be used to open aliased locations.

    $ grandall play

##"Manual" Setup

Install the grandalld binary.

    go get -u github.com/bmatsuo/grandall/cmd/grandalld

Create a configuration file somewhere like `~/.config/grandall/grandalld.conf`.

Create a sites directory somewhere like `~/.config/grandall/sites`.

Then start grandalld.

```sh
    grandalld \
        -config ~/.config/grandall/grandalld.conf \
        -sites ~/.config/grandall/sites
```

##Integration examples

The `examples/` directory contains examples for integrating grandalld and
ensuring it is always running.

- [lsb-init](https://github.com/bmatsuo/grandall.exp/tree/master/examples/lsb-init)
  contains an [LSB init script](https://wiki.debian.org/LSBInitScripts) for
  grandalld.

- OS X? I don't like launchd...
