FROM gcr.io/distroless/base
LABEL maintainer = "Jack Laxson <jackjrabbit+spongebob_bot@gmail.com>"
COPY spongebob_exe /app/
ENTRYPOINT ["/app/spongebob_exe"]