# Use an official Node.js runtime as a parent image
FROM node:16-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json to the container
COPY package*.json ./

# Install project dependencies
RUN npm install

# Copy the entire application code to the container
COPY . .

# Expose a port for the application to listen on (usually 80)
EXPOSE 4000

# Define the command to run your React app
CMD ["npm", "start"]
