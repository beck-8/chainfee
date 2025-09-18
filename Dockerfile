FROM ubuntu AS builder
WORKDIR /app
COPY . .
RUN apt update && apt install mesa-opencl-icd ocl-icd-opencl-dev gcc git jq pkg-config curl clang build-essential hwloc libhwloc-dev wget -y 
RUN wget -c https://golang.org/dl/go1.23.10.linux-amd64.tar.gz -O - | tar -xz -C /usr/local && \
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
RUN . ~/.bashrc && make

FROM ubuntu
ENV TZ=Asia/Shanghai
RUN apt-get update && apt-get install -y ca-certificates tzdata && \
    ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    rm -rf /var/cache/apt/*  && \
    mkdir -p /app   
COPY --from=builder /app/sectors_penalty /app/sectors_penalty
CMD /app/sectors_penalty
EXPOSE 8099