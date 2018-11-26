# PMR-Kafka-Writer
Example Configuration:
```
broker.address = localhost:9092;
topic = test;
sasl = {
        username = admin;
        password = admin-secret;
}
webserver = {
        path = "/";
        port = 81;
        basicauthfile = "auth.file";
        tls = {
                certfile = "cert.pem"
                keyfile = "key.pem"
        }
}
```

- `broker.address` required; can be a list of strings,  separated by comma.
- `topic` required; is the topic where the producers works on
- `sasl.username` & `sasl.password` optional; when not set sasl will not be used for kafka
- `webserver.path` required; the path on which the requests will be handled
- `webserver.port` required; the port on which the webserver will be running
- `webserver.basicauthfile` optional; when set this file will be used for http basic auth
- `webserver.tls.certfile` & `webserver.tls.keyfile` optional; when set the server will listen on https, else on http

Example basic auth file:
```
test = $apr1$i0XX0wZf$BOTqstyd9oOki/7v4MENU/
john = $1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1
```
Each user is in a separated line. The first column is for the username and the second for the password hash.
