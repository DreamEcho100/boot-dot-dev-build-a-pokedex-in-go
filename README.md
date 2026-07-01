## Notes

When running locally _(on Ubuntu)_, a lot of times I get

```
Error executing command: Get "https://pokeapi.co/api/v2/location-area/?limit=20&offset=0": read tcp 192.168.1.6:36816->172.67.195.193:443: read: connection reset by peer
```

That happens when I'm connecting through Wi-Fi. it temporary works after `sudo systemctl restart NetworkManager` but, it worked with me through my pone hotspot better _(I think it is a problem with my ISP, I have no idea why it happens, but it is not related to the code of this project. also, I didn't test to see if it works using a VPN)_.
