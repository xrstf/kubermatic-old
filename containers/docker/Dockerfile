FROM docker:18.09.0-dind

RUN apk add --no-cache -U ca-certificates wget unzip jq bash
RUN wget https://releases.hashicorp.com/vault/1.1.3/vault_1.1.3_linux_amd64.zip && \
    unzip vault_1.1.3_linux_amd64.zip && \
    rm vault_1.1.3_linux_amd64.zip && \
    mv vault /bin/
