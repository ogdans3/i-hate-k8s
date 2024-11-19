FROM debian:bullseye

LABEL maintainer="opsxcq@strm.sh"

# Install Python and clean up in a single RUN step to minimize image size
RUN apt-get update && \
    apt-get upgrade -y && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y python3 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY main.sh /main.sh

# Create /www directory
RUN mkdir /www

EXPOSE 80

WORKDIR /www

ENTRYPOINT ["/main.sh"]
