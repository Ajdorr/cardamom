FROM ubuntu:22.04

ENV CARDAMOM_CONFIG=config.yaml 

RUN apt-get update && apt-get install -y curl

COPY build/cardamom-core config.yaml .

CMD ["./cardamom-core"]