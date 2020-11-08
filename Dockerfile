FROM golang:1.15-alpine3.12 AS build

RUN apk add \
	make

WORKDIR /go/src/app
COPY . .

RUN make TAGS=production


FROM alpine:3.12

WORKDIR /app
COPY --from=build /go/src/app/assets assets
COPY --from=build /go/src/app/templates templates
COPY --from=build /go/src/app/data/cities.json data/
COPY --from=build /go/src/app/datetime .

CMD ["./datetime"]
