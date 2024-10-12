# Use a specific version of Node.js alpine image for consistency
FROM node:16-alpine

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm ci --only=production

# Copy the rest of the application code
COPY --chown=appuser:appgroup . .

# Set correct permissions
RUN chown -R appuser:appgroup /app

# Use the non-root user
USER appuser

# Expose the port the app runs on
EXPOSE 3000

# Start the app
CMD ["npm", "start"]