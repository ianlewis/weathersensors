FROM fluent/fluentd:latest
USER ubuntu
WORKDIR /home/ubuntu
ENV PATH /home/ubuntu/ruby/bin:$PATH
RUN gem install fluent-plugin-bigquery
EXPOSE 24284
CMD fluentd -c /fluentd/etc/$FLUENTD_CONF -p /fluentd/plugins $FLUENTD_OPT

# NOTE: the fluent/fluentd image has ONBUILD commands that expect fluent.conf
#       and a plugins directory with something in it.

# TODO: Can i just install the fluent-plugin-bigquery plugin in the plugins directory?
