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

# Make port 4000 available to the world outside this container
EXPOSE 4000

# Define the command to run your application with Gunicorn
CMD ["gunicorn", "--bind", "0.0.0.0:4000", "app:create_app()"]
