FROM python:3.9-alpine

RUN apk update && apk add --no-cache \
    cmake \
    gcc \
    g++ \
    make 

WORKDIR /app

COPY . .
RUN pip install --upgrade pip
RUN pip install -r requirements.txt

EXPOSE 8000

CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--reload", "--port", "8000"]