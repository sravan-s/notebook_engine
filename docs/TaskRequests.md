## Task Requests

Basic format:
```
{
    "action": ActionType,
    "payload": Payload,
}
```

## Tasks

1. Create VM
{
    action: "CREATE_VM",
    payload: {
        notebook_id: i64,
    } 
}

2. Stop VM
{
    action: "STOP_VM",
    payload: { notebook_id } 
}

3. Run paragraph
{
    action: "RUN_PARAGRAPH",
    payload: { notebook_id, paragraph_id, code }
}

