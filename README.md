# etcd-snapshot-to-json

`etcd-snapshot-to-json` is a lightweight tool for displaying `etcd` snapshots in a JSON format.

This repository has an associated container image for convenience.

## Image Overview

- **Image**: `spurin/etcd-snapshot-to-json:latest`
- **Supported Architectures**: `linux/amd64`, `linux/arm64/v8`, `linux/ppc64le`, `linux/s390x`
- **Tool Function**: Converts `etcd` snapshot files to JSON output

## Limitations

While `etcd-snapshot-to-json` aims to provide a human-readable JSON representation of the `etcd` snapshot data, it is important to note that this conversion from the original protobuf format might not be entirely lossless. Certain data types and structures in protobuf may not have direct JSON equivalents, which can lead to subtle discrepancies or loss of fidelity in the translation process. As such, this tool is best used for inspection and debugging purposes rather than for backing up or replicating `etcd` data in a lossless manner.

## Usage

### Running the Container with a Snapshot File

To use `etcd-snapshot-to-json` mount a directory containing your `etcd` snapshot file and run the container, specifying the path to the snapshot file inside the container.

#### Command-Line Options

- `-keys`: Optional. Comma-separated list of keys to filter in the output.
- `--latest`: Optional. When used with `-keys`, outputs only the highest version of the specified keys.

#### Example Command

```bash
docker run --rm -v /path/to/snapshot:/snapshots spurin/etcd-snapshot-to-json:latest /snapshots/snapshot.db --keys key1,key2 --latest
```

- **Explanation**:
  - `--rm`: Automatically removes the container after it finishes.
  - `-v /path/to/snapshot:/snapshots`: Mounts the local directory containing `snapshot.db` to `/snapshots` in the container.
  - `spurin/etcd-snapshot-to-json:latest`: Specifies the container image.
  - `/snapshots/snapshot.db`: The path to the snapshot file within the container.
  - `-keys key1,key2`: Filters output to include only key1 and key2.
  - `--latest`: Ensures only the highest version of each specified key is included.

#### Example Output

Running the command will output the snapshot contents in JSON format:

```json
[
    {
        "key": "key1",
        "value": "value1",
        "create_revision": 1024,
        "mod_revision": 1040,
        "version": 4
    },
    {
        "key": "key2",
        "value": "value2",
        "create_revision": 1032,
        "mod_revision": 1078,
        "version": 3
    }
]
```

### Example Usage

Assume you have an `etcd` snapshot file located at `/home/user/etcd_snapshots/snapshot.db` on your local machine. To view this snapshot as JSON, run:

```bash
docker run --rm -v /home/user/etcd_snapshots:/snapshots spurin/etcd-snapshot-to-json:latest /snapshots/snapshot.db
```

The output will be a JSON-formatted string representing the key-value data stored in the `etcd` snapshot.

## Building the Image Locally (Optional)

If you want to build the image locally for testing or development, clone this repository and use the following command:

```bash
docker build -t etcd-snapshot-to-json .
```

Then, run the image with:

```bash
docker run --rm -v /path/to/snapshot:/snapshots etcd-snapshot-to-json /snapshots/snapshot.db
```

## Copying the Binary in a Dockerfile

To copy the `etcd-snapshot-to-json` binary directly into another Dockerfile based image, you can use `COPY --from` option.

```dockerfile
COPY --from=spurin/etcd-snapshot-to-json:latest /etcd-snapshot-to-json /usr/local/bin/etcd-snapshot-to-json
```

## Copying the Binary Locally (via Docker)

If you want to copy the binary to your local machine without creating a new image, you can use the `docker cp` command with a temporary container:

1. **Create a Temporary Container**:
   ```bash
   docker create --name temp_container spurin/etcd-snapshot-to-json:latest
   ```

2. **Copy the Binary from the Container**:
   ```bash
   docker cp temp_container:/etcd-snapshot-to-json ./etcd-snapshot-to-json
   ```

3. **Remove the Temporary Container**:
   ```bash
   docker rm temp_container
   ```
