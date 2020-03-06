FROM osixia/openldap:latest

# Install the necessary tools
RUN apt-get -y update \
	&& LC_ALL=C DEBIAN_FRONTEND=noninteractive apt-get -t buster install -y --no-install-recommends \
	ca-certificates \
	build-essential \
	libltdl-dev \
	groff-base \
	git \
	&& apt-get clean \
	&& rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /opt/openldap-build
# Clone the openldap and openldap-bcrypt repos
RUN git clone https://github.com/openldap/openldap.git
WORKDIR /opt/openldap-build/openldap/contrib/slapd-modules/passwd
RUN git clone https://github.com/wclarie/openldap-bcrypt.git bcrypt

# Configure and build openldap
WORKDIR /opt/openldap-build/openldap
RUN ./configure --prefix=/opt/openldap-build --enable-modules --enable-shared
RUN make depend && make && make install

# Make the bcrypt module
WORKDIR /opt/openldap-build/openldap/contrib/slapd-modules/passwd/bcrypt
RUN LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/opt/openldap-build make && make install

