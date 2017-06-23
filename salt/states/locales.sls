locales packages:
    pkg.installed:
        - pkgs:
            - locales
set locale:
    locale.present:
        - name: en_US.UTF-8
default_locale:
    locale.system:
        - name: en_US.UTF-8
        - require:
            - locale: set locale
