# App to test databox store performance

Chris Greenhalgh, The University of Nottingham

status: starting...

## test plan

(1/minutes = 1,440/day, 43,200/month, 518,400/year)
(1/second = 3,600/hour, 86,400/day, 2,592,000/month, 31,104,000/year)

- insert time vs no items previously inserted (0, 10, 1000, 10,000, 100,000, 1,000,000)
- insert time vs item size (text 0, 10b, 100b, 1000b, 10,000b, 100,000b)
- read time, latest vs no items in store
- read time, latest vs item size
- read time, range, one item vs timestamp (how long ago) (cf no items)
- read time, range vs number of items in range
- read rate vs no concurrent readers (1, 10, 100)

1) store-json
2) store-timeseries
3) core-store

## Build

```
docker build -t app-test-store -f Dockerfile.dev .
```

Upload manifest.
