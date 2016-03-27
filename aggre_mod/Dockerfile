FROM scratch
# ca-certs are required for outbound HTTP connections.
ADD aggre_mod /aggre_mod
CMD ["/aggre_mod", "-access-token-path=/secrets/token"]
