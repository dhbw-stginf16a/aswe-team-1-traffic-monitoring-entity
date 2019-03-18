FROM scratch

ADD traffic-monitor .

EXPOSE 8080

CMD ["./traffic-monitor"]