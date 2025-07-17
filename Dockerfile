FROM alpine:latest

# Install GitHub CLI
RUN apk add --no-cache github-cli

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Copy the binary
COPY gh-actions-tui /usr/local/bin/gh-actions-tui

# Make it executable
RUN chmod +x /usr/local/bin/gh-actions-tui

# Switch to non-root user
USER appuser

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/gh-actions-tui"]