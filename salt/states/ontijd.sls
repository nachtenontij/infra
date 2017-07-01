ontijd:
    user.present:
        - home: /home/ontijd
/etc/nginx/site.d/10-ontijd.conf:
    file.managed:
        - source: salt://ontijd.nginx.conf
        - template: jinja
/etc/systemd/system/ontijd.service:
    file.managed:
        - source: salt://ontijd.service
ontijd running:
    service.running:
        - name: ontijd
