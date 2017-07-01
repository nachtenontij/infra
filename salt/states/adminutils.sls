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
{% if grains['vagrant'] %}
/etc/profile.d/go.sh:
    file.managed:
        - source: salt://go.profile
{% endif %}
