This tool can be used to benchmark your QuadrilleDB deployment. It can also be used in local to get a sense of the performance it can offer.

To run the benchmarks, use below commands. (You can download the binary from Download link on Quadrille documentation website or build the tool from source)

QuadrilleDB running on localhost

```bash
quadrille-benchmark
```

QuadrilleDB running on some other host

```bash
quadrille-benchmark leader_hostname:tcp_port
```

The tcp_port is set to 5679 by default