mailman packages:
    pkg.installed:
        - pkgs:
            - mailman
fcgiwrap:
    service.running
/etc/nginx/cetana.d/10-mailman.conf:
    file.managed:
        - source: salt://mailman.nginx.conf
        - template: jinja
/etc/mailman/mm_cfg.py:
    file.managed:
        - source: salt://mm_cfg.py
        - template: jinja
