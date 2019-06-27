# Sensu Performance Testing

This repository contains the Sensu performance testing assets used to
stress test and measure Sensu's capabilities. Performance testing is
done for every Sensu major and minor release to help guard against
performance regressions.

## The Sensu Testbed

The Sensu Testbed is comprised of five bare metal hosts and two
gigabit ethernet network switches. Bare metal is used for increased
control and consistency between testing runs (single tenant, no
hypervisor, etc.). One host is for running thousands of Sensu Agent
sessions (A1), three hosts are for running the Sensu Backend cluster
(B1, B2, B3), and the final host runs Postgres for the Sensu
Enterprise Event Store (P). One of the network switches is used for
SSH access to each node and the Sensu Agent sessions traffic. The
other network switch is used for Sensu Backend etcd and Postgres
traffic. The Postgres host uses three 1 gigabit ethernet cards,
round-robin bonded (bond0), to increase its network bandwidth.

![Network Diagram](https://raw.githubusercontent.com/sensu/sensu-perf/master/images/network.png)

### Hardware

#### Agents (agents1)

- AMD Ryzen Threadripper 2990WX Processor, 32 Cores, 3.0 GHz, 83MB Cache

- Gigabyte X399 AORUS PRO, DDR4 2666MHz, Triple M.2

- Corsair Vengeance LPX 32GB DDR4 2666MHz CL16 Quad Channel Kit (4x 8GB)

- Intel 660p Series M.2 PCIe 512GB Solid State Drive

- GeForce GT 710, 1GB DDR3

- Cooler Master Wraith Ripper Ryzen ThreadRipper CPU Cooler

- EVGA SuperNOVA 850W Power Supply

- Cooler Master MasterCase H500P Mesh E-ATX Case

#### Backends (backend1, backend2, backend3)

- AMD Ryzen Threadripper 2920X Processor, 12 Cores, 3.5GHz, 39MB Cache

- Gigabyte X399 AORUS PRO, DDR4 2666MHz, Triple M.2

- Corsair Vengeance LPX 16GB DDR4 2666MHz CL16 Dual Channel Kit (2x 8GB)

- Two Intel 660p Series M.2 PCIe 512GB Solid State Drives

- Intel Gigabit CT PCIe Network Card

- GeForce GT 710, 1GB DDR3

- Noctua NH-U12S TR4-SP3 CPU Cooler

- Corsair CX Series 650W Power Supply

- Corsair Carbide Series 270R Mid Tower ATX Case

#### Postgres (postgres)

- AMD Ryzen Threadripper 2920X Processor, 12 Cores, 3.5GHz, 39MB Cache

- Gigabyte X399 AORUS PRO, DDR4 2666MHz, Triple M.2

- Corsair Vengeance LPX 16GB DDR4 2666MHz CL16 Dual Channel Kit (2x 8GB)

- Two Intel 660p Series M.2 PCIe 512GB Solid State Drives

- Samsung 970 PRO NVMe M.2 PCIe 1TB Solid State Drive

- Three Intel Gigabit CT PCIe Network Cards

- GeForce GT 710, 1GB DDR3

- Noctua NH-U12S TR4-SP3 CPU Cooler

#### Network

- Two Ubiquiti UniFi 8 Port 60W Switches (US-8-60W)

- Eleven Cat 6 5ft Ethernet Cables
