# NAT Type 

## NAT Variations

link to [RFC 3489](https://datatracker.ietf.org/doc/html/rfc3489#section-5)

### Full Cone

A full cone NAT is one where all requests from the
same internal IP address and port are mapped to the same external
IP address and port.  Furthermore, any external host can send a
packet to the internal host, by sending a packet to the mapped
external address

### Restricted Cone

A restricted cone NAT is one where all requests
from the same internal IP address and port are mapped to the same
external IP address and port.  Unlike a full cone NAT, an external
host (with IP address X) can send a packet to the internal host
only if the internal host had previously sent a packet to IP
address X

### Port Restricted Cone

A port restricted cone NAT is like a
restricted cone NAT, but the restriction includes port numbers.
Specifically, an external host can send a packet, with source IP
address X and source port P, to the internal host only if the
internal host had previously sent a packet to IP address X and
port P

### Symmetric

A symmetric NAT is one where all requests from the
same internal IP address and port, to a specific destination IP
address and port, are mapped to the same external IP address and
port.  If the same host sends a packet with the same source
address and port, but to a different destination, a different
mapping is used.  Furthermore, only the external host that
receives a packet can send a UDP packet back to the internal host