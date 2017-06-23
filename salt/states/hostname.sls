/etc/mailname:
    file.managed:
        - contents:
            {{ grains['fqdn'] }}
