# Holy Internet

Checks the internet connection at configed interval and sends email notification to your friends
and curses to your admin.

# Config

The app will look for a configuration file into the following locations:

* /etc/holy-internet/config.yml
* $HOME/.holy-internet/config.yml
* ./config.yml

```$xslt
saints:
    - 'facebook.com'
    - 'google.com'
    - 'wikipedia.com'
    - 'reddit.com'
    - 'instagram.com'
    - 'gmail.com'
    - 'github.com'
    - 'bing.com'
    - 'cnn.com'
    - 'quora.com'
    - 'youtube.com'
    - 'vimeo.com'
    - 'ask.com'
    - 'bitbucket.org'
pray:
    count: 3
    every: 30
curses:
    target: "your@admin.com"
    messages:
        - "Holy Cow! My internet was down!!!"
        - "Jesus Christ, where's my internet?"
followers:
    people:
        - "your@buddy.com"
    message: "My Internet was down for %s"

email:
    host: smtp.gmail.com
    username: your@username.com
    password: yourpassword
    port: 587
    tls: true
```
