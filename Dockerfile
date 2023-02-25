FROM public.ecr.aws/lambda/provided:al2 as build

# install go
RUN yum install -y golang && yum clean all

# build
WORKDIR /src
COPY . .
COPY main.go main.go
RUN go build -o /src/main

# copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2

COPY --from=build /src/main /main

ENTRYPOINT [ "/main" ]