ARG IMAGE
FROM ${IMAGE}
ENV PLANTUML_VERSION=1.2023.1
ENV LANG en_US.UTF-8
RUN apk add --no-cache graphviz wget ca-certificates && \
    apk add --no-cache graphviz wget ca-certificates ttf-dejavu fontconfig && \
    wget "http://downloads.sourceforge.net/project/plantuml/${PLANTUML_VERSION}/plantuml.${PLANTUML_VERSION}.jar" -O plantuml.jar && \
    apk del wget ca-certificates && \
    mkdir ~/.fonts; cd ~/.fonts;wget https://moji.or.jp/wp-content/ipafont/IPAexfont/IPAexfont00401.zip; unzip IPAexfont00401.zip; fc-cache -fv