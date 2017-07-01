nginx packages:
    pkg.installed:
        - pkgs:
            - nginx
            - fcgiwrap
            - letsencrypt
/etc/nginx/sites-enabled/default:
    file.absent
/srv/default/htdocs:
    file.directory:
        - makedirs: true
/etc/nginx/site.d:
    file.directory
/etc/nginx/backends:
    file.directory
/etc/nginx/backends/fcgiwrap:
    file.managed:
        - source: salt://fcgiwrap.nginx-backend
/etc/nginx/sites-enabled/site.conf:
    file.absent
/etc/nginx/sites-enabled/00-site.conf:
    file.managed:
        - source: salt://site.nginx.conf
        - template: jinja
nginx running:
    service.running:
        - name: nginx
        - watch:
            - file: /etc/nginx/sites-enabled/*.conf
            - file: /etc/nginx/site.d/*.conf
