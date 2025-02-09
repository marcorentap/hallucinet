# Hallucinet

Hallucinet allows you to use `.test` domain names for your Docker containers simply by connecting them on the `hallucinet` Docker network.

## Features
- Automatic DNS resolution for `.test` domains within Docker.
- No need to expose ports manually.
- Lightweight and easy to set up.

## Usage

### 1. Start Hallucinet

Run the following command to start Hallucinet:

```sh
docker compose up -d
```

This will start a DNS server on `192.168.100.2`.

### 2. Configure Your System to Use Hallucinet

#### On Systemd-Based Distributions
If your system uses `systemd-resolved`, modify `/etc/systemd/resolved.conf`:

1. Add the following line under the `[Resolve]` section:
   
   ```ini
   [Resolve]
   DNS=192.168.100.2
   ```

2. Restart the resolver service:
   
   ```sh
   sudo systemctl restart systemd-resolved
   ```

3. Verify that `192.168.100.2` is being used:

   ```sh
   resolvectl status
   ```

### 3. Running Containers with Hallucinet

To use Hallucinet's `.test` domains, connect your containers to the `hallucinet` network.

Run a container with:

```sh
docker run -it --rm --network hallucinet --name myapp nginx
```

Your container will now be accessible at `myapp.test`. 

