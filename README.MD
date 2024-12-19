Log format:
```json

{
    Host: hostname or ip (WIP)
    Timestamp: unix ts
    Type: str (HbSeqErr, HbTimeout, HbGood, etc)
    Payload : {
        <this is the unique error struct>
    }
}
```

