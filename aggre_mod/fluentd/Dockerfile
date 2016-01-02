FROM fluent/fluentd:latest
USER root
WORKDIR /home/ubuntu
ENV PATH /home/ubuntu/ruby/bin:$PATH
RUN gem install fluent-plugin-bigquery
RUN mkdir -p /var/log/fluent/
EXPOSE 24284
CMD fluentd -c /fluentd/etc/$FLUENTD_CONF -p /fluentd/plugins $FLUENTD_OPT

# NOTE: the fluent/fluentd image has ONBUILD commands that expect fluent.conf
#       and a plugins directory with something in it.
