# etcd_snapshot_to_json

`etcd_snapshot_to_json` is a lightweight tool for converting `etcd` snapshots to JSON format. This tool allows you to inspect `etcd` snapshot contents by outputting key-value data in a human-readable JSON format.

This repository has an associated container image for convenience.

## Image Overview

- **Image**: `spurin/etcd_snapshot_to_json:latest`
- **Supported Architectures**: `linux/amd64`, `linux/arm64/v8`, `linux/ppc64le`, `linux/s390x`
- **Tool Function**: Converts `etcd` snapshot files to JSON output

## Usage

### Running the Container with a Snapshot File

To use `etcd_snapshot_to_json` mount a directory containing your `etcd` snapshot file and run the container, specifying the path to the snapshot file inside the container.

#### Example Command

```bash
docker run --rm -v /path/to/snapshot:/snapshots spurin/etcd_snapshot_to_json:latest /snapshots/snapshot.db
```

- **Explanation**:
  - `--rm`: Automatically removes the container after it finishes.
  - `-v /path/to/snapshot:/snapshots`: Mounts the local directory containing `snapshot.db` to `/snapshots` in the container.
  - `spurin/etcd_snapshot_to_json:latest`: Specifies the Docker image.
  - `/snapshots/snapshot.db`: The path to the snapshot file within the container.

#### Example Output

Running the command will output the snapshot contents in JSON format:

```json
{
  "key1": "value1",
  "key2": "value2",
  ...
}
```

### Example Usage

Assume you have an `etcd` snapshot file located at `/home/user/etcd_snapshots/snapshot.db` on your local machine. To view this snapshot as JSON, run:

```bash
docker run --rm -v /home/user/etcd_snapshots:/snapshots spurin/etcd_snapshot_to_json:latest /snapshots/snapshot.db
```

The output will be a JSON-formatted string representing the key-value data stored in the `etcd` snapshot.

## Building the Image Locally (Optional)

If you want to build the image locally for testing or development, clone this repository and use the following command:

```bash
docker build -t etcd_snapshot_to_json .
```

Then, run the image with:

```bash
docker run --rm -v /path/to/snapshot:/snapshots etcd_snapshot_to_json /snapshots/snapshot.db
```

## Copying the Binary in a Dockerfile

To copy the `etcd_snapshot_to_json` binary directly into another Dockerfile based image, you can use `COPY --from` option.

```dockerfile
COPY --from=spurin/etcd_snapshot_to_json:latest /etcd_snapshot_to_json /usr/local/bin/etcd_snapshot_to_json
```

## Copying the Binary Locally (via Docker)

If you want to copy the binary to your local machine without creating a new image, you can use the `docker cp` command with a temporary container:

1. **Create a Temporary Container**:
   ```bash
   docker create --name temp_container spurin/etcd_snapshot_to_json:latest
   ```

2. **Copy the Binary from the Container**:
   ```bash
   docker cp temp_container:/etcd_snapshot_to_json ./etcd_snapshot_to_json
   ```

3. **Remove the Temporary Container**:
   ```bash
   docker rm temp_container
   ```
