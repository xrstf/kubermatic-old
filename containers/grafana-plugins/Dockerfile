FROM alpine:3.10

RUN mkdir -p /plugins
WORKDIR /plugins

RUN wget -O- https://grafana.com/api/plugins/grafana-piechart-panel/versions/1.3.9/download | unzip - && \
    mv grafana-piechart-panel-* grafana-piechart-panel

RUN wget -O- https://grafana.com/api/plugins/farski-blendstat-panel/versions/1.0.1/download | unzip - && \
    mv farski-blendstat-grafana-* farski-blendstat-panel

RUN wget -O- https://grafana.com/api/plugins/michaeldmoore-multistat-panel/versions/1.2.3/download | unzip - && \
    mv michaeldmoore-michaeldmoore-multistat-panel-* michaeldmoore-multistat-panel

RUN wget -O- https://grafana.com/api/plugins/vonage-status-panel/versions/1.0.9/download | unzip - && \
    mv Vonage-Grafana_Status_panel-* vonage-status-panel
