FROM chromedp/headless-shell:latest

ENV DB_HOST localhost
EXPOSE 9222

WORKDIR /code
COPY gospider /gospider
RUN chmod +x /gospider

CMD ["/gospider"]