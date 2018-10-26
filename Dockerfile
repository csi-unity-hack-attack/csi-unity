FROM ubuntu:16.04
LABEL description="Unity CSI Driver"

COPY ./_output/csi-unity /csi-unity

RUN apt update && apt install -y nfs-common && apt clean

ENTRYPOINT ["/csi-unity", "/bin/bash"]
