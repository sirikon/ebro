FROM httpd:2.4
RUN apt-get update && apt-get install -y git
WORKDIR /var/www/git/test.git
COPY ./content/* .
RUN git config --global user.email "you@example.com"
RUN git config --global user.name "Your Name"
RUN git init && git add . && git commit -m "first commit"
COPY ./httpd.conf /usr/local/apache2/conf/httpd.conf
WORKDIR /usr/local/apache2
