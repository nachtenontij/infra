mattermost:
    user.present:
        - home: /opt/mattermost
extract mattermost:
    archive.extracted:
        - name: /opt/mattermost
        - source: https://releases.mattermost.com/3.10.0/mattermost-3.10.0-linux-amd64.tar.gz
        - user: mattermost
        - source_hash: 3977cb70b88a6def7009176bf23880fe5ad864cead05a1f2cae7792c8ac9148c
        - if_missing: /opt/mattermost/mattermost
mattermost mysql:
    mysql_user.present:
        - host: localhost
        - name: mattermost
        - password: {{ pillar['secrets']['mysql']['mattermost'] }}
    mysql_database.present:
        - name: mattermost
mattermost grant:
    mysql_grants.present:
        - grant: all privileges
        - database: mattermost.*
        - user: mattermost
/opt/mattermost/config.json:
    file.managed:
        - source: salt://mattermostConfig.json
        - user: mattermost
        - replace: false
        - mode: 600
        - template: jinja
/etc/nginx/sites-enabled/mattermost.conf:
    file.managed:
        - source: salt://mattermost.nginx.conf
        - template: jinja
