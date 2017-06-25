base:
    '*':
        {% if grains['vagrant'] %}
        # contains auto-generated passwords.  In production there are stored
        # out of the repository.
        - vagrant
        {% else %}
        - secrets
        {% endif %}
