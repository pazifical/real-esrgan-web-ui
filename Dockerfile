# Building backend

FROM golang:1.21 AS build

WORKDIR /app
COPY main.go .
COPY go.mod .
COPY internal internal
RUN go build .

# Building Real-ESRGAN container

FROM python:3.11

WORKDIR /app
RUN git clone https://github.com/xinntao/Real-ESRGAN.git
WORKDIR /app/Real-ESRGAN
RUN pip install basicsr
RUN pip install facexlib
RUN pip install gfpgan
RUN pip install -r requirements.txt
RUN python setup.py develop
RUN apt-get update && apt-get install ffmpeg libsm6 libxext6 wget -y && apt-get clean
RUN wget https://github.com/xinntao/Real-ESRGAN/releases/download/v0.1.0/RealESRGAN_x4plus.pth -P /app/Real-ESRGAN/weights

WORKDIR /app
COPY --from=build /app/backend .
COPY static static

CMD [ "/app/backend"]