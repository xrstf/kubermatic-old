FROM alpine:3.10

ENV MC_VERSION=RELEASE.2019-09-20T00-09-55Z \
    KUBECTL_VERSION=v1.16.2 \
    HELM_VERSION=v2.16.1 \
    VAULT_VERSION=1.2.3 \
    YQ_VERSION=2.4.0

RUN apk add --no-cache -U \
    bash \
    ca-certificates \
    curl \
    git \
    iproute2 \
    iptables \
    ipvsadm \
    jq \
    make \
    openssh-client \
    rsync \
    socat \
    unzip

RUN curl -Lo /usr/bin/yq https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_amd64 && \
    chmod +x /usr/bin/yq && \
    yq --version

RUN curl -Lo /usr/bin/mc https://dl.minio.io/client/mc/release/linux-amd64/archive/mc.${MC_VERSION} && \
    chmod +x /usr/bin/mc && \
    mc --version

RUN curl -Lo /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl && \
    chmod +x /usr/bin/kubectl && \
    kubectl version --short --client

RUN curl -L https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz | tar -xvz && \
    mv linux-amd64/helm /usr/bin/helm && \
    rm -rf linux-amd64 && \
    helm version --client --short

RUN curl -Lo vault.zip https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip && \
    unzip vault.zip && \
    rm vault.zip && \
    mv vault /usr/bin/vault && \
    vault version
