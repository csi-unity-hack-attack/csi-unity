#TODO
FROM ubuntu:16.04

COPY _output/csi-unity /csi-unity

ENTRYPOINT ["/csi-unity"]
