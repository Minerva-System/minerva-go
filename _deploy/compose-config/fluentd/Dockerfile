FROM fluentd:v1.16.2-debian-1.1
USER root
RUN gem install elasticsearch --no-document --version 8.1.2
RUN gem install fluent-plugin-elasticsearch --no-document --version 5.3.0
RUN gem install fluent-plugin-multi-format-parser --no-document --version 1.1.0
RUN gem install fluent-plugin-docker
USER fluent

