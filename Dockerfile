FROM scratch

COPY triggerserver /app/
WORKDIR "/app"
EXPOSE 44444
CMD ["./triggerserver"]
