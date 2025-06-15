# Webhook Forwarder:
- I Made this project to learn GO. and also i am going to have webhooks related works in couple months.

# How it works?
1. You give `this server` webhook endpoint/{uuid} to your push service (stripe or whatever)
2. You subscribe to `this server` event_pool/{uuid} api (SSE)
3. internally i put webhook request in redis pubsub on uuid channel 
4. `this server` event_pool is subscribed to that channel and as soon as message comes , from that client `you` get result via SSE


# Client Side:
- [webhook_forwarder_client](https://github.com/amzker/webhook_forwarder_client).
 
- basic idea is that i connect to my server SSE and as soon as request comes i forward to local endpoint 
Example:
```bash
./client -server_host <forwarder_host> -listen_token <uuid> -local_forward_url <your_local_webhook_url>
```
