/.pachy-filter:
    file.managed:
        - source: salt://pachy-filter
backup authorized keys:
    ssh_auth.present:
        - user: root
        - names:
            - ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMsbyPatMkNdiMpvbOjq4XqS4LWWY6CsrGcstSA2Lrbg root@serf
