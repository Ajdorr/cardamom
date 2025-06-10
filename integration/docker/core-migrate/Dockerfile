FROM ubuntu:22.04

ENV CARDAMOM_CONFIG=config.yaml 

COPY build/cardamom-core-migrate config.yaml .

CMD ["./cardamom-core-migrate"]