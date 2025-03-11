This Poc is a simple proposal for a new API for the Tiger.

It doens't contain a real world example, it's just an example of how the API would look like.

The idea is to have a single endpoint that can be used to CRUD the relationship between a user entity and a policy entity.

The user entity is the entity that will be used to identify the user.

The policy entity is the entity that will be used to identify the policy.

The relationship between the user and the policy is many to one.

The user can have multiple policies and a policy has one user.

The user entity will have the following fields:

- id
- name
- email
- created_at
- updated_at

The policy entity will have the following fields:

- id
- user_id
- type
- amount_cents
- created_at
- updated_at
