var APIBuilder = require('@axway/api-builder-runtime');
var Model = APIBuilder.createModel('RegisteredUsers', {
    "connector": "memory",
    "fields": {
        "name": {
            "type": "string",
            "required": true
        },
        "secret": {
            "type": "object",
            "required": true
        }
    },
    "actions": []
});
module.exports = Model;