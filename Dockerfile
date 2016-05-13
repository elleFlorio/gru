FROM ubuntu
COPY gru /bin/gru
ENTRYPOINT ["gru"]
CMD ["--help"]