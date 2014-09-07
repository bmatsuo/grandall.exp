##grandalld

A simple grandalld Dockerfile.

##Simple deployment

The simplest usage binds a directory of site aliases to the container.  In the
example below, sites are defined in `~/.config/grandall/sites/`.

    docker run --name grandall.exp \
        -p 8080:8080 \
        -v ~/.config/grandall/sites:/etc/grandall/sites-enabled \
        bmatsuo/grandall.exp

##Complex deployment

Similar to nginx, grandall allows maintainence of both 'available' and
'enabled' aliases.  The idea is that the enabled aliases are symbolic links to
available aliases.

    docker run --name grandall.exp-complex \
        -p 8080:8080 \
        -v ~/.config/grandall/sites-enabled:/etc/grandall/sites-enabled \
        -v ~/.config/grandall/sites-available:/etc/grandall/sites-available \
        bmatsuo/grandall.exp

The easiest way to achieve this is to have adjacent folders on the host machine
named "sites-enabled" and "sites-available".  The "sites-enabled" directory
contains symlinks of the form `../sites-available/*`.  This allows symlinks to
work inside the container.

    cd ~/.config/grandall/sites-enabled
    ln -s ../sites-available/foo foo
