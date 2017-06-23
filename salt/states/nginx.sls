nginx packages:
    pkg.installed:
        - pkgs:
            - nginx
            - fcgiwrap
            - letsencrypt
/etc/nginx/sites-enabled/default:
    file.absent
/etc/nginx/cetana.d:
    file.directory
/etc/nginx/backends:
    file.directory
/etc/nginx/backends/fcgiwrap:
    file.managed:
        - source: salt://fcgiwrap.nginx-backend
/etc/nginx/sites-enabled/cetana.conf:
    file.managed:
        - source: salt://site.nginx.conf
        - template: jinja
nginx running:
    service.running:
        - name: nginx
        - watch:
            - file: /etc/nginx/sites-enabled/cetana.conf
            - file: /etc/nginx/cetana.d/*.conf
