##grandall

Grandall is a self-managed url shortener inspired by Justin Abrahms'
[randall](https://github.com/justinabrahms/randall).

##Usage

Bind a short URL to a frequented website.

    $ cat > ~/.config/grandall/sites
    bind = "/goplay"
    url = "http://play.golang.org/"
    ^D
    $

After grandalld is restarted visiting the bound URL will redirect to the
destination URL.

    open http://<grandalld-host-addr>/goplay

##Setup

Install the grandalld binary.

    go get github.com/bmatsuo/grandall/cmd/grandalld

Create a configuration file somewhere, say `~/.config/grandall/grandalld.conf`.

Create a sites directory somewhere, say `~/.config/grandall/sites`.

Then start grandalld.

```sh
    grandalld \
        -config ~/.config/grandall/grandalld.conf \
        -sites ~/.config/grandall/sites
```
