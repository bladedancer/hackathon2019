const sdk = require('axway-flow-sdk');
const action = require('./action');

function getFlowNodes() {
	const flownodes = sdk.init(module);

	flownodes.add('speakeasy', {
		name: 'Speakeasy',
		icon: 'icon.svg',
		description: 'Time-base one-time password support using Speakeasy.',
		category: 'security'
	});

	// Secret method
	flownodes
		.method('secret', {
			name: 'Generate Secret',
			description: 'Create a new secret.'
		})
		.parameter('name', {
			description: 'The name to associate with the secret.',
			type: 'string'
		})
		.output('next', {
			name: 'Next',
			description: 'Return the fresh secret.',
			context: '$.secret',
			schema: {
				type: 'object'
			}
		})
		.action(action.secret);

	// Generate a token
	flownodes
		.method('token', {
			name: 'Generate Token',
			description: 'Generate a timebased onetime password.'
		})
		.parameter('secret', {
			description: 'Base32 encoded secret.',
			type: 'string'
		}, true)
		.parameter('ttl', {
			description: 'The validity period of the token.',
			type: 'number'
		}, false)
		.output('next', {
			name: 'Next',
			description: 'Return the fresh token.',
			context: '$.token',
			schema: {
				type: 'string'
			}
		})
		.action(action.token);

	// Validate the token
	flownodes
		.method('validate', {
			name: 'Validate Token',
			description: 'Validate a token.'
		})
		.parameter('secret', {
			description: 'Base32 encoded secret.',
			type: 'string'
		}, true)
		.parameter('token', {
			description: 'Token.',
			type: 'string'
		}, true)
		.output('invalid', {
			name: 'Invalid',
			description: 'The token is invalid.'
		})
		.output('valid', {
			name: 'Valid',
			description: 'The token is valid.'
		})
		.action(action.validate);
	return Promise.resolve(flownodes);
}

exports = module.exports = getFlowNodes;
