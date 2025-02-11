# Hallucinet  

Hallucinet allows you to use `.test` domain names for your Docker containers simply by connecting them to a Docker network.  

## Usage  

### 1. Create a Network and Set Hallucinet's IP Address  

Example `docker-compose.yaml`:  

```yaml
services:
  hallucinet:
    image: marcorentap/hallucinet
    build: .
    container_name: hallucinet
    restart: always
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    networks:
      hallucinet:
        ipv4_address: 192.168.100.2
      
networks:
  hallucinet:
    name: hallucinet
    ipam:
      config:
        - subnet: 192.168.100.0/24
```

This will start Hallucinet on `192.168.100.2`.  

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

To use Hallucinet's `.test` domains, connect your containers to the network.  

Run a container with:  

```sh
docker run -it --rm --network hallucinet --name myapp nginx
```

Your container will now be accessible at `myapp.test`.  

If you also want to make it accessible from other containers on the network, modify `/etc/docker/daemon.json` and add:  

```json
{
  "dns": ["192.168.100.2", "8.8.8.8"]
}
```

