FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD aggre_mod /aggre_mod
CMD ["/aggre_mod", "-access-token-path=/secrets/token"]
