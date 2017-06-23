mail packages:
    pkg.installed:
        - pkgs:
            - postfix-pcre
postfix:
    service:
        - running
/etc/postfix/main.cf:
    file.managed:
        - source: salt://mail/main.cf
        - template: jinja
        - watch_in:
            - service: postfix
/etc/postfix/virtual:
    file.directory
{% for file in ['transport', 'sender_canonical_map', 'virtual/map', 'virtual/domains' ] %}
/etc/postfix/{{ file }}:
    file.managed:
        - source: salt://mail/{{ file }}
        - template: jinja
    cmd.wait:
        - name: postmap /etc/postfix/{{ file }}
        - watch:
            - file: /etc/postfix/{{ file }}
        - watch_in:
            - service: postfix
        - require:
            - file: /etc/postfix/virtual
            - pkg: mail packages
{% endfor %}
