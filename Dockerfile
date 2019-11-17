FROM gcr.io/distroless/base
LABEL maintainer = "Jack Laxson <jackjrabbit@gmail.com>"
COPY spongebob_exe /app/
ENTRYPOINT ["/app/spongebob_exe"]