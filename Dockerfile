FROM python:3.11

# RUN apt update
# RUN apt install git -y
WORKDIR /app

# Real-ESRGAN
RUN git clone https://github.com/xinntao/Real-ESRGAN.git

WORKDIR /app/Real-ESRGAN
RUN pip install basicsr
RUN pip install facexlib
RUN pip install gfpgan
RUN pip install -r requirements.txt
RUN python setup.py develop

# API
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY main.py .

RUN apt-get update 
RUN apt-get install ffmpeg libsm6 libxext6 wget -y

RUN wget https://github.com/xinntao/Real-ESRGAN/releases/download/v0.1.0/RealESRGAN_x4plus.pth -P /app/Real-ESRGAN/weights


CMD [ "uvicorn", "main:app", "--host","0.0.0.0"]