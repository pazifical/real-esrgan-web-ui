```
docker build -t real_esrgan .
docker run --rm -it -v ./mount:/mount -p 8000:8000 --runtime=nvidia --gpus all --name real_esrgan real_esrgan
```