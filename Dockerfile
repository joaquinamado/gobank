FROM golang:1.23.4

# Install air
RUN go install github.com/air-verse/air@latest

RUN apt-get update && \ 
    apt-get install --no-install-recommends -y \
    apt-utils ghostscript fontforge cabextract wget

ARG WKHTML2PDF_VERSION=0.12.6.1-3

# Obtener VERSION_CODENAME y descargar .deb
RUN . /etc/os-release && \
    wget https://github.com/wkhtmltopdf/packaging/releases/download/${WKHTML2PDF_VERSION}/wkhtmltox_${WKHTML2PDF_VERSION}.${VERSION_CODENAME}_amd64.deb -O wkhtmltox.deb

# Instalar el paquete .deb
RUN dpkg -i wkhtmltox.deb || apt-get install -f -y

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air"]

