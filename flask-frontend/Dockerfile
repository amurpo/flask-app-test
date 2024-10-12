# Base image for building the application
FROM python:3.13-slim AS builder

# Set the working directory within the container
WORKDIR /app

# Copy the requirements file into the container at /app
COPY requirements.txt .

# Install any needed packages specified in requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Copy the current directory contents into the container at /app
COPY . /app

# Compile translations
RUN pip install --no-cache-dir babel
RUN pybabel compile -d translations

# Install gunicorn
RUN pip install gunicorn

# Final image based on Python slim image
FROM python:3.13-slim

# Create a non-root user
RUN adduser --disabled-password --gecos '' appuser

# Set the working directory
WORKDIR /app

# Copy the compiled application and dependencies from the builder stage
COPY --from=builder /app /app
COPY --from=builder /usr/local/lib/python3.13/site-packages /usr/local/lib/python3.13/site-packages
COPY --from=builder /usr/local/bin/gunicorn /usr/local/bin/gunicorn

# Set the owner of the directory to the non-root user
RUN chown -R appuser:appuser /app

# Switch to the non-root user
USER appuser

# Expose the port for the application
EXPOSE 4000

# Define the command to run your application with Gunicorn
CMD ["gunicorn", "--bind", "0.0.0.0:4000", "app:create_app()"]
