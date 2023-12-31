# Real-ESRGAN-Web-UI

A containerized web UI with public API for upscaling images.

## Why?

This project is meant to be run on a server or at least containerized.

## Dependencies

Since Real-ESRGAN can utilize the graphics card, the NVIDIA Container Toolkit has to be installed and configured first:
https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/latest/install-guide.html

## Docker

### Build the container

```
docker build -t real_esrgan_web_ui .
```

### Run the container in interactive mode
```
docker run --rm -it -p 18080:8080 --runtime=nvidia --gpus all --name real_esrgan_web_ui real_esrgan_web_ui
```

### Run the container in the background
```
docker run --rm -d -p 8080:8080 --runtime=nvidia --gpus all --name real_esrgan_web_ui real_esrgan_web_ui
```
