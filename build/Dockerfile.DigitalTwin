# Build Stage
FROM python:3.8-slim as builder
RUN apt-get update
RUN apt-get install -y --no-install-recommends build-essential
WORKDIR /workspace
COPY ./apps/digitaltwin/requirements.txt .
RUN pip install --user -r requirements.txt

# Production Stage
FROM python:3.8-slim
COPY --from=builder /root/.local /root/.local
COPY ./apps/digitaltwin/ /digitaltwin/
ENV PATH=/root/.local/bin:$PATH \
    PYTHONUNBUFFERED=1
ENTRYPOINT ["python3"]
CMD ["/digitaltwin/main.py"]
