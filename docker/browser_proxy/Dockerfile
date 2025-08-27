FROM kasmweb/brave:1.17.0-rolling-daily

USER root

RUN apt-get update && apt-get install -y nginx

COPY nginx-chrome-proxy.conf /etc/nginx/chrome-debug-proxy.conf

COPY start_browser_with_proxy.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/start_browser_with_proxy.sh

USER kasm-user

CMD ["/usr/local/bin/start_browser_with_proxy.sh"]
