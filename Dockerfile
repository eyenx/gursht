FROM scratch

COPY gursht /
COPY etc/passwd /etc/passwd

USER nobody

ENTRYPOINT ["/gursht"]
