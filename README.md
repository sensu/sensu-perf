# Sensu Performance Testing

This repository contains the Sensu performance testing assets used to
stress test and measure Sensu's capabilities. Performance testing is
done for every Sensu major and minor release to help guard against
performance regressions.

## The Sensu Testbed

TODO: Diagram

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

- Intel Gigabit CT PCI-E Network Card

- GeForce GT 710, 1GB DDR3

- Noctua NH-U12S TR4-SP3 CPU cooler

- Corsair CX Series 650W Power Supply

- Corsair Carbide Series 270R Mid Tower ATX Case

#### Postgres (postgres)

- AMD Ryzen Threadripper 2920X Processor, 12 Cores, 3.5GHz, 39MB Cache

- Gigabyte X399 AORUS PRO, DDR4 2666MHz, Triple M.2

- Corsair Vengeance LPX 16GB DDR4 2666MHz CL16 Dual Channel Kit (2x 8GB)

- Two Intel 660p Series M.2 PCI-E 512GB Solid State Drives

- Samsung 970 PRO NVMe M.2 PCI-E 1TB Solid State Drive

- Three Intel Gigabit CT PCI-E Network Cards

- GeForce GT 710, 1GB DDR3

- Noctua NH-U12S TR4-SP3 CPU cooler

#### Network

- Two Ubiquiti UniFi Switch 8 60W

- Eleven Cat 6 5ft Ethernet Cables
