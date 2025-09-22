FROM python:3.10-slim

RUN apt-get update && apt-get install -y \
    libgl1 \
    libglib2.0-0 \
    libsm6 \
    libxext6 \
    libxrender1 \
    && rm -rf /var/lib/apt/lists/*

RUN pip install --upgrade pip setuptools wheel

RUN pip install --no-cache-dir \
    numpy==1.25.0 \
    pillow==10.0.0 \
    opencv-python==4.8.1.78 \
    tensorflow-cpu==2.15.0

COPY cmd/worker/worker.py /worker.py
COPY cmd/worker/segmentator_weights_converted.h5 /segmentator_weights_converted.h5