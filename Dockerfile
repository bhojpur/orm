FROM moby/buildkit:v0.9.3
WORKDIR /orm
COPY orm README.md /orm/
ENV PATH=/orm:$PATH
ENTRYPOINT [ "/bhojpur/orm" ]