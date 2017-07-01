go packages:
    pkg.installed:
        - pkgs:
            - golang
            - git
/usr/local/go:
    file.directory

{% if grains['vagrant'] %}
/usr/local/go/src/github.com/nachtenontij/infra:
    file.directory:
        - makedirs: True
        - user: vagrant
    mount.mounted:
        - device: /vagrant
        - opts: bind
        - fstype: none
{% endif %}

{% for cmd in ['ontijd'] %}
build {{ cmd }}:
    cmd.run:
        - name: . /etc/profile && go get github.com/nachtenontij/infra/cmd/{{ cmd }}
{% endfor %}
