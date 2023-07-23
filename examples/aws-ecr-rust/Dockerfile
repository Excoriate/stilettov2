# Use the official Rust image as the Builder
FROM rust:1.53 as builder

# Create a new empty shell project
RUN USER=root cargo new --bin hello
WORKDIR /hello

# Copy our source code
COPY src ./src

# Build for release
RUN cargo build --release

# Our final base
FROM debian:buster-slim

# Copy the build artifact from the builder stage and create a new user.
RUN useradd rust
USER rust
COPY --from=builder /hello/target/release/hello /usr/local/bin

# Command to run when the container starts
CMD ["hello"]
