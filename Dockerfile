FROM golang

ARG CI_BUILD_REF_NAME

RUN rm /bin/sh && ln -s /bin/bash /bin/sh

# install necessary libs
RUN apt-get update && \
    apt-get install -y --no-install-recommends libssl-dev

# install pip for aws cli
RUN wget 'https://bootstrap.pypa.io/get-pip.py' && python2.7 get-pip.py && pip install awscli

# setup ssh keys
RUN mkdir /root/.ssh && chmod 700 /root/.ssh
COPY docker/keys/id_rsa /root/.ssh/id_rsa
COPY docker/keys/id_rsa.pub /root/.ssh/id_rsa.pub
RUN chmod 600 /root/.ssh/id_rsa && chmod 644 /root/.ssh/id_rsa.pub

# accept gitlab.org ssh key
RUN mkdir -p ~/.ssh && ssh-keyscan gitlab.com >> ~/.ssh/known_hosts

# download git repository
WORKDIR /go/src/service-recordingStorage
COPY . .

# install dep and project dependencies
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN go install -v

# other stuff
EXPOSE 8094

COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod 776 /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]