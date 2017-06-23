adminutils packages:
    pkg.installed:
        - pkgs:
            - htop
            - iftop
            - iotop
            - ncdu
            - screen
            - vim
            - psmisc
            - socat
/etc/vim/vimrc.local:
    file.managed:
        - source: salt://vimrc
