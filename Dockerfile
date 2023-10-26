# Use the official Python runtime as a parent image
FROM python:3.11-slim

# Set the working directory within the container
WORKDIR /app

# Copy the requirements file into the container at /app
COPY requirements.txt . /app/

# Install any needed packages specified in requirements.txt
RUN pip install -r requirements.txt

# Copy the current directory contents into the container at /app
COPY . /app

# Make port 80 available to the world outside this container
EXPOSE 80

# Define the command to run your application with Gunicorn
CMD ["gunicorn", "--bind", "0.0.0.0:80", "app:create_app()"]
