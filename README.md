# pfcpsim for Travelping UPG-VPP

This repository is a compatibility build of
[omec-project/pfcpsim](https://github.com/omec-project/pfcpsim) for PFCP session
provisioning and GTP-U load testing against
[Travelping UPG-VPP](https://github.com/travelping/upg-vpp).

The `main` branch is based on upstream pfcpsim `v1.4.4`. It keeps the upstream
Apache-2.0 license and history, with a small set of changes required by
UPG-VPP `v2.0.0`.

## Container Image

```text
ghcr.io/infinitydon/pfcpsim-travelping:v1.4.4-7
```

Build locally:

```sh
docker build -t pfcpsim-travelping:local .
```

## Compatibility Changes

The following differences from upstream `v1.4.4` are intentional:

1. Add PFCP Network Instance IEs:
   - `access` for access-side PDRs and FARs.
   - `internet` for core-side PDRs and FARs.
2. Mark the UE IPv4 address in downlink PDRs as a destination address using
   the PFCP `SD` flag.
3. Do not add a zero-address Outer Header Creation IE to the initial dropped
   downlink FAR. UPG-VPP rejects that placeholder.
4. Respond to PFCP Heartbeat Requests so the UPF can retain the association.
5. Omit periodic URRs from bulk session establishment. They can generate
   enough reports to exhaust pfcpsim's synchronous response timeout while
   creating many sessions.
6. Set session and application QER MBR values to 1 Gbps so load tests measure
   the data path rather than the low upstream test defaults.

These changes target forwarding load tests. This branch is not intended to
replace upstream pfcpsim for general PFCP conformance testing.

## Running

The container includes both `pfcpsim` and `pfcpctl`.

```sh
docker run --rm --network host \
  --name pfcpsim-travelping \
  ghcr.io/infinitydon/pfcpsim-travelping:v1.4.4-7 \
  --interface n4
```

Configure and associate it with UPG-VPP:

```sh
docker exec pfcpsim-travelping pfcpctl -s localhost:54321 \
  service configure \
  --n3-addr 10.0.3.10 \
  --remote-peer-addr 10.0.4.9:8805

docker exec pfcpsim-travelping pfcpctl -s localhost:54321 \
  service associate
```

Create 1,000 sessions:

```sh
docker exec pfcpsim-travelping pfcpctl -s localhost:54321 \
  session create \
  --count 1000 \
  --baseID 1 \
  --gnb-addr 10.0.3.1 \
  --ue-pool 48.0.0.0/16 \
  --qfi 9
```

The generated uplink TEIDs begin at `1` and increment by `10`. UE addresses
begin at `48.0.0.1`.

## Kubernetes

For a Multus-attached N4 interface, run the simulator with:

```yaml
command: ["pfcpsim"]
args: ["--interface", "n4"]
```

The associated Helm chart is maintained separately in the
`trex-travelping` workspace. It binds the N4 PCI device, provisions 1,000
sessions, and runs a Cisco TRex GTP-U profile.

The validated setup used:

- Travelping UPG-VPP `v2.0.0`
- pfcpsim base `v1.4.4`
- Cisco TRex `v3.06`
- 1,000 PFCP sessions
- 100,000 packets per second for 15 seconds
- 1,500,001 transmitted test packets
- 0% measured packet loss

## Development

```sh
go test ./...
docker build -t pfcpsim-travelping:dev .
```

The upstream remote can be retained as:

```sh
git remote add upstream https://github.com/omec-project/pfcpsim.git
```

## Upstream And License

Original project:
[github.com/omec-project/pfcpsim](https://github.com/omec-project/pfcpsim)

This repository preserves the upstream Apache-2.0 licensing information in
`LICENSES/` and the source file SPDX headers.
