# BackerHub

## Technology

## API Testing

```js
// example scripts
let res = pm.response.json();

pm.test("Status code is 400", function () {
    pm.response.to.have.status(400);
});

pm.test("Response has required fields", function () {
    pm.expect(res).to.have.property("success");
    pm.expect(res).to.have.property("message");
    pm.expect(res).to.have.property("errors");
});

pm.test("Message check", function () {
    pm.expect(res.message).to.eql("Name is required");
});

pm.collectionVariables.set("USER_1_ID", res.data.id);
```

## References
