FROM fluent/fluentd:v0.12.29
MAINTAINER Ian Lewis <ianmlewis@gmail.com>

WORKDIR /home/fluent
ENV PATH /home/fluent/.gem/ruby/2.3.0/bin:$PATH

USER root
# For some reason the json gem doesn't seem to be installed though it
# looks like it's being installed in the parent's Dockerfile.
# For some reason active support doesn't install bigdecimal either?
RUN apk --no-cache --update add sudo build-base ruby-dev && \
    apk --no-cache add ca-certificates && \
    sudo -u fluent gem install json && \
    sudo -u fluent gem install fluent-plugin-secure-forward && \
    sudo -u fluent gem install bigdecimal && \
    sudo -u fluent gem install fluent-plugin-bigquery -v 0.2.16 && \
    rm -rf /home/fluent/.gem/ruby/2.3.0/cache/*.gem && sudo -u fluent gem sources -c && \
    apk del sudo build-base ruby-dev && rm -rf /var/cache/apk/*

COPY fluent.conf /fluentd/etc/

EXPOSE 24284

CMD exec fluentd -c /fluentd/etc/$FLUENTD_CONF -p /fluentd/plugins $FLUENTD_OPT

# NOTE: the fluent/fluentd image has ONBUILD commands that expect fluent.conf
#       and a plugins directory with something in it.
